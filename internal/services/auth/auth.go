package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/nhassl3/sso/internal/domain/models"
	"github.com/nhassl3/sso/internal/lib/jwt"
	sl "github.com/nhassl3/sso/internal/lib/logger/sl"
	"github.com/nhassl3/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

const (
	opRegisterUser = "auth.RegisterNewUser"
	opLogin        = "auth.Login"
	opIsAmdin      = "auth.IsAdmin"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrGenerateJWToken    = errors.New("token generation error")
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, password []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (user models.User, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

type AppProvider interface {
	App(ctx context.Context, appD int) (app models.App, err error)
}

// New returns a new instance of the Auth service
func New(
	log *slog.Logger,
	usrSaver UserSaver,
	usrProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    usrSaver,
		usrProvider: usrProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

// Login checks if user with given credentials exists in the system
//
// If user exists with given email, but password is incorrect, returns error
// If user doesn't exists, returns error
// Else returns string value (token) and nil for error object
func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, error) {
	log := a.log.With(
		slog.String("op", opLogin),
		slog.String("email", email),
		slog.Int("AppID", appID),
	)

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", sl.ErrLog(err))

			return "", fmt.Errorf("%s: %w", opLogin, ErrInvalidCredentials)
		}
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Warn("invalid credentials", sl.ErrLog(err))

		return "", fmt.Errorf("%s: %w", opLogin, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", opLogin, err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", opLogin, err)
	}

	return token, nil
}

// RegisterNewUser lets user register in system with given credentials
func (a *Auth) RegisterNewUser(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	log := a.log.With(
		slog.String("op", opRegisterUser),
		slog.String("email", email),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.ErrLog(err))

		return 0, fmt.Errorf("%s: %w", opRegisterUser, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user in database", sl.ErrLog(err))

		return 0, fmt.Errorf("%s: %w", opRegisterUser, err)
	}

	log.Info("User registered", slog.Int64("id", id))

	return id, nil
}

// IsAdmin checks user if is have a admin permissions'
//
// If it's ordinary user returns false value
// Else true value for admin
func (a *Auth) IsAdmin(
	ctx context.Context,
	userID int64,
) (bool, error) {
	// TODO: implement
	return false, nil
}
