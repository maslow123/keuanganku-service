package services

import (
	"context"
	"log"
	"net/http"
	"testing"

	"github.com/maslow123/transactions/pkg/pb"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {
	testCases := []struct {
		name string
		req  *pb.CreateTransactionRequest
		resp *pb.CreateTransactionResponse
	}{
		{
			"OK",
			&pb.CreateTransactionRequest{
				UserId:     1,
				PosId:      1,
				Total:      2000,
				Details:    "Dikasih ibu",
				ActionType: 0,
			},
			&pb.CreateTransactionResponse{
				Status: int32(http.StatusCreated),
				Error:  "",
			},
		},
		{
			"OK",
			&pb.CreateTransactionRequest{
				UserId:     1,
				PosId:      1,
				Total:      2000,
				Details:    "Beli cireng",
				ActionType: 1,
			},
			&pb.CreateTransactionResponse{
				Status: int32(http.StatusCreated),
				Error:  "",
			},
		},
		{
			"Invalid UserID",
			&pb.CreateTransactionRequest{
				UserId:     0,
				PosId:      1,
				Total:      2000,
				Details:    "Beli cireng",
				ActionType: 0,
			},
			&pb.CreateTransactionResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-user-id",
			},
		},
		{
			"Invalid PosID",
			&pb.CreateTransactionRequest{
				UserId:     1,
				PosId:      0,
				Total:      2000,
				Details:    "Beli cireng",
				ActionType: 0,
			},
			&pb.CreateTransactionResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-pos-id",
			},
		},
		{
			"Invalid Total",
			&pb.CreateTransactionRequest{
				UserId:     1,
				PosId:      1,
				Total:      0,
				Details:    "Beli cireng",
				ActionType: 0,
			},
			&pb.CreateTransactionResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-total",
			},
		},
		{
			"Invalid Balance Type",
			&pb.CreateTransactionRequest{
				UserId:     1,
				PosId:      1,
				Total:      5000,
				Details:    "Beli cireng",
				ActionType: 3,
			},
			&pb.CreateTransactionResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-action-type",
			},
		},
		{
			"Invalid Details",
			&pb.CreateTransactionRequest{
				UserId:     1,
				PosId:      1,
				Total:      2000,
				Details:    "",
				ActionType: 0,
			},
			&pb.CreateTransactionResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-details",
			},
		},
		{
			"Pos Not Found",
			&pb.CreateTransactionRequest{
				UserId:     1,
				PosId:      9999999,
				Total:      2000,
				Details:    "Beli cireng",
				ActionType: 0,
			},
			&pb.CreateTransactionResponse{
				Status: int32(http.StatusNotFound),
				Error:  "pos-not-found",
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewTransactionServiceClient(conn)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			response, err := client.CreateTransaction(ctx, tc.req)
			log.Println(response)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)
		})
	}
}

func TestGetTransactionList(t *testing.T) {
	testCases := []struct {
		name string
		req  *pb.GetTransactionListRequest
		resp *pb.GetTransactionListResponse
	}{
		{
			"OK",
			&pb.GetTransactionListRequest{
				UserId: 1,
				Page:   1,
				Limit:  5,
			},
			&pb.GetTransactionListResponse{
				Status: int32(http.StatusOK),
				Error:  "",
			},
		},
		{
			"Invalid UserID",
			&pb.GetTransactionListRequest{
				UserId: 0,
				Page:   1,
				Limit:  5,
			},
			&pb.GetTransactionListResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-user-id",
			},
		},
		{
			"Invalid Page",
			&pb.GetTransactionListRequest{
				UserId: 1,
				Page:   0,
				Limit:  5,
			},
			&pb.GetTransactionListResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-page",
			},
		},
		{
			"Invalid Limit",
			&pb.GetTransactionListRequest{
				UserId: 1,
				Page:   1,
				Limit:  0,
			},
			&pb.GetTransactionListResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-limit",
			},
		},
		{
			"Transaction Not Found",
			&pb.GetTransactionListRequest{
				UserId: 1,
				Page:   100,
				Limit:  100,
			},
			&pb.GetTransactionListResponse{
				Status: int32(http.StatusNotFound),
				Error:  "transaction-not-found",
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewTransactionServiceClient(conn)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			response, err := client.GetTransactionByUser(ctx, tc.req)
			log.Println(response)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)
			if response.Status == int32(http.StatusOK) {
				require.NotEmpty(t, response.Transaction)
			}
		})
	}
}
