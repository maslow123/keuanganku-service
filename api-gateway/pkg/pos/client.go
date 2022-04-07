package pos

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/config"
	"github.com/maslow123/api-gateway/pkg/pos/pb"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.PosServiceClient
	Router *gin.Engine
}

func InitServiceClient(c *config.Config) pb.PosServiceClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.PosServiceUrl, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewPosServiceClient(cc)
}
