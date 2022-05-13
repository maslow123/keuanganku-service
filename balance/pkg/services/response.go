package services

import "github.com/maslow123/balance/pkg/pb"

func genericUpsertBalanceResponse(statusCode int, errorMessage string) (*pb.UpsertBalanceResponse, error) {
	return &pb.UpsertBalanceResponse{
		Status: int32(statusCode),
		Error:  errorMessage,
	}, nil
}

func genericGetUserBalanceResponse(statusCode int, errorMessage string) (*pb.GetUserBalanceResponse, error) {
	return &pb.GetUserBalanceResponse{
		Status: int32(statusCode),
		Error:  errorMessage,
	}, nil
}
