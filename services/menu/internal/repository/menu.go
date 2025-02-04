package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/thejixer/jixifood/shared/models"
)

func (s *PostgresStore) createMenuTables() error {

	query := `
		CREATE TABLE if not exists categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL UNIQUE,
			description TEXT,
			is_quantifiable BOOLEAN NOT NULL DEFAULT FALSE,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := s.db.Exec(query)

	return err

}

type MenuRepo struct {
	db *sql.DB
}

func NewMenuRepo(db *sql.DB) *MenuRepo {
	return &MenuRepo{
		db: db,
	}
}

func (r *MenuRepo) CreateCategory(ctx context.Context, category *models.CategoryEntity) (*models.CategoryEntity, error) {

	NewCategory := &models.CategoryEntity{
		Name:           category.Name,
		Description:    category.Description,
		IsQuantifiable: category.IsQuantifiable,
		CreatedAt:      time.Now().UTC(),
	}

	query := `
	INSERT INTO categories (name, description, is_quantifiable, created_at)
	VALUES ($1, $2, $3, $4) RETURNING id`
	lastInsertId := 0
	insertErr := r.db.QueryRowContext(
		ctx,
		query,
		NewCategory.Name,
		NewCategory.Description,
		NewCategory.IsQuantifiable,
		NewCategory.CreatedAt,
	).Scan(&lastInsertId)

	if insertErr != nil {
		return nil, fmt.Errorf("error in menuRepo.createCategory: %w", insertErr)
	}
	NewCategory.ID = uint64(lastInsertId)
	return NewCategory, nil

}
