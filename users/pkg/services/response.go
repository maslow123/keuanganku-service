package services

import (
	"errors"

	"github.com/maslow123/users/pkg/pb"
)

func genericRegisterResponse(statusCode int, errorMessage string) (*pb.RegisterResponse, error) {
	return &pb.RegisterResponse{
		Status: int32(statusCode),
		Error:  errorMessage,
	}, nil
}

func genericLoginResponse(statusCode int, errorMessage string) (*pb.LoginResponse, error) {
	return &pb.LoginResponse{
		Status: int32(statusCode),
		Error:  errorMessage,
	}, nil
}

func genericUpdateProfileResponse(statusCode int, errorMessage string) (*pb.UpdateProfileResponse, error) {
	return &pb.UpdateProfileResponse{
		Status: int32(statusCode),
		Error:  errorMessage,
	}, nil
}

func genericChangePasswordResponse(statusCode int, errorMessage string) (*pb.ChangePasswordResponse, error) {
	return &pb.ChangePasswordResponse{
		Status: int32(statusCode),
		Error:  errorMessage,
	}, nil
}

func genericUploadImageResponse(errorMessage string) error {
	return errors.New(errorMessage)
}
