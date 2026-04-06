package repository

import (
	"context"
	"fmt"

	"goledger-challenge-besu/internal/database"
)

type StorageRepository interface {
	GetValue(ctx context.Context) (string, error)
	SetValue(ctx context.Context, value string) error
}

type storageRepository struct {
	db *database.Postgres
}

func NewStorageRepository(db *database.Postgres) StorageRepository {
	return &storageRepository{db: db}
}

func (r *storageRepository) GetValue(ctx context.Context) (string, error) {
	var value string
	err := r.db.Pool.QueryRow(ctx, "SELECT value FROM StorageState WHERE id = 1").Scan(&value)
	if err != nil {
		return "", fmt.Errorf("failed to get value from database: %w", err)
	}
	return value, nil
}

func (r *storageRepository) SetValue(ctx context.Context, value string) error {
	_, err := r.db.Pool.Exec(ctx,
		"UPDATE StorageState SET value = $1, updated_at = NOW() WHERE id = 1",
		value,
	)
	if err != nil {
		return fmt.Errorf("failed to set value in database: %w", err)
	}
	return nil
}
