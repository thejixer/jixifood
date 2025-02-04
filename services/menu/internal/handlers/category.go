package handlers

import (
	"context"

	pb "github.com/thejixer/jixifood/generated/menu"
	"github.com/thejixer/jixifood/shared/constants"
	apperrors "github.com/thejixer/jixifood/shared/errors"
	"github.com/thejixer/jixifood/shared/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *MenuHandler) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.Category, error) {

	if req.Name == "" || req.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "bad request : "+apperrors.ErrInputRequirements.Error())
	}

	NewCtx, err := s.MenuLogic.WrapContextAroundNewContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, apperrors.ErrUnauthorized.Error())
	}

	resp, err := s.MenuLogic.CheckPermission(NewCtx, constants.PermissionManageMenu)
	if err != nil || !resp.HasPermission {
		return nil, status.Error(codes.PermissionDenied, apperrors.ErrForbidden.Error())
	}

	c, err := s.MenuLogic.CreateCategory(NewCtx, req.Name, req.Description, req.IsQuantifiable)
	if err != nil {
		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())
	}

	category := s.MenuLogic.MapCategoryEntityToPBCategory(c)
	return category, nil

}

func (s *MenuHandler) EditCategory(ctx context.Context, req *pb.EditCategoryRequest) (*pb.Category, error) {
	if req.Name == "" || req.Description == "" {
		return nil, status.Error(codes.InvalidArgument, "bad request : "+apperrors.ErrInputRequirements.Error())
	}

	NewCtx, err := s.MenuLogic.WrapContextAroundNewContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, apperrors.ErrUnauthorized.Error())
	}
	resp, err := s.MenuLogic.CheckPermission(NewCtx, constants.PermissionManageMenu)
	if err != nil || !resp.HasPermission {
		return nil, status.Error(codes.PermissionDenied, apperrors.ErrForbidden.Error())
	}

	d := &models.CategoryEntity{
		ID:             req.Id,
		Name:           req.Name,
		Description:    req.Description,
		IsQuantifiable: req.IsQuantifiable,
	}
	c, err := s.MenuLogic.EditCategory(NewCtx, d)
	if err != nil {
		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())
	}
	category := s.MenuLogic.MapCategoryEntityToPBCategory(c)
	return category, nil
}

func (s *MenuHandler) GetCategories(ctx context.Context, req *pb.Empty) (*pb.GetCategoriesResponse, error) {
	resp, err := s.MenuLogic.GetCategories(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, apperrors.ErrUnexpected.Error())
	}
	categories := s.MenuLogic.MapCategoriesToPB(resp)
	return &pb.GetCategoriesResponse{
		Categories: categories,
	}, nil
}

func (s *MenuHandler) GetCategory(ctx context.Context, req *pb.GetCategoryRequest) (*pb.Category, error) {
	c, err := s.MenuLogic.GetCategory(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, apperrors.ErrNotFound.Error())
	}
	category := s.MenuLogic.MapCategoryEntityToPBCategory(c)
	return category, nil
}
