package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/nhassl3/sso/internal/domain/models"
	"github.com/stretchr/testify/assert"
)

func TestNewToken(t *testing.T) {
	testCases := []struct {
		name      string
		user      models.User
		app       models.App
		duration  time.Duration
		expectErr bool
	}{
		{
			name: "Valid token generation",
			user: models.User{
				ID:    3131,
				Email: "user@example.com",
			},
			app: models.App{
				ID:     3131,
				Secret: "mysecret",
			},
			duration:  time.Hour,
			expectErr: false,
		},
		{
			name: "Invalid token generation - empty secret",
			user: models.User{
				ID:    1354,
				Email: "user@example.com",
			},
			app: models.App{
				ID:     1723813,
				Secret: "",
			},
			duration:  time.Hour,
			expectErr: true,
		},
		{
			name: "Invalid token generation - empty duration",
			user: models.User{
				ID:    1354,
				Email: "user@example.com",
			},
			app: models.App{
				ID:     1723813,
				Secret: "adad",
			},
			duration:  0,
			expectErr: true,
		},
		{
			name: "Invalid token generation - empty secret and duration",
			user: models.User{
				ID:    1354,
				Email: "user@example.com",
			},
			app: models.App{
				ID:     1723813,
				Secret: "",
			},
			duration:  0,
			expectErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := NewToken(tt.user, tt.app, tt.duration)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Empty(t, tokenString)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tokenString)

				// Validate token claims
				parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
					return []byte(tt.app.Secret), nil
				})
				assert.NoError(t, err)
				assert.NotNil(t, parsedToken)
				assert.True(t, parsedToken.Valid)

				claims := parsedToken.Claims.(jwt.MapClaims)
				assert.NoError(t, claims.Valid())
				assert.EqualValues(t, tt.user.ID, claims["uid"])
				assert.Equal(t, tt.user.Email, claims["email"])
				assert.EqualValues(t, tt.app.ID, claims["app_id"])

				exp := claims["exp"].(float64)
				assert.WithinDuration(t, time.Unix(int64(exp), 0), time.Now().Add(tt.duration), time.Minute)
			}
		})
	}
}
