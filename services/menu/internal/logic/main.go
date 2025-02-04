package logic

import (
	"context"
	"fmt"

	authPB "github.com/thejixer/jixifood/generated/auth"
	pb "github.com/thejixer/jixifood/generated/menu"
	apperrors "github.com/thejixer/jixifood/shared/errors"
	"github.com/thejixer/jixifood/shared/models"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	grpcclient "github.com/thejixer/jixifood/services/menu/internal/grpc-client"
	"github.com/thejixer/jixifood/services/menu/internal/repository"
)

type MenuLogic struct {
	dbStore *repository.PostgresStore
	gc      *grpcclient.GrpcClient
}

func NewMenuLogic(dbStore *repository.PostgresStore, gc *grpcclient.GrpcClient) *MenuLogic {
	return &MenuLogic{
		dbStore: dbStore,
		gc:      gc,
	}
}

func (l *MenuLogic) WrapContextAroundNewContext(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, apperrors.ErrMissingMetaData
	}
	return metadata.NewOutgoingContext(ctx, md), nil
}

func (l *MenuLogic) CheckPermission(ctx context.Context, permissionName string) (*authPB.CheckPermissionResponse, error) {

	d := &authPB.CheckPermissionRequest{
		PersmissionName: permissionName,
	}
	resp, err := l.gc.AuthClient.CheckPermission(ctx, d)
	if err != nil {
		return nil, fmt.Errorf("error in MenuLogic.CheckPermission: %w", err)
	}

	return resp, nil
}

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
