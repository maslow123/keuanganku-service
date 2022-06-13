package services

import (
	"context"
	"database/sql"
	"log"
	"net"
	"testing"

	_ "github.com/lib/pq"
	"github.com/maslow123/users/pkg/config"
	"github.com/maslow123/users/pkg/pb"
	"github.com/maslow123/users/pkg/utils"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/maslow123/users/pkg/client"
)

func dialer(t *testing.T) func(context.Context, string) (net.Conn, error) {
	c, err := config.LoadConfig("../config/envs", "test")
	require.NoError(t, err)

	listener := bufconn.Listen(1024 * 1024)

	jwt := utils.JwtWrapper{
		SecretKey:       c.JWTSecretKey,
		Issuer:          "user-service",
		ExpirationHours: 24 * 365,
	}

	db, err := sql.Open("postgres", c.DBUrl)
	if err != nil {
		require.NoError(t, err)
	}

	balanceService := client.InitBalanceServiceClient(c.BalanceServiceUrl)
	imageStore := NewDiskImageStore("../tmp")
	s := Server{
		DB:             db,
		Jwt:            jwt,
		BalanceService: balanceService,
		ImageStore:     imageStore,
	}

	server := grpc.NewServer()
	pb.RegisterUserServiceServer(server, &s)
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
