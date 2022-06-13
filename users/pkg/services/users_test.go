package services

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/maslow123/users/pkg/pb"
	"github.com/maslow123/users/pkg/utils"
	"github.com/stretchr/testify/require"
)

var randUser, randPass string

func TestRegister(t *testing.T) {
	randUser = utils.RandomString(10)
	randPass = utils.RandomString(10)

	testCases := []struct {
		name string
		req  *pb.RegisterRequest
		resp *pb.RegisterResponse
	}{
		{
			"OK",
			&pb.RegisterRequest{
				Name:            utils.RandomString(10),
				Email:           randUser,
				Password:        randPass,
				ConfirmPassword: randPass,
			},
			&pb.RegisterResponse{
				Status: int32(http.StatusCreated),
				Error:  "",
			},
		},
		{
			"Invalid Name",
			&pb.RegisterRequest{
				Name:            "",
				Email:           randUser,
				Password:        randPass,
				ConfirmPassword: randPass,
			},
			&pb.RegisterResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-name",
			},
		},
		{
			"Invalid Email",
			&pb.RegisterRequest{
				Name:            utils.RandomString(10),
				Email:           "",
				Password:        randPass,
				ConfirmPassword: randPass,
			},
			&pb.RegisterResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-email",
			},
		},
		{
			"Invalid Password",
			&pb.RegisterRequest{
				Name:            utils.RandomString(10),
				Email:           randUser,
				Password:        "",
				ConfirmPassword: randPass,
			},
			&pb.RegisterResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-password",
			},
		},
		{
			"Invalid Confirm Password",
			&pb.RegisterRequest{
				Name:            utils.RandomString(10),
				Email:           randUser,
				Password:        randPass,
				ConfirmPassword: "",
			},
			&pb.RegisterResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-confirm-password",
			},
		},
		{
			"Password Not Match",
			&pb.RegisterRequest{
				Name:            utils.RandomString(10),
				Email:           randUser,
				Password:        randPass,
				ConfirmPassword: "not match",
			},
			&pb.RegisterResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "password-not-match",
			},
		},
		{
			"Existing user",
			&pb.RegisterRequest{
				Name:            utils.RandomString(10),
				Email:           randUser,
				Password:        randPass,
				ConfirmPassword: randPass,
			},
			&pb.RegisterResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "email-already-exists",
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			response, err := client.Register(ctx, tc.req)
			log.Println(response)
			require.NoError(t, err)
			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)
		})
	}
}

func TestLogin(t *testing.T) {
	randUser = "user1@gmail.com"
	randPass = "111111"

	testCases := []struct {
		name string
		req  *pb.LoginRequest
		resp *pb.LoginResponse
	}{
		{
			"OK",
			&pb.LoginRequest{
				Email:    randUser,
				Password: randPass,
			},
			&pb.LoginResponse{
				Status: int32(http.StatusOK),
				Error:  "",
			},
		},
		{
			"Invalid Email",
			&pb.LoginRequest{
				Email:    "",
				Password: randPass,
			},
			&pb.LoginResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-email",
			},
		},
		{
			"Invalid Password",
			&pb.LoginRequest{
				Email:    randUser,
				Password: "",
			},
			&pb.LoginResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-password",
			},
		},
		{
			"Wrong Password",
			&pb.LoginRequest{
				Email:    randUser,
				Password: "wrong password",
			},
			&pb.LoginResponse{
				Status: int32(http.StatusUnauthorized),
				Error:  "password-not-match",
			},
		},
		{
			"User Not Found",
			&pb.LoginRequest{
				Email:    "xxxx",
				Password: "xxxx",
			},
			&pb.LoginResponse{
				Status: int32(http.StatusNotFound),
				Error:  "user-not-found",
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			response, err := client.Login(ctx, tc.req)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)

			if response.Status == int32(http.StatusOK) {
				require.NotEmpty(t, response.Token)
				require.NotNil(t, response.User)
			}
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	testCases := []struct {
		name string
		req  *pb.UpdateProfileRequest
		resp *pb.UpdateProfileResponse
	}{
		{
			"OK",
			&pb.UpdateProfileRequest{
				Id:    2,
				Name:  "User Updated",
				Email: "user2@gmail.com",
			},
			&pb.UpdateProfileResponse{
				Status: http.StatusOK,
				Error:  "",
			},
		},
		{
			"Invalid User ID",
			&pb.UpdateProfileRequest{
				Id:    0,
				Name:  "User Updated",
				Email: "user2@gmail.com",
			},
			&pb.UpdateProfileResponse{
				Status: http.StatusBadRequest,
				Error:  "invalid-user-id",
			},
		},
		{
			"Invalid Name",
			&pb.UpdateProfileRequest{
				Id:    1,
				Name:  "",
				Email: "user2@gmail.com",
			},
			&pb.UpdateProfileResponse{
				Status: http.StatusBadRequest,
				Error:  "invalid-name",
			},
		},
		{
			"Invalid Email",
			&pb.UpdateProfileRequest{
				Id:    1,
				Name:  "User Updated",
				Email: "",
			},
			&pb.UpdateProfileResponse{
				Status: http.StatusBadRequest,
				Error:  "invalid-email",
			},
		},
		{
			"Invalid User Not Found",
			&pb.UpdateProfileRequest{
				Id:    9999,
				Name:  "User Updated",
				Email: "user2@gmail.com",
			},
			&pb.UpdateProfileResponse{
				Status: http.StatusNotFound,
				Error:  "user-not-found",
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			response, err := client.UpdateProfile(ctx, tc.req)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)
		})
	}
}

