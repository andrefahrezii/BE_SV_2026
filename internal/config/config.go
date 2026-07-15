package config

import (
	"log"
	"os"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Name string
		Port string
		Env  string
	}
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
	JWT struct {
		Secret string
		Expiry string
	}
	RateLimit struct {
		PublicRPS   int
		PublicBurst int
		AdminRPS    int
		AdminBurst  int
	}
}

var (
	Conf *Config
	Log  *zap.Logger
)

func Load() {
	// Load .env file first
	_ = godotenv.Load()

	Conf = &Config{}

	Conf.App.Name = getEnv("APP_NAME", "sv-backend")
	Conf.App.Port = getEnv("APP_PORT", "8081")
	Conf.App.Env = getEnv("APP_ENV", "development")

	Conf.DB.Host = getEnv("DB_HOST", "localhost")
	Conf.DB.Port = getEnv("DB_PORT", "5432")
	Conf.DB.User = getEnv("DB_USER", "sv_user")
	Conf.DB.Password = getEnv("DB_PASSWORD", "sv_password")
	Conf.DB.Name = getEnv("DB_NAME", "sv_portal")
	Conf.DB.SSLMode = getEnv("DB_SSLMODE", "require")

	Conf.Redis.Addr = getEnv("REDIS_ADDR", "localhost:6379")
	Conf.Redis.Password = getEnv("REDIS_PASSWORD", "")
	Conf.Redis.DB = getEnvInt("REDIS_DB", 0)

	Conf.JWT.Secret = getEnv("JWT_SECRET", "change-me")
	Conf.JWT.Expiry = getEnv("JWT_EXPIRY", "24h")

	Conf.RateLimit.PublicRPS = getEnvInt("RATE_LIMIT_PUBLIC_RPS", 100)
	Conf.RateLimit.PublicBurst = getEnvInt("RATE_LIMIT_PUBLIC_BURST", 10)
	Conf.RateLimit.AdminRPS = getEnvInt("RATE_LIMIT_ADMIN_RPS", 300)
	Conf.RateLimit.AdminBurst = getEnvInt("RATE_LIMIT_ADMIN_BURST", 30)

	initLogger()
}

func initLogger() {
	level := zap.InfoLevel
	if Conf.App.Env == "production" {
		level = zap.WarnLevel
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       Conf.App.Env == "development",
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	var err error
	Log, err = cfg.Build()
	if err != nil {
		log.Fatalf("logger init failed: %v", err)
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func getEnvInt(k string, def int) int {
	if v := os.Getenv(k); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return def
}
