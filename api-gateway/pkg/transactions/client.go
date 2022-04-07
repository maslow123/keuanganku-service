package transactions

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/config"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.TransactionServiceClient
	Router *gin.Engine
}

func InitServiceClient(c *config.Config) pb.TransactionServiceClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.TransactionServiceUrl, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewTransactionServiceClient(cc)
}
