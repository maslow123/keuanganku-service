package services

import "github.com/maslow123/pos/pkg/pb"

func genericCreatePosResponse(statusCode int, errorMessage string) (*pb.CreatePosResponse, error) {
	return &pb.CreatePosResponse{
		Status: int64(statusCode),
		Error:  errorMessage,
	}, nil
}

func genericPosDetailResponse(statusCode int, errorMessage string) (*pb.PosDetailResponse, error) {
	return &pb.PosDetailResponse{
		Status: int64(statusCode),
		Error:  errorMessage,
	}, nil
}

func genericListPosByUserResponse(statusCode int, errorMessage string) (*pb.GetPosListResponse, error) {
	return &pb.GetPosListResponse{
		Status: int64(statusCode),
		Error:  errorMessage,
	}, nil
}

func genericUpdatePosByUserResponse(statusCode int, errorMessage string) (*pb.UpdatePosResponse, error) {
	return &pb.UpdatePosResponse{
		Status: int64(statusCode),
		Error:  errorMessage,
	}, nil
}

func genericDeletePosByUserResponse(statusCode int, errorMessage string) (*pb.DeletePosResponse, error) {
	return &pb.DeletePosResponse{
		Status: int64(statusCode),
		Error:  errorMessage,
	}, nil
}
