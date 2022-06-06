package services

import (
	"context"
	"net/http"
	"testing"
	"time"

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
				Type:       0,
				Date:       int32(time.Now().Unix()),
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
				Type:       0,
				Date:       int32(time.Now().Unix()),
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
				Type:       0,
				Date:       int32(time.Now().Unix()),
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
				Type:       0,
				Date:       int32(time.Now().Unix()),
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
				Type:       0,
				Date:       int32(time.Now().Unix()),
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
				Type:       0,
				Date:       int32(time.Now().Unix()),
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
				Type:       0,
				Date:       int32(time.Now().Unix()),
			},
			&pb.CreateTransactionResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-details",
			},
		},
		{
			"Invalid Type",
			&pb.CreateTransactionRequest{
				UserId:     1,
				PosId:      1,
				Total:      2000,
				Details:    "test",
				ActionType: 0,
				Type:       2,
				Date:       int32(time.Now().Unix()),
			},
			&pb.CreateTransactionResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-type",
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
				Type:       0,
				Date:       int32(time.Now().Unix()),
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
				UserId:    1,
				Page:      1,
				Limit:     5,
				Action:    0,
				StartDate: int32(time.Now().Unix()),
				EndDate:   int32(time.Now().Unix()),
			},
			&pb.GetTransactionListResponse{
				Status: int32(http.StatusOK),
				Error:  "",
			},
		},
		{
			"Invalid UserID",
			&pb.GetTransactionListRequest{
				UserId:    0,
				Page:      1,
				Limit:     5,
				Action:    0,
				StartDate: int32(time.Now().Unix()),
				EndDate:   int32(time.Now().Unix()),
			},
			&pb.GetTransactionListResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-user-id",
			},
		},
		{
			"Invalid Page",
			&pb.GetTransactionListRequest{
				UserId:    1,
				Page:      0,
				Limit:     5,
				Action:    0,
				StartDate: int32(time.Now().Unix()),
				EndDate:   int32(time.Now().Unix()),
			},
			&pb.GetTransactionListResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-page",
			},
		},
		{
			"Invalid Limit",
			&pb.GetTransactionListRequest{
				UserId:    1,
				Page:      1,
				Limit:     0,
				Action:    0,
				StartDate: int32(time.Now().Unix()),
				EndDate:   int32(time.Now().Unix()),
			},
			&pb.GetTransactionListResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-limit",
			},
		},
		{
			"Invalid Type",
			&pb.GetTransactionListRequest{
				UserId:    1,
				Page:      1,
				Limit:     10,
				Action:    3,
				StartDate: int32(time.Now().Unix()),
				EndDate:   int32(time.Now().Unix()),
			},
			&pb.GetTransactionListResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-type",
			},
		},
		{
			"Transaction Not Found",
			&pb.GetTransactionListRequest{
				UserId:    1,
				Page:      100,
				Limit:     100,
				Action:    0,
				StartDate: int32(time.Now().Unix()),
				EndDate:   int32(time.Now().Unix()),
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
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)
			if response.Status == int32(http.StatusOK) {
				require.NotEmpty(t, response.Transaction)
			}
		})
	}
}

func TestDetailTransaction(t *testing.T) {
	testCases := []struct {
		name             string
		userId           int32
		getTransactionId func(t *testing.T, ctx context.Context, client pb.TransactionServiceClient) int32
		resp             *pb.DetailTransactionResponse
	}{
		{
			"OK",
			1,
			func(t *testing.T, ctx context.Context, client pb.TransactionServiceClient) int32 {
				arg := &pb.CreateTransactionRequest{
					UserId:     1,
					PosId:      1,
					Total:      2000,
					Details:    "Test Detail Transaction",
					ActionType: 0,
					Type:       0,
				}
				tx, err := client.CreateTransaction(ctx, arg)
				require.NoError(t, err)

				return tx.Id

			},
			&pb.DetailTransactionResponse{
				Status: int32(http.StatusOK),
				Error:  "",
				Transaction: &pb.Transaction{
					Details: "Test Detail Transaction",
				},
			},
		},
		{
			"Invalid User ID",
			0,
			func(t *testing.T, ctx context.Context, client pb.TransactionServiceClient) int32 {
				return 1

			},
			&pb.DetailTransactionResponse{
				Status:      int32(http.StatusBadRequest),
				Error:       "invalid-user-id",
				Transaction: &pb.Transaction{},
			},
		},
		{
			"Invalid Transaction ID",
			1,
			func(t *testing.T, ctx context.Context, client pb.TransactionServiceClient) int32 {
				return 0

			},
			&pb.DetailTransactionResponse{
				Status:      int32(http.StatusBadRequest),
				Error:       "invalid-transaction-id",
				Transaction: &pb.Transaction{},
			},
		},
		{
			"Transaction Not Found",
			1,
			func(t *testing.T, ctx context.Context, client pb.TransactionServiceClient) int32 {
				return 99999

			},
			&pb.DetailTransactionResponse{
				Status:      int32(http.StatusNotFound),
				Error:       "transaction-not-found",
				Transaction: &pb.Transaction{},
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
			transactionId := tc.getTransactionId(t, ctx, client)
			req := &pb.DetailTransactionRequest{
				Id:     transactionId,
				UserId: tc.userId,
			}
			response, err := client.DetailTransaction(ctx, req)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)
		})
	}
}
func TestDeleteTransaction(t *testing.T) {
	testCases := []struct {
		name             string
		getTransactionId func(t *testing.T, ctx context.Context, client pb.TransactionServiceClient) int32
		resp             *pb.DeletePosResponse
	}{
		{
			"OK",
			func(t *testing.T, ctx context.Context, client pb.TransactionServiceClient) int32 {
				arg := &pb.CreateTransactionRequest{
					UserId:     1,
					PosId:      1,
					Total:      2000,
					Details:    "Dikasih ibu",
					ActionType: 0,
					Type:       0,
				}
				tx, err := client.CreateTransaction(ctx, arg)
				require.NoError(t, err)

				return tx.Id

			},
			&pb.DeletePosResponse{
				Status: int32(http.StatusOK),
				Error:  "",
			},
		},
		{
			"Invalid Transaction ID",
			func(t *testing.T, ctx context.Context, client pb.TransactionServiceClient) int32 {
				return 0

			},
			&pb.DeletePosResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-transaction-id",
			},
		},
		{
			"Transaction Not Found",
			func(t *testing.T, ctx context.Context, client pb.TransactionServiceClient) int32 {
				return 99999

			},
			&pb.DeletePosResponse{
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
			transactionId := tc.getTransactionId(t, ctx, client)
			req := &pb.DeleteTransactionRequest{
				Id:     transactionId,
				UserId: 1,
			}
			response, err := client.DeleteTransactionByUser(ctx, req)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)
		})
	}
}
