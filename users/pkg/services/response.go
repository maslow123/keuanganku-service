package services

import "github.com/maslow123/users/pkg/pb"

func genericRegisterResponse(statusCode int, errorMessage string) (*pb.RegisterResponse, error) {
	return &pb.RegisterResponse{
		Status: int64(statusCode),
		Error:  errorMessage,
	}, nil
}

func genericLoginResponse(statusCode int, errorMessage string) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{
		Status: int64(statusCode),
		Error:  errorMessage,
	}, nil
}
