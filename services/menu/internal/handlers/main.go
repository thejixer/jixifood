package handlers

import (
	"context"

	authPB "github.com/thejixer/jixifood/generated/auth"
	pb "github.com/thejixer/jixifood/generated/menu"
	"github.com/thejixer/jixifood/services/menu/internal/logic"
	"github.com/thejixer/jixifood/shared/models"
)

type MenuLogicInterface interface {
	CheckPermission(ctx context.Context, permissionName string) (*authPB.CheckPermissionResponse, error)
	CreateCategory(ctx context.Context, name, description string, isQuantifiable bool) (*models.CategoryEntity, error)
	MapCategoryEntityToPBCategory(*models.CategoryEntity) *pb.Category
	MapCategoriesToPB(categories []*models.CategoryEntity) []*pb.Category
	WrapContextAroundNewContext(ctx context.Context) (context.Context, error)
	EditCategory(ctx context.Context, category *models.CategoryEntity) (*models.CategoryEntity, error)
	GetCategories(ctx context.Context) ([]*models.CategoryEntity, error)
	GetCategory(ctx context.Context, id uint64) (*models.CategoryEntity, error)
}

type MenuHandler struct {
	pb.UnimplementedMenuServiceServer
	MenuLogic MenuLogicInterface
}

func NewMenuHandler(menuLogic *logic.MenuLogic) *MenuHandler {
	return &MenuHandler{
		MenuLogic: menuLogic,
	}
}
