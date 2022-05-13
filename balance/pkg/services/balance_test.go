package services

import (
	"context"
	"net/http"
	"testing"

	"github.com/maslow123/balance/pkg/pb"
	"github.com/stretchr/testify/require"
)

func TestUpsertBalance(t *testing.T) {
	testCases := []struct {
		name string
		req  *pb.UpsertBalanceRequest
		resp *pb.UpsertBalanceResponse
	}{
		{
			"OK",
			&pb.UpsertBalanceRequest{
				UserId: 1,
				Type:   0,
				Total:  3000,
				Action: pb.UpsertBalanceRequest_ActionType(pb.UpsertBalanceRequest_ActionType_value["INCREASE"]),
			},
			&pb.UpsertBalanceResponse{
				Status:         http.StatusCreated,
				Error:          "",
				CurrentBalance: 3000,
			},
		},
		{
			"OK",
			&pb.UpsertBalanceRequest{
				UserId: 1,
				Type:   0,
				Total:  3000,
				Action: pb.UpsertBalanceRequest_ActionType(pb.UpsertBalanceRequest_ActionType_value["DECREASE"]),
			},
			&pb.UpsertBalanceResponse{
				Status:         http.StatusCreated,
				Error:          "",
				CurrentBalance: 0,
			},
		},
		{
			"Invalid User ID",
			&pb.UpsertBalanceRequest{
				UserId: 0,
				Type:   0,
				Total:  3000,
				Action: pb.UpsertBalanceRequest_ActionType(pb.UpsertBalanceRequest_ActionType_value["DECREASE"]),
			},
			&pb.UpsertBalanceResponse{
				Status:         http.StatusBadRequest,
				Error:          "invalid-user-id",
				CurrentBalance: 0,
			},
		},
		{
			"Invalid Type",
			&pb.UpsertBalanceRequest{
				UserId: 1,
				Type:   3,
				Total:  3000,
				Action: pb.UpsertBalanceRequest_ActionType(pb.UpsertBalanceRequest_ActionType_value["DECREASE"]),
			},
			&pb.UpsertBalanceResponse{
				Status:         http.StatusBadRequest,
				Error:          "invalid-type",
				CurrentBalance: 0,
			},
		},
		{
			"Invalid Action",
			&pb.UpsertBalanceRequest{
				UserId: 1,
				Type:   0,
				Total:  3000,
				Action: 3,
			},
			&pb.UpsertBalanceResponse{
				Status:         http.StatusBadRequest,
				Error:          "invalid-action",
				CurrentBalance: 0,
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewBalanceServiceClient(conn)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			response, err := client.UpsertBalance(ctx, tc.req)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)
		})
	}
}

func TestGetUserBalance(t *testing.T) {
	testCases := []struct {
		name string
		req  *pb.GetUserBalanceRequest
		resp *pb.GetUserBalanceResponse
	}{
		{
			"OK",
			&pb.GetUserBalanceRequest{
				UserId: 1,
			},
			&pb.GetUserBalanceResponse{
				Status: http.StatusOK,
				Error:  "",
			},
		},
		{
			"Invalid UserID",
			&pb.GetUserBalanceRequest{
				UserId: 0,
			},
			&pb.GetUserBalanceResponse{
				Status: http.StatusBadRequest,
				Error:  "invalid-user-id",
			},
		},
		{
			"Balance not found",
			&pb.GetUserBalanceRequest{
				UserId: 99999,
			},
			&pb.GetUserBalanceResponse{
				Status: http.StatusNotFound,
				Error:  "user-balance-not-found",
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewBalanceServiceClient(conn)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			response, err := client.GetUserBalance(ctx, tc.req)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)
		})
	}
}
