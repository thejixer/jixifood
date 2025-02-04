package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	apperrors "github.com/thejixer/jixifood/shared/errors"
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

func (r *MenuRepo) EditCategory(ctx context.Context, category *models.CategoryEntity) (*models.CategoryEntity, error) {
	query := `
		UPDATE categories
		SET name = $2, description = $3, is_quantifiable = $4
		WHERE id = $1
		RETURNING *;
	`
	rows, err := r.db.QueryContext(ctx, query, category.ID, category.Name, category.Description, category.IsQuantifiable)

	if err != nil {
		return nil, fmt.Errorf("error in menuRepo.EditCategory: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		return ScanIntoCategoryEntity(rows)
	}
	return nil, fmt.Errorf("error in menuRepo.EditCategory: %w: %v", apperrors.ErrNotFound, err)
}

func (r *MenuRepo) GetCategories(ctx context.Context) ([]*models.CategoryEntity, error) {
	query := `SELECT * FROM categories`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error in menuRepo.GetCategories: %w", err)
	}
	defer rows.Close()
	var categories []*models.CategoryEntity
	for rows.Next() {
		c, err := ScanIntoCategoryEntity(rows)
		if err != nil {
			return nil, fmt.Errorf("error in menuRepo.GetCategories: %w", err)
		}
		categories = append(categories, c)
	}

	return categories, nil

}
func ScanIntoCategoryEntity(rows *sql.Rows) (*models.CategoryEntity, error) {
	c := new(models.CategoryEntity)
	if err := rows.Scan(
		&c.ID,
		&c.Name,
		&c.Description,
		&c.IsQuantifiable,
		&c.CreatedAt,
	); err != nil {
		return nil, err
	}
	return c, nil
}
