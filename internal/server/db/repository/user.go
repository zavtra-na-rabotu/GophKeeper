package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/model"
	"go.uber.org/zap"
)

var (
	ErrUserAlreadyExists = errors.New("user already exist")
	ErrUserNotFound      = errors.New("user not found")
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*model.User, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, login, password_hash FROM users WHERE login = $1`, login)

	var user model.User
	err := row.Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		zap.L().Error("Failed to query user by login", zap.String("login", login), zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, login string, password string) (uint64, error) {
	row := r.db.QueryRowContext(ctx, `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id`, login, password)

	var userID uint64
	err := row.Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, ErrUserAlreadyExists
		}
		return 0, err
	}

	return userID, nil
}
