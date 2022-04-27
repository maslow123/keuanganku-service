package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/maslow123/pos/pkg/pb"
	"github.com/maslow123/pos/utils"
	"github.com/stretchr/testify/require"
)

var lastInsertedId int

func TestCreatePos(t *testing.T) {
	testCases := []struct {
		name string
		req  *pb.CreatePosRequest
		resp *pb.CreatePosResponse
	}{
		{
			"OK",
			&pb.CreatePosRequest{
				Name:   utils.RandomString(10),
				Type:   0, // outflow
				Color:  "#FF00FF",
				UserId: 1,
			},
			&pb.CreatePosResponse{
				Status: int64(http.StatusCreated),
				Error:  "",
			},
		},
		{
			"Invalid UserID",
			&pb.CreatePosRequest{
				Name:   utils.RandomString(10),
				Type:   0, // outflow
				Color:  "#FF00FF",
				UserId: 0,
			},
			&pb.CreatePosResponse{
				Status: int64(http.StatusBadRequest),
				Error:  "invalid-user",
			},
		},
		{
			"Invalid Pos Name",
			&pb.CreatePosRequest{
				Name:   "",
				Type:   0, // inflow
				Color:  "#FF00FF",
				UserId: 1,
			},
			&pb.CreatePosResponse{
				Status: int64(http.StatusBadRequest),
				Error:  "invalid-name",
			},
		},
		{
			"Invalid Pos Type",
			&pb.CreatePosRequest{
				Name:   utils.RandomString(10),
				Type:   3,
				Color:  "#FF00FF",
				UserId: 1,
			},
			&pb.CreatePosResponse{
				Status: int64(http.StatusBadRequest),
				Error:  "invalid-type",
			},
		},
		{
			"Invalid Pos Color",
			&pb.CreatePosRequest{
				Name:   utils.RandomString(10),
				Type:   1, // outflow
				Color:  "",
				UserId: 1,
			},
			&pb.CreatePosResponse{
				Status: int64(http.StatusBadRequest),
				Error:  "invalid-color",
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewPosServiceClient(conn)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			response, err := client.CreatePos(ctx, tc.req)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)

			lastInsertedId = int(response.Id)
		})
	}
}

