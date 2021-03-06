package client

import (
	"context"
	"fmt"

	"github.com/maslow123/transactions/pkg/pb"
	"google.golang.org/grpc"
)

type PosServiceClient struct {
	Client pb.PosServiceClient
}

func InitPosServiceClient(url string) PosServiceClient {
	cc, err := grpc.Dial(url, grpc.WithInsecure())

	if err != nil {
		fmt.Println("Could not connect:", err)
	}

	c := PosServiceClient{
		Client: pb.NewPosServiceClient(cc),
	}

	return c
}

func (c *PosServiceClient) PosDetail(posId int32) (*pb.PosDetailResponse, error) {
	req := &pb.PosDetailRequest{
		Id: posId,
	}

	return c.Client.PosDetail(context.Background(), req)
}

func (c *PosServiceClient) UpdateTotalPosByUser(posId, amount int32, action pb.UpdateTotalPosRequest_ActionTransaction) (*pb.UpdateTotalPosResponse, error) {
	req := &pb.UpdateTotalPosRequest{
		Id:     posId,
		Amount: amount,
		Action: action,
	}

	return c.Client.UpdateTotalPosByUser(context.Background(), req)
}
