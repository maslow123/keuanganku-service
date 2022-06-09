package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
	"github.com/maslow123/users/pkg/client"
	"github.com/maslow123/users/pkg/config"
	"github.com/maslow123/users/pkg/pb"
	"github.com/maslow123/users/pkg/services"
	"github.com/maslow123/users/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	c, err := config.LoadConfig("./pkg/config/envs", "dev")

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	db, err := sql.Open("postgres", c.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "user-service",
		ExpirationHours: 24 * 1,
	}

	listen, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalln(err)
	}

	balanceService := client.InitBalanceServiceClient(c.BalanceServiceUrl)
	opts := []grpc.ServerOption{}
	api := services.Server{
		DB:             db,
		Jwt:            jwt,
		BalanceService: balanceService,
	}
	server := grpc.NewServer(opts...)
	pb.RegisterUserServiceServer(server, &api)

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	ctx := context.Background()

	go func() {
		for range channel {
			log.Println("Shutting down gRPC server...")
			server.GracefulStop()
			<-ctx.Done()
		}
	}()

	log.Println("Starting gRPC server on port: ", c.Port)
	server.Serve(listen)
}