func TestPosList(t *testing.T) {
	testCases := []struct {
		name string
		req  *pb.GetPosListRequest
		resp *pb.GetPosListResponse
	}{
		{
			"OK",
			&pb.GetPosListRequest{
				UserId: 1,
				Type:   0,
				Limit:  10,
				Page:   1,
			},
			&pb.GetPosListResponse{
				Status: int32(http.StatusOK),
				Error:  "",
			},
		},
		{
			"Invalid UserID",
			&pb.GetPosListRequest{
				UserId: 0,
				Type:   0,
				Limit:  10,
				Page:   1,
			},
			&pb.GetPosListResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-user-id",
			},
		},
		{
			"Invalid Type",
			&pb.GetPosListRequest{
				UserId: 0,
				Type:   2,
				Limit:  10,
				Page:   1,
			},
			&pb.GetPosListResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-user-id",
			},
		},
		{
			"Invalid Limit",
			&pb.GetPosListRequest{
				UserId: 1,
				Type:   0,
				Limit:  0,
				Page:   1,
			},
			&pb.GetPosListResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-limit",
			},
		},
		{
			"Invalid Page",
			&pb.GetPosListRequest{
				UserId: 1,
				Type:   0,
				Limit:  10,
				Page:   0,
			},
			&pb.GetPosListResponse{
				Status: int32(http.StatusBadRequest),
				Error:  "invalid-page",
			},
		},
		{
			"Pos Not found",
			&pb.GetPosListRequest{
				UserId: 999,
				Type:   0,
				Limit:  10,
				Page:   1,
			},
			&pb.GetPosListResponse{
				Status: int32(http.StatusNotFound),
				Error:  "pos-not-found",
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewPosServiceClient(conn)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			response, err := client.GetPosByUser(ctx, tc.req)
			log.Println("Response: ", response)
			log.Println("err: ", err)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)

			if response.Status == int32(http.StatusOK) {
				require.NotEmpty(t, response.Pos)
			}
		})
	}
}
func TestPosDetail(t *testing.T) {
	lastInsertedId = 1

	testCases := []struct {
		name string
		req  *pb.PosDetailRequest
		resp *pb.PosDetailResponse
	}{
		{
			"OK",
			&pb.PosDetailRequest{
				Id: int64(lastInsertedId),
			},
			&pb.PosDetailResponse{
				Status: int64(http.StatusOK),
				Error:  "",
			},
		},
		{
			"Invalid Pos ID",
			&pb.PosDetailRequest{
				Id: int64(0),
			},
			&pb.PosDetailResponse{
				Status: int64(http.StatusBadRequest),
				Error:  "invalid-id",
			},
		},
		{
			"Pos Not Found",
			&pb.PosDetailRequest{
				Id: int64(9999),
			},
			&pb.PosDetailResponse{
				Status: int64(http.StatusNotFound),
				Error:  "pos-not-found",
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewPosServiceClient(conn)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			response, err := client.PosDetail(ctx, tc.req)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)

			if response.Status == int64(http.StatusOK) {
				require.NotNil(t, response.Pos)
			}
		})
	}
}
func TestUpdatePos(t *testing.T) {
	testCases := []struct {
		name string
		req  *pb.UpdatePosRequest
		resp *pb.UpdatePosResponse
	}{
		{
			"OK",
			&pb.UpdatePosRequest{
				Id:    2,
				Name:  utils.RandomString(10),
				Color: fmt.Sprintf("#%s", utils.RandomString(6)),
			},
			&pb.UpdatePosResponse{
				Status: int64(http.StatusOK),
				Error:  "",
			},
		},
		{
			"Invalid ID",
			&pb.UpdatePosRequest{
				Id:    0,
				Name:  utils.RandomString(10),
				Color: fmt.Sprintf("#%s", utils.RandomString(6)),
			},
			&pb.UpdatePosResponse{
				Status: int64(http.StatusBadRequest),
				Error:  "invalid-id",
			},
		},
		{
			"Invalid Name",
			&pb.UpdatePosRequest{
				Id:    1,
				Name:  "",
				Color: fmt.Sprintf("#%s", utils.RandomString(6)),
			},
			&pb.UpdatePosResponse{
				Status: int64(http.StatusBadRequest),
				Error:  "invalid-name",
			},
		},

		{
			"Invalid Color",
			&pb.UpdatePosRequest{
				Id:    1,
				Name:  utils.RandomString(10),
				Color: "",
			},
			&pb.UpdatePosResponse{
				Status: int64(http.StatusBadRequest),
				Error:  "invalid-color",
			},
		},
		{
			"Not Found",
			&pb.UpdatePosRequest{
				Id:    99999,
				Name:  utils.RandomString(10),
				Color: fmt.Sprintf("#%s", utils.RandomString(6)),
			},
			&pb.UpdatePosResponse{
				Status: int64(http.StatusNotFound),
				Error:  "pos-not-found",
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewPosServiceClient(conn)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			response, err := client.UpdatePosByUser(ctx, tc.req)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)

			if response.Status == int64(http.StatusOK) {
				require.NotNil(t, response.Pos)
				require.Equal(t, response.Pos.Name, tc.req.Name)
			}
		})
	}
}
func TestDeletePos(t *testing.T) {
	testCases := []struct {
		name     string
		getPosID func(t *testing.T, ctx context.Context, client pb.PosServiceClient) int64
		resp     *pb.DeletePosResponse
	}{
		{
			"OK",
			func(t *testing.T, ctx context.Context, client pb.PosServiceClient) int64 {
				arg := &pb.CreatePosRequest{
					UserId: 1,
					Name:   utils.RandomString(10),
					Type:   0,
					Color:  fmt.Sprintf("#%s", utils.RandomString(6)),
				}
				todo, err := client.CreatePos(ctx, arg)
				require.NoError(t, err)

				return todo.Id

			},
			&pb.DeletePosResponse{
				Status: int64(http.StatusOK),
				Error:  "",
			},
		},
		{
			"Invalid ID",
			func(t *testing.T, ctx context.Context, client pb.PosServiceClient) int64 {
				return 0
			},
			&pb.DeletePosResponse{
				Status: int64(http.StatusBadRequest),
				Error:  "invalid-id",
			},
		},
		{
			"Pos Not Found",
			func(t *testing.T, ctx context.Context, client pb.PosServiceClient) int64 {
				return 999
			},
			&pb.DeletePosResponse{
				Status: int64(http.StatusNotFound),
				Error:  "pos-not-found",
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewPosServiceClient(conn)
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			posId := tc.getPosID(t, ctx, client)

			req := &pb.DeletePosRequest{
				Id: posId,
			}

			response, err := client.DeletePosByUser(ctx, req)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)
		})
	}
}

