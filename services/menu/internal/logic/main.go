package logic

import (
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
