package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"github.com/nhassl3/sso/internal/domain/models"
	"github.com/nhassl3/sso/internal/storage"
)

const (
	opNew      = "storage.sqlite.New"
	opSaveUser = "storage.sqlite.SaveUser"
	opUser     = "storage.sqlite.User"
	opIsAdmin  = "storage.sqlite.IsAdmin"
	opApp      = "storage.sqlite.App"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", opNew, err)
	}

	return &Storage{db: db}, nil
}

// SaveUser save user in database with given credentials
//
// returns user ID if function successfully complete. Type: int64
// else return error
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", opSaveUser, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", opSaveUser, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", opSaveUser, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", opSaveUser, err)
	}

	return id, nil
}

// User returns user by email
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	var user models.User

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", opUser, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	if err = row.Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, storage.ErrUserNotFound
		}

		return models.User{}, fmt.Errorf("%s: %w", opUser, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	var isAdmin bool

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: %w", opIsAdmin, err)
	}

	res, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, storage.ErrUserNotFound
		}

		return false, fmt.Errorf("%s: %w", opIsAdmin, err)
	}

	if err := res.Scan(&isAdmin); err != nil {
		return false, fmt.Errorf("%s: %w", opIsAdmin, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	var app models.App

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps where id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", opApp, err)
	}

	row := stmt.QueryRowContext(ctx, appID)
	if err = row.Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, storage.ErrUserNotFound
		}

		return models.App{}, fmt.Errorf("%s: %w", opApp, err)
	}

	return app, nil
}
