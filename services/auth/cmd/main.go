package main

import (
	"fmt"
	"log"
	"net"

	"github.com/joho/godotenv"
	pb "github.com/thejixer/jixifood/generated/auth"
	"github.com/thejixer/jixifood/services/auth/internal/config"
	"github.com/thejixer/jixifood/services/auth/internal/handlers"
	"github.com/thejixer/jixifood/services/auth/internal/redis"
	"github.com/thejixer/jixifood/services/auth/internal/repository"
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

	dbStore, err := repository.NewPostgresStore(cfg)
	if err != nil {
		log.Fatal("could not connect to the database: ", err)
	}

	if err := dbStore.Init(); err != nil {
		log.Fatal("could not connect to the database: ", err)
	}

	redisStore, err := redis.NewRedisStore(cfg)
	if err != nil {
		log.Fatalf("can not connect to redis, %s", err)
	}

	s := grpc.NewServer()

	service := handlers.NewAuthServiceServer(dbStore, redisStore)
	pb.RegisterAuthServiceServer(s, service)

	fmt.Printf("serving on port : %v \n", cfg.ListenAddr)
	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("impossible to server: %s", err)
	}
}