func TestUpdateTotalPos(t *testing.T) {
	var lastInsertedId int64
	testCases := []struct {
		name     string
		getPosID func(t *testing.T, ctx context.Context, client pb.PosServiceClient) int64
		action   pb.UpdateTotalPosRequest_ActionTransaction
		amount   int64
		resp     *pb.UpdateTotalPosResponse
	}{
		{
			"OK Increment",
			func(t *testing.T, ctx context.Context, client pb.PosServiceClient) int64 {
				arg := &pb.CreatePosRequest{
					UserId: 1,
					Name:   utils.RandomString(10),
					Type:   0,
					Color:  fmt.Sprintf("#%s", utils.RandomString(6)),
				}
				todo, err := client.CreatePos(ctx, arg)
				require.NoError(t, err)

				lastInsertedId = todo.Id
				return todo.Id

			},
			pb.UpdateTotalPosRequest_INCREASE,
			5000,
			&pb.UpdateTotalPosResponse{
				Status: int64(http.StatusOK),
				Error:  "",
				Total:  5000,
			},
		},
		{
			"OK Decrement",
			func(t *testing.T, ctx context.Context, client pb.PosServiceClient) int64 {
				return lastInsertedId

			},
			pb.UpdateTotalPosRequest_DECREASE,
			5000,
			&pb.UpdateTotalPosResponse{
				Status: int64(http.StatusOK),
				Error:  "",
				Total:  0,
			},
		},
		{
			"Invalid ID",
			func(t *testing.T, ctx context.Context, client pb.PosServiceClient) int64 {
				return 0

			},
			pb.UpdateTotalPosRequest_INCREASE,
			5000,
			&pb.UpdateTotalPosResponse{
				Status: int64(http.StatusBadRequest),
				Error:  "invalid-id",
				Total:  0,
			},
		},
		{
			"Invalid Action",
			func(t *testing.T, ctx context.Context, client pb.PosServiceClient) int64 {
				return lastInsertedId

			},
			3,
			5000,
			&pb.UpdateTotalPosResponse{
				Status: int64(http.StatusBadRequest),
				Error:  "invalid-action",
				Total:  0,
			},
		},
		{
			"Invalid Amount",
			func(t *testing.T, ctx context.Context, client pb.PosServiceClient) int64 {
				return lastInsertedId

			},
			pb.UpdateTotalPosRequest_INCREASE,
			0,
			&pb.UpdateTotalPosResponse{
				Status: int64(http.StatusBadRequest),
				Error:  "invalid-amount",
				Total:  0,
			},
		},
		{
			"Pos Not Found",
			func(t *testing.T, ctx context.Context, client pb.PosServiceClient) int64 {
				return 99999

			},
			pb.UpdateTotalPosRequest_INCREASE,
			5000,
			&pb.UpdateTotalPosResponse{
				Status: int64(http.StatusNotFound),
				Error:  "pos-not-found",
				Total:  0,
			},
		},
	}

	ctx := context.Background()
	conn := checkConnection(ctx, t)
	defer conn.Close()

	client := pb.NewPosServiceClient(conn)
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			posId := tc.getPosID(t, ctx, client)

			req := &pb.UpdateTotalPosRequest{
				Action: tc.action,
				Amount: tc.amount,
				Id:     posId,
			}

			response, err := client.UpdateTotalPosByUser(ctx, req)
			require.NoError(t, err)

			require.Equal(t, tc.resp.Status, response.Status)
			require.Equal(t, tc.resp.Error, response.Error)
			if response.Status == int64(http.StatusOK) {
				require.Equal(t, tc.resp.Total, response.Total)
			}
		})
	}
}
