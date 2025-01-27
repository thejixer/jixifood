package main

import (
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	pb "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/services/auth/internal/config"
	"github.com/thejixer/jixifood/services/auth/internal/handlers"
	"google.golang.org/grpc"
)

func init() {
	godotenv.Load()
}

func main() {

	cfg := config.NewAuthServiceConfig()
	lis, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		log.Fatalf("can't create a listener: %s", err)
	}
	s := grpc.NewServer()
	service := handlers.NewAuthServiceServer()
	pb.RegisterAuthServiceServer(s, service)

	fmt.Printf("serving on port : %v \n", cfg.ListenAddr)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to server: %s", err)
	}
}
