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
			},
			&pb.UpsertBalanceResponse{
				Status: http.StatusCreated,
				Error:  "",
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
