package auth

import (
	"fmt"
	"time"

	"github.com/you/sharing-vision-backend-v2/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func ComparePassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func GenerateToken(userID int, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   fmt.Sprint(userID),
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(jwtDuration()).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Conf.JWT.Secret))
}

func jwtDuration() time.Duration {
	d, _ := time.ParseDuration(config.Conf.JWT.Expiry)
	if d == 0 {
		d = 24 * time.Hour
	}
	return d
}
