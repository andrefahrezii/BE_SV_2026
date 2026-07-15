package middleware

import (
	"strings"

	"github.com/you/sharing-vision-backend-v2/internal/config"
	"github.com/you/sharing-vision-backend-v2/internal/model"
	"github.com/you/sharing-vision-backend-v2/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (a *AuthMiddleware) JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.Conf.JWT.Secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			email := jwtToString(claims["email"])
			role := jwtToString(claims["role"])
			sub := jwtToString(claims["sub"])
			userID := 0
			for _, ch := range sub {
				if ch >= '0' && ch <= '9' {
					userID = userID*10 + int(ch-'0')
				}
			}
			c.Set("user", &model.User{
				ID:    userID,
				Email: email,
				Role:  role,
			})
		}
		c.Next()
	}
}

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := getUser(c)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		c.Next()
	}
}

func (a *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, ok := getUser(c)
		if !ok || user.Role != "admin" {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}

func getUser(c *gin.Context) (*model.User, bool) {
	v, ok := c.Get("user")
	if !ok {
		return nil, false
	}
	user, ok := v.(*model.User)
	return user, ok
}

func jwtToString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
