package logic

import (
	"context"
	"fmt"

	pb "github.com/thejixer/jixifood/generated/menu"
	"github.com/thejixer/jixifood/shared/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (l *MenuLogic) CreateCategory(ctx context.Context, name, description string, isQuantifiable bool) (*models.CategoryEntity, error) {

	c := &models.CategoryEntity{
		Name:           name,
		Description:    description,
		IsQuantifiable: isQuantifiable,
	}

	category, err := l.dbStore.MenuRepo.CreateCategory(ctx, c)
	if err != nil {
		return nil, fmt.Errorf("error in MenuLogic.CreateCategory: %w", err)
	}

	return category, nil
}

func (l *MenuLogic) MapCategoryEntityToPBCategory(c *models.CategoryEntity) *pb.Category {

	return &pb.Category{
		Id:             c.ID,
		Name:           c.Name,
		Description:    c.Description,
		IsQuantifiable: c.IsQuantifiable,
		CreatedAt:      timestamppb.New(c.CreatedAt),
	}
}

func (l *MenuLogic) MapCategoriesToPB(categories []*models.CategoryEntity) []*pb.Category {
	var result []*pb.Category
	for _, item := range categories {
		c := l.MapCategoryEntityToPBCategory(item)
		result = append(result, c)
	}
	return result
}

func (l *MenuLogic) EditCategory(ctx context.Context, category *models.CategoryEntity) (*models.CategoryEntity, error) {

	c, err := l.dbStore.MenuRepo.EditCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("error in MenuLogic.EditCategory: %w", err)
	}

	return c, nil
}

func (l *MenuLogic) GetCategories(ctx context.Context) ([]*models.CategoryEntity, error) {
	categories, err := l.dbStore.MenuRepo.GetCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in MenuLogic.GetCategories: %w", err)
	}
	return categories, nil
}