func TestChangePassword(t *testing.T) {
	testCases := []struct {
		name string
		req  *pb.ChangePasswordRequest
		resp *pb.ChangePasswordResponse
	}{
		{
			"OK",
			&pb.ChangePasswordRequest{
				Id:              1,
				OldPassword:     "111111",
				Password:        "111111",
				ConfirmPassword: "111111",
			},
			&pb.ChangePasswordResponse{
				Status: http.StatusOK,
				Error:  "",
			},
		},
		{
			"Invalid User ID",
			&pb.ChangePasswordRequest{
				Id:              0,
				OldPassword:     "111111",
				Password:        "111111",
				ConfirmPassword: "111111",
			},
			&pb.ChangePasswordResponse{
				Status: http.StatusBadRequest,
				Error:  "invalid-user-id",
			},
		},
		{
			"Invalid Old Password",
			&pb.ChangePasswordRequest{
				Id:              1,
				OldPassword:     "",
				Password:        "111111",
				ConfirmPassword: "111111",
			},
			&pb.ChangePasswordResponse{
				Status: http.StatusBadRequest,
				Error:  "invalid-old-password",
			},
		},
		{
			"Invalid Password",
			&pb.ChangePasswordRequest{
				Id:              1,
				OldPassword:     "111111",
				Password:        "",
				ConfirmPassword: "111111",
			},
			&pb.ChangePasswordResponse{
				Status: http.StatusBadRequest,
				Error:  "invalid-password",
			},
		},
		{
			"Invalid Confirm Password",
			&pb.ChangePasswordRequest{
				Id:              1,
				OldPassword:     "111111",
				Password:        "111111",
				ConfirmPassword: "",
			},
			&pb.ChangePasswordResponse{
				Status: http.StatusBadRequest,
				Error:  "invalid-confirm-password",
			},
		},
		{
			"Password doesnt match with confirm password",
			&pb.ChangePasswordRequest{
				Id:              1,
				OldPassword:     "111111",
				Password:        "111111",
				ConfirmPassword: "11111",
			},
			&pb.ChangePasswordResponse{
				Status: http.StatusBadRequest,
				Error:  "password-does'nt-match-with-confirm-password",
			},
		},
		{
			"User not found",
			&pb.ChangePasswordRequest{
				Id:              9999,
				OldPassword:     "111111",
				Password:        "111111",
				ConfirmPassword: "111111",
			},
			&pb.ChangePasswordResponse{
				Status: http.StatusNotFound,
				Error:  "user-not-found",
			},
		},
		{
			"Password doesnt match",
			&pb.ChangePasswordRequest{
				Id:              1,
				OldPassword:     "wrong password",
				Password:        "111111",
				ConfirmPassword: "111111",
			},
			&pb.ChangePasswordResponse{
				Status: http.StatusBadRequest,
				Error:  "password-not-match",
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			response, err := client.ChangePassword(ctx, tc.req)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)
		})
	}
}

func TestUploadImage(t *testing.T) {
	t.Parallel()

	testImageFolder := "../tmp"

	imagePath := fmt.Sprintf("%s/avatar.jpg", testImageFolder)
	file, err := os.Open(imagePath)
	require.NoError(t, err)
	defer file.Close()

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	stream, err := client.UploadImage(ctx)

	imageType := filepath.Ext(imagePath)
	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				UserId:    "1",
				ImageType: imageType,
			},
		},
	}

	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}

		require.NoError(t, err)
		size += n
		req := &pb.UploadImageRequest{
			Data: &pb.UploadImageRequest_ChunkData{
				ChunkData: buffer[:n],
			},
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}

	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.NotZero(t, res.GetId())
	require.EqualValues(t, size, res.GetSize())

	savedImagePath := fmt.Sprintf("%s/%s%s", testImageFolder, res.GetId(), imageType)
	require.FileExists(t, savedImagePath)

	file, err = os.Open(savedImagePath)
	require.NoError(t, err)

	err = file.Close()
	require.NoError(t, err)

	// Comments for windows only
	// err = os.Remove(savedImagePath)
	// require.NoError(t, err)

}
