package services

import (
	"context"
	"database/sql"
	"log"
	"net"
	"testing"

	_ "github.com/lib/pq"
	"github.com/maslow123/transactions/pkg/client"
	"github.com/maslow123/transactions/pkg/config"
	"github.com/maslow123/transactions/pkg/pb"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func dialer(t *testing.T) func(context.Context, string) (net.Conn, error) {
	c, err := config.LoadConfig("../config/envs", "test")
	require.NoError(t, err)

	listener := bufconn.Listen(1024 * 1024)

	posService := client.InitPosServiceClient(c.PosServiceUrl)

	db, err := sql.Open("postgres", c.DBUrl)
	if err != nil {
		require.NoError(t, err)
	}

	s := Server{
		DB:         db,
		PosService: posService,
	}

	server := grpc.NewServer()
	pb.RegisterTransactionServiceServer(server, &s)
	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func checkConnection(ctx context.Context, t *testing.T) *grpc.ClientConn {
	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(t)))
	require.NoError(t, err)

	return conn
}
