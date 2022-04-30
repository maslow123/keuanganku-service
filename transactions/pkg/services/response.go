package services

import "github.com/maslow123/transactions/pkg/pb"

func genericCreateTransactionResponse(statusCode int, errorMessage string) (*pb.CreateTransactionResponse, error) {
	return &pb.CreateTransactionResponse{
		Status: int32(statusCode),
		Error:  errorMessage,
	}, nil
}

func genericGetTransactionListByUserResponse(statusCode int, errorMessage string) (*pb.GetTransactionListResponse, error) {
	return &pb.GetTransactionListResponse{
		Status: int32(statusCode),
		Error:  errorMessage,
	}, nil
}
