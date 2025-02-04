package models

import (
	"context"
	"time"
)

type MenuRepository interface {
	CreateCategory(ctx context.Context, category *CategoryEntity) (*CategoryEntity, error)
	EditCategory(ctx context.Context, category *CategoryEntity) (*CategoryEntity, error)
	GetCategories(ctx context.Context) ([]*CategoryEntity, error)
	GetCategory(ctx context.Context, id uint64) (*CategoryEntity, error)
}

type CategoryEntity struct {
	ID             uint64
	Name           string
	Description    string
	IsQuantifiable bool
	CreatedAt      time.Time
}

// ######### DTOS #########

type CreateCategoryDto struct {
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description" validate:"required"`
	IsQuantifiable bool   `json:"is_quantifiable" validate:"omitempty"`
}

type EditCategoryDto struct {
	Name           string `json:"name" validate:"required"`
	Description    string `json:"description" validate:"required"`
	IsQuantifiable bool   `json:"is_quantifiable" validate:"omitempty"`
}

type CategoryDto struct {
	ID             uint64    `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	IsQuantifiable bool      `json:"is_quantifiable"`
	CreatedAt      time.Time `json:"created_at"`
}
