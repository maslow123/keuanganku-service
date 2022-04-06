package services

import (
	"context"
	"log"
	"net/http"
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
				Status: int64(http.StatusCreated),
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
				Status: int64(http.StatusBadRequest),
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
				Status: int64(http.StatusBadRequest),
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
				Status: int64(http.StatusBadRequest),
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
				Status: int64(http.StatusBadRequest),
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
				Status: int64(http.StatusBadRequest),
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
				Status: int64(http.StatusBadRequest),
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
				Status: int64(http.StatusOK),
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
				Status: int64(http.StatusBadRequest),
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
				Status: int64(http.StatusBadRequest),
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
				Status: int64(http.StatusUnauthorized),
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
				Status: int64(http.StatusNotFound),
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

			if response.Status == int64(http.StatusOK) {
				require.NotEmpty(t, response.Token)
				require.NotNil(t, response.User)
			}
		})
	}
}
