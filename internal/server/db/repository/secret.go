package repository

import (
	"context"
	"database/sql"
	"github.com/zavtra-na-rabotu/GophKeeper/internal/model"
	"go.uber.org/zap"
)

type SecretRepository struct {
	db *sql.DB
}

func NewSecretRepository(db *sql.DB) *SecretRepository {
	return &SecretRepository{db: db}
}

func (r *SecretRepository) Save(ctx context.Context, secret *model.Secret) (int, error) {
	row := r.db.QueryRowContext(
		ctx,
		`INSERT INTO secrets (user_id, title, type, content, metadata, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		secret.UserID,
		secret.Title,
		secret.Type,
		secret.Content,
		secret.Metadata,
		secret.CreatedAt,
		secret.UpdatedAt,
	)

	var secretID int
	err := row.Scan(&secretID)
	if err != nil {
		return 0, err
	}

	return secretID, nil
}

func (r *SecretRepository) GetAllByUserID(ctx context.Context, userID uint64) ([]*model.Secret, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT id, title, type, content, metadata, created_at, updated_at FROM secrets WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var secrets []*model.Secret
	for rows.Next() {
		secret := &model.Secret{}
		err := rows.Scan(&secret.ID, &secret.Title, &secret.Type, &secret.Content, &secret.Metadata, &secret.CreatedAt, &secret.UpdatedAt)
		if err != nil {
			zap.L().Error("Failed to get secrets", zap.Error(err))
			return nil, err
		}
		secrets = append(secrets, secret)
	}

	if err := rows.Err(); err != nil {
		zap.L().Error("Failed to get secrets", zap.Error(err))
		return nil, err
	}

	return secrets, nil
}

func (r *SecretRepository) DeleteByUserID(ctx context.Context, secretID uint64, userID uint64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM secrets WHERE id = $1 AND user_id = $2", secretID, userID)
	if err != nil {
		return err
	}

	return nil
}
