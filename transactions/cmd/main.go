package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"os/signal"

	_ "github.com/lib/pq"
	"github.com/maslow123/transactions/pkg/client"
	"github.com/maslow123/transactions/pkg/config"
	"github.com/maslow123/transactions/pkg/pb"
	"github.com/maslow123/transactions/pkg/services"
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

	listen, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalln(err)
	}

	posService := client.InitPosServiceClient(c.PosServiceUrl)

	opts := []grpc.ServerOption{}
	api := services.Server{
		DB:         db,
		PosService: posService,
	}
	server := grpc.NewServer(opts...)
	pb.RegisterTransactionServiceServer(server, &api)

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
