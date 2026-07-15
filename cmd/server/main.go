package main

import (
	"flag"

	"github.com/gin-contrib/cors"
	"go.uber.org/zap"

	"github.com/you/sharing-vision-backend-v2/internal/auth"
	"github.com/you/sharing-vision-backend-v2/internal/cache"
	"github.com/you/sharing-vision-backend-v2/internal/config"
	"github.com/you/sharing-vision-backend-v2/internal/db"
	"github.com/you/sharing-vision-backend-v2/internal/handler"
	"github.com/you/sharing-vision-backend-v2/internal/middleware"
	"github.com/you/sharing-vision-backend-v2/internal/repository"
	"github.com/you/sharing-vision-backend-v2/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	var migrationFlag string
	flag.StringVar(&migrationFlag, "migrate", "", "run migrations: up or down")
	flag.Parse()

	config.Load()
	defer config.Log.Sync()

	if migrationFlag != "" {
		if err := runMigrations(); err != nil {
			config.Log.Fatal("migrations failed", zap.Error(err))
		}
		return
	}

	runServer()
}

func runMigrations() error {
	dbConn, err := db.Connect()
	if err != nil {
		return err
	}
	defer dbConn.Close()
	return auth.RunMigrations(dbConn)
}

func runServer() {
	dbConn, err := db.Connect()
	if err != nil {
		config.Log.Fatal("db connect failed", zap.Error(err))
	}
	defer dbConn.Close()

	redisClient, err := cache.ConnectRedis()
	if err != nil {
		config.Log.Fatal("redis connect failed", zap.Error(err))
	}

	articleRepo := repository.NewArticleRepo(dbConn)
	userRepo := repository.NewUserRepo(dbConn)
	categoryRepo := repository.NewCategoryRepo(dbConn)
	articleCache := cache.NewCache(redisClient)
	articleService := service.NewArticleService(articleRepo, articleCache)
	authService := service.NewAuthService(userRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	if err := auth.RunMigrations(dbConn); err != nil {
		config.Log.Fatal("migrations failed", zap.Error(err))
	}
	if err := authService.BootstrapAdmin(); err != nil {
		config.Log.Warn("bootstrap admin skipped", zap.Error(err))
	}

	limiter := middleware.NewRateLimiter()
	authMw := middleware.NewAuthMiddleware(authService)

	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"https://fdz.antasource.xyz",
			"https://be.antasource.xyz",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))
	r.Use(middleware.LogRequest())
	r.Use(gin.Recovery())

	public := r.Group("/api/v1")
	public.Use(middleware.RateLimitPublic(limiter))
	{
		articleHandler := handler.NewArticleHandler(articleService)
		categoryHandler := handler.NewCategoryHandler(categoryService)
		public.POST("/auth/register", handler.NewAuthHandler(authService).Register)
		public.POST("/auth/login", handler.NewAuthHandler(authService).Login)
		public.GET("/articles", articleHandler.PublicList)
		public.GET("/articles/page/:limit/:offset", articleHandler.PublicList)
		public.GET("/articles/:id", articleHandler.GetByID)
		public.GET("/categories", categoryHandler.PublicList)
	}

	admin := r.Group("/api/v1/admin")
	admin.Use(middleware.RateLimitAdmin(limiter))
	admin.Use(authMw.JWTAuth())
	admin.Use(authMw.RequireAuth())
	admin.Use(authMw.RequireAdmin())
	{
		articleHandler := handler.NewArticleHandler(articleService)
		categoryHandler := handler.NewCategoryHandler(categoryService)
		admin.POST("/articles", articleHandler.Create)
		admin.PUT("/articles/:id", articleHandler.Update)
		admin.DELETE("/articles/:id", articleHandler.Delete)
		admin.GET("/articles", articleHandler.List)
		admin.GET("/dashboard", articleHandler.Dashboard)
		admin.GET("/audit-logs", articleHandler.AuditLogs)
		admin.GET("/articles/:id", articleHandler.GetByID)
		admin.POST("/categories", categoryHandler.Create)
		admin.GET("/categories", categoryHandler.List)
		admin.GET("/categories/:id", categoryHandler.GetByID)
		admin.PUT("/categories/:id", categoryHandler.Update)
		admin.DELETE("/categories/:id", categoryHandler.Delete)
	}

	addr := "0.0.0.0:" + config.Conf.App.Port
	config.Log.Info("starting server", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		config.Log.Fatal("server failed", zap.Error(err))
	}
}
