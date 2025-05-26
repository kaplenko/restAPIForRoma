package jwt

import (
	"os"
	"time"
	"user-service/internal/entity"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func Init(JWTSecret string) {
	jwtSecret = []byte(os.Getenv(JWTSecret))

}

func NewToken(user entity.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
