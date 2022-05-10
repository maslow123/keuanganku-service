package pos

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"github.com/stretchr/testify/require"
)

func TestCreatePos(t *testing.T) {
	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"type":   0,
				"total":  3000,
				"action": pb.UpsertBalanceRequest_ActionType(pb.UpsertBalanceRequest_ActionType_value["INCREASE"]),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "OK",
			body: gin.H{
				"type":   0,
				"total":  3000,
				"action": pb.UpsertBalanceRequest_ActionType(pb.UpsertBalanceRequest_ActionType_value["DECREASE"]),
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
	}

	// set authorizationHeader
	server := NewServer(t)
	authorizationHeader := addAuthorization(t, server)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server = NewServer(t)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/balance/upsert"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Authorization", authorizationHeader)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
