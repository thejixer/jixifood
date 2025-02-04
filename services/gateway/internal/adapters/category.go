package adapters

import (
	menuPB "github.com/thejixer/jixifood/generated/menu"
	"github.com/thejixer/jixifood/shared/models"
)

func MapPBCategoryToCategoryDTO(c *menuPB.Category) *models.CategoryDto {

	return &models.CategoryDto{
		ID:             c.Id,
		Name:           c.Name,
		Description:    c.Description,
		IsQuantifiable: c.IsQuantifiable,
		CreatedAt:      c.CreatedAt.AsTime(),
	}
}

func MapPBtoCategoriesDto(categories []*menuPB.Category) []*models.CategoryDto {
	var result []*models.CategoryDto
	for _, item := range categories {
		c := MapPBCategoryToCategoryDTO(item)
		result = append(result, c)
	}
	return result
}
