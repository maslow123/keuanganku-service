package pos

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/config"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.BalanceServiceClient
	Router *gin.Engine
}

func InitServiceClient(c *config.Config) pb.BalanceServiceClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.BalanceServiceUrl, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewBalanceServiceClient(cc)
}
