package logic

import (
	"context"
	"fmt"

	authPB "github.com/thejixer/jixifood/generated/auth"
	apperrors "github.com/thejixer/jixifood/shared/errors"
	"google.golang.org/grpc/metadata"

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
