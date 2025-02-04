package handlers

import (
	"context"

	pb "github.com/thejixer/jixifood/generated/menu"
	"github.com/thejixer/jixifood/shared/constants"
	apperrors "github.com/thejixer/jixifood/shared/errors"
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
