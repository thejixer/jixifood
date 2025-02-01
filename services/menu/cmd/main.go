package main

import (
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	pb "github.com/thejixer/jixifood/generated/menu"
	"github.com/thejixer/jixifood/services/menu/internal/config"
	grpcclient "github.com/thejixer/jixifood/services/menu/internal/grpc-client"
	"github.com/thejixer/jixifood/services/menu/internal/handlers"
	"github.com/thejixer/jixifood/services/menu/internal/logic"
	"github.com/thejixer/jixifood/services/menu/internal/repository"
	"google.golang.org/grpc"
)

func init() {
	godotenv.Load()
}

func main() {
	cfg := config.NewMenuServiceConfig()
	lis, err := net.Listen("tcp", cfg.ListenAddr)
	if err != nil {
		log.Fatalf("can't create a listener: %s", err)
	}

	dbStore, err := repository.NewPostgresStore(cfg)
	if err != nil {
		log.Fatal("could not connect to the database: ", err)
	}

	if err := dbStore.Init(); err != nil {
		log.Fatal("could not connect to the database: ", err)
	}

	gc := grpcclient.NewGRPCClient(cfg)
	defer gc.Shutdown()

	MenuLogic := logic.NewMenuLogic(dbStore, gc)

	menuHandler := handlers.NewMenuHandler(MenuLogic)
	s := grpc.NewServer()
	pb.RegisterMenuServiceServer(s, menuHandler)

	fmt.Printf("serving on port : %v \n", cfg.ListenAddr)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to server: %s", err)
	}
}
