package client

import (
	"context"
	"fmt"
	"log"

	"github.com/maslow123/transactions/pkg/pb"
	"google.golang.org/grpc"
)

type BalanceServiceClient struct {
	Client pb.BalanceServiceClient
}

func InitBalanceServiceClient(url string) BalanceServiceClient {
	cc, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := BalanceServiceClient{
		Client: pb.NewBalanceServiceClient(cc),
	}

	return c
}

func (c *BalanceServiceClient) UpsertBalance(userId, balanceType, action, total int32) (*pb.UpsertBalanceResponse, error) {
	actionType := pb.UpsertBalanceRequest_ActionType(pb.UpsertBalanceRequest_ActionType_value["INCREASE"])
	if action == 1 {
		actionType = pb.UpsertBalanceRequest_ActionType(pb.UpsertBalanceRequest_ActionType_value["DECREASE"])
	}

	req := &pb.UpsertBalanceRequest{
		UserId: userId,
		Type:   balanceType,
		Action: actionType,
		Total:  total,
	}

	log.Println("Request: ", req)

	return c.Client.UpsertBalance(context.Background(), req)
}
