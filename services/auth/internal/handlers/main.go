package handlers

import pb "github.com/thejixer/jixifood/generated/auth"

type authServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

func NewAuthServiceServer() *authServiceServer {
	return &authServiceServer{}
}
