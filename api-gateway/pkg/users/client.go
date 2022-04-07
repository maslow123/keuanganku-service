package users

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/config"
	"github.com/maslow123/api-gateway/pkg/users/pb"
	"google.golang.org/grpc"
)

type ServiceClient struct {
	Client pb.UserServiceClient
	Router *gin.Engine
}

func InitServiceClient(c *config.Config) pb.UserServiceClient {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.UserServiceUrl, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	return pb.NewUserServiceClient(cc)
}
