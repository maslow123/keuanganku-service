package users

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maslow123/api-gateway/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email":    "user1@gmail.com",
				"password": "111111",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Invalid Email",
			body: gin.H{
				"email":    "",
				"password": "111111",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Password",
			body: gin.H{
				"email":    "user1",
				"password": "",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Password not match",
			body: gin.H{
				"email":    "user1@gmail.com",
				"password": "wrong password",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "User not found",
			body: gin.H{
				"email":    "user9999",
				"password": "wrong password",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := NewServer(t)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users/login"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestRegister(t *testing.T) {
	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"name":             utils.RandomString(10),
				"email":            utils.RandomString(10),
				"password":         "111111",
				"confirm_password": "111111",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
			},
		},
		{
			name: "Invalid Name",
			body: gin.H{
				"name":             "",
				"email":            utils.RandomString(10),
				"password":         "111111",
				"confirm_password": "111111",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Username",
			body: gin.H{
				"name":             utils.RandomString(10),
				"email":            "",
				"password":         "111111",
				"confirm_password": "111111",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Password",
			body: gin.H{
				"name":             utils.RandomString(10),
				"email":            utils.RandomString(10),
				"password":         "",
				"confirm_password": "111111",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Invalid Confirm Password",
			body: gin.H{
				"name":             utils.RandomString(10),
				"email":            utils.RandomString(10),
				"password":         "111111",
				"confirm_password": "",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Confirm Password not match",
			body: gin.H{
				"name":             utils.RandomString(10),
				"email":            utils.RandomString(10),
				"password":         "111111",
				"confirm_password": "not match",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "Existing User",
			body: gin.H{
				"name":             utils.RandomString(10),
				"email":            "user1@gmail.com",
				"password":         "111111",
				"confirm_password": "111111",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := NewServer(t)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users/register"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"email": "user2@gmail.com",
				"name":  "Omama Olala",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				log.Println(recorder)
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
			server := NewServer(t)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users/update"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Authorization", authorizationHeader)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestChangePassword(t *testing.T) {
	testCases := []struct {
		name          string
		body          gin.H
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"old_password":     "111111",
				"password":         "111111",
				"confirm_password": "111111",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				log.Println(recorder)
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
			server := NewServer(t)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users/change-password"
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
			require.NoError(t, err)

			request.Header.Set("Authorization", authorizationHeader)
			server.Router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestUploadImage(t *testing.T) {

	server := NewServer(t)
	authorizationHeader := addAuthorization(t, server)

	filePath := "./images/avatar.jpg"
	fieldName := "file"
	body := new(bytes.Buffer)

	mw := multipart.NewWriter(body)

	file, err := os.Open(filePath)
	require.NoError(t, err)
	defer file.Close()

	w, err := mw.CreateFormFile(fieldName, filePath)
	require.NoError(t, err)

	_, err = io.Copy(w, file)
	require.NoError(t, err)

	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/users/upload", body)
	req.Header.Add("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", authorizationHeader)

	res := httptest.NewRecorder()
	server.Router.ServeHTTP(res, req)

	require.Equal(t, http.StatusOK, res.Code)
}
