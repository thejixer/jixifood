package handlers

import (
	pb "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/services/auth/internal/repository"
)

type authServiceServer struct {
	pb.UnimplementedAuthServiceServer
	dbStore *repository.PostgresStore
}

func NewAuthServiceServer(dbStore *repository.PostgresStore) *authServiceServer {
	return &authServiceServer{
		dbStore: dbStore,
	}
}
