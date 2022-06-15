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

func TestDetailUserTransaction(t *testing.T) {
	testCases := []struct {
		name             string
		getTransactionId func(t *testing.T, server *ServiceClient, authorizationHeader string) int32
		checkResponse    func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			getTransactionId: func(t *testing.T, server *ServiceClient, authorizationHeader string) int32 {
				return createRandomTransaction(t, server, authorizationHeader, int32(time.Now().Unix()), 10000)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid Transaction ID",
			getTransactionId: func(t *testing.T, server *ServiceClient, authorizationHeader string) int32 {
				return 0
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Transaction Not Found",
			getTransactionId: func(t *testing.T, server *ServiceClient, authorizationHeader string) int32 {
				return 999999
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
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
			url := fmt.Sprintf("/transactions/detail/%d", id)
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
				return createRandomTransaction(t, server, authorizationHeader, int32(time.Now().Unix()), 10000)
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

func TestGetPercentageExpenditure(t *testing.T) {
	testCases := []struct {
		name          string
		query         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			query: "start_date=2022-01-02&end_date=2022-01-01",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var response pb.GetPercentageExpenditureResponse
				err = json.Unmarshal(data, &response)
				require.NoError(t, err)

				require.Equal(t, float32(90), response.Percentage)
			},
		},
		{
			name:  "Empty Start Date",
			query: "start_date=&end_date=2022-01-01",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var response pb.GetPercentageExpenditureResponse
				err = json.Unmarshal(data, &response)
				require.NoError(t, err)

				require.Equal(t, "invalid-start-date", response.Error)
			},
		},
		{
			name:  "Empty End Date",
			query: "start_date=2022-01-02&end_date=",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var response pb.GetPercentageExpenditureResponse
				err = json.Unmarshal(data, &response)
				require.NoError(t, err)

				require.Equal(t, "invalid-end-date", response.Error)
			},
		},
		{
			name:  "Invalid Start Date",
			query: "start_date=invalid&end_date=2022-01-01",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var response pb.GetPercentageExpenditureResponse
				err = json.Unmarshal(data, &response)
				require.NoError(t, err)

				require.Equal(t, "invalid-start-date", response.Error)
			},
		},
		{
			name:  "Invalid End Date",
			query: "start_date=2022-01-02&end_date=invalid",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var response pb.GetPercentageExpenditureResponse
				err = json.Unmarshal(data, &response)
				require.NoError(t, err)

				require.Equal(t, "invalid-end-date", response.Error)
			},
		},
		{
			name:  "Transaction Not Found",
			query: "start_date=2022-03-03&end_date=2022-03-03",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)

				var response pb.GetPercentageExpenditureResponse
				err = json.Unmarshal(data, &response)
				require.NoError(t, err)

				require.Equal(t, "transaction-not-found", response.Error)
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
			var transactionsId []int32
			// create dummy transaction
			if tc.name == "OK" {
				dates := []string{"2022-01-02", "2022-01-01"}
				totals := []int32{10000, 100000}
				for i, date := range dates {
					unixDate, err := time.Parse("2006-01-02", date)
					require.NoError(t, err)

					transactionId := createRandomTransaction(t, server, authorizationHeader, int32(unixDate.Unix()), totals[i])
					transactionsId = append(transactionsId, transactionId)
				}
			}

			url := fmt.Sprintf("/transactions/expenditure?%s", tc.query)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			request.Header.Set("Authorization", authorizationHeader)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)

			// delete transaction
			if len(transactionsId) > 0 {
				for _, txID := range transactionsId {
					err := deleteTransction(t, server, authorizationHeader, txID)
					require.NoError(t, err)
				}
			}
		})
	}
}

func createRandomTransaction(t *testing.T, server *ServiceClient, authorizationHeader string, createdAt, total int32) int32 {
	recorder := httptest.NewRecorder()

	body := gin.H{
		"pos_id":      1,
		"total":       total,
		"details":     "Beli cireng",
		"action_type": 1,
		"type":        0,
		"date":        createdAt,
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

func deleteTransction(t *testing.T, server *ServiceClient, authorizationHeader string, transactionId int32) error {
	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("/transactions/%d", transactionId)
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)

	request.Header.Set("Authorization", authorizationHeader)
	server.Router.ServeHTTP(recorder, request)

	data, err := ioutil.ReadAll(recorder.Body)
	require.NoError(t, err)

	log.Println("==== DATA ====", string(data))

	var tx pb.DeleteTransactionResponse
	err = json.Unmarshal(data, &tx)
	return err
}
