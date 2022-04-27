package pos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/transactions/pb"
	"github.com/maslow123/api-gateway/pkg/utils"
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
				"name":  fmt.Sprintf("pos %s", utils.RandomString(10)),
				"type":  0,
				"color": "#FF00FF",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "Invalid Name",
			body: gin.H{
				"name":  "",
				"type":  0,
				"color": "#FF00FF",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Type",
			body: gin.H{
				"name":  fmt.Sprintf("pos %s", utils.RandomString(10)),
				"type":  3,
				"color": "#FF00FF",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Color",
			body: gin.H{
				"name":  fmt.Sprintf("pos %s", utils.RandomString(10)),
				"type":  0,
				"color": "",
			},
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

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/pos/create"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Authorization", authorizationHeader)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestPostList(t *testing.T) {
	testCases := []struct {
		name          string
		query         string
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			query: "page=1&limit=10&type=0",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "Invalid Type",
			query: "page=1&limit=10&type=2",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Invalid Page",
			query: "page=0&limit=10&type=0",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Invalid Limit",
			query: "page=1&limit=0&type=0",
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Pos Not Found",
			query: "page=100&limit=10&type=0",
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

			url := fmt.Sprintf("/pos/list?%s", tc.query)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			request.Header.Set("Authorization", authorizationHeader)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
func TestPosDetail(t *testing.T) {
	testCases := []struct {
		name          string
		posID         int
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			posID: 1,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "Invalid Pos ID",
			posID: 0,
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Pos Not Found",
			posID: 99999,
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

			url := fmt.Sprintf("/pos/%d", tc.posID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			request.Header.Set("Authorization", authorizationHeader)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdatePos(t *testing.T) {
	testCases := []struct {
		name          string
		posID         int64
		body          gin.H
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			posID: 1,
			body: gin.H{
				"name":  fmt.Sprintf("pos %s", utils.RandomString(10)),
				"color": "#FF00FF",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name:  "Invalid Pos ID",
			posID: 0,
			body: gin.H{
				"name":  fmt.Sprintf("pos %s", utils.RandomString(10)),
				"color": "#FF00FF",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Invalid Name",
			posID: 1,
			body: gin.H{
				"name":  "",
				"color": "#FF00FF",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Invalid Color",
			posID: 1,
			body: gin.H{
				"name":  fmt.Sprintf("pos %s", utils.RandomString(10)),
				"color": "",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name:  "Pos Not Found",
			posID: 9999,
			body: gin.H{
				"name":  fmt.Sprintf("pos %s", utils.RandomString(10)),
				"color": "#FF00FF",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	server := NewServer(t)
	authorizationHeader := addAuthorization(t, server)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server = NewServer(t)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/pos/%d", tc.posID)
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Authorization", authorizationHeader)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeletePos(t *testing.T) {
	testCases := []struct {
		name          string
		posID         func(t *testing.T, server *ServiceClient, authorizationHeader string) int64
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			posID: func(t *testing.T, server *ServiceClient, authorizationHeader string) int64 {
				return createRandomPOS(t, server, authorizationHeader)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid ID",
			posID: func(t *testing.T, server *ServiceClient, authorizationHeader string) int64 {
				return 0
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Pos Not Found",
			posID: func(t *testing.T, server *ServiceClient, authorizationHeader string) int64 {
				return 9999
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	server := NewServer(t)
	authorizationHeader := addAuthorization(t, server)

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server = NewServer(t)
			recorder := httptest.NewRecorder()

			id := tc.posID(t, server, authorizationHeader)
			url := fmt.Sprintf("/pos/%d", id)
			request, err := http.NewRequest(http.MethodDelete, url, nil)
			require.NoError(t, err)

			request.Header.Set("Authorization", authorizationHeader)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func createRandomPOS(t *testing.T, server *ServiceClient, authorizationHeader string) int64 {
	recorder := httptest.NewRecorder()

	body := gin.H{
		"user_id": 1,
		"name":    fmt.Sprintf("pos %s", utils.RandomString(10)),
		"type":    0,
		"color":   "#FF00FF",
	}

	data, err := json.Marshal(body)
	require.NoError(t, err)

	url := "/pos/create"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	require.NoError(t, err)

	request.Header.Set("Authorization", authorizationHeader)
	server.Router.ServeHTTP(recorder, request)

	data, err = ioutil.ReadAll(recorder.Body)
	require.NoError(t, err)

	log.Println("==== DATA ====", string(data))

	var pos pb.CreatePosResponse
	err = json.Unmarshal(data, &pos)
	require.NoError(t, err)

	return pos.Id
}
