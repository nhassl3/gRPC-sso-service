package jwt

import (
	"time"

	JWT "github.com/golang-jwt/jwt"
	"github.com/nhassl3/sso/internal/domain/models"
)

// NewToken generate JWToken that let user to get some actions in some services
func NewToken(user models.User, app models.App, duration time.Duration) (string, error) {
	token := JWT.New(JWT.SigningMethodHS256)

	claims := token.Claims.(JWT.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()
	claims["app_id"] = app.ID

	tokenString, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
