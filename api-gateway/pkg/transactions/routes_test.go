package transactions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {
	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"pos_id":      1,
				"total":       10000,
				"details":     "Beli cireng",
				"action_type": 0,
				"type":        1,
				"date":        time.Now().Unix(),
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

			url := "/transactions/create"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Authorization", authorizationHeader)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetUserTransaction(t *testing.T) {
	testCases := []struct {
		name          string
		query         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			query: "page=1&limit=10&action=0&start_date=0&end_date=0",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "Invalid Page",
			query: "page=0&limit=10&action=0&start_date=0&end_date=0",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Invalid Limit",
			query: "page=1&limit=0&action=0&start_date=0&end_date=0",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Invalid Type",
			query: "page=1&limit=5&action=3&start_date=0&end_date=0",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := fmt.Sprintf("/transactions/list?%s", tc.query)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			request.Header.Set("Authorization", authorizationHeader)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteTransaction(t *testing.T) {
	testCases := []struct {
		name             string
		getTransactionId func(t *testing.T, server *ServiceClient, authorizationHeader string) int32
		checkResponse    func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			getTransactionId: func(t *testing.T, server *ServiceClient, authorizationHeader string) int32 {
				return createRandomTransaction(t, server, authorizationHeader)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
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

			id := tc.getTransactionId(t, server, authorizationHeader)
			url := fmt.Sprintf("/transactions/%d", id)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			request.Header.Set("Authorization", authorizationHeader)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func createRandomTransaction(t *testing.T, server *ServiceClient, authorizationHeader string) int32 {
	recorder := httptest.NewRecorder()

	body := gin.H{
		"pos_id":      1,
		"total":       10000,
		"details":     "Beli cireng",
		"action_type": 0,
		"type":        1,
	}

	data, err := json.Marshal(body)
	log.Println(err)
	require.NoError(t, err)

	url := "/transactions/create"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	request.Header.Set("Authorization", authorizationHeader)
	server.Router.ServeHTTP(recorder, request)

	data, err = ioutil.ReadAll(recorder.Body)
	require.NoError(t, err)

	log.Println("==== DATA ====", string(data))

	var tx pb.CreateTransactionResponse
	err = json.Unmarshal(data, &tx)
	require.NoError(t, err)

	return tx.Id
}
