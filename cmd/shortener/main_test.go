package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/levshindenis/sprint1/internal/app/handlers"
)

func TestHSStorage_PostHandler(t *testing.T) {
	var serv struct {
		handlers.HStorage
	}
	serv.InitStorage()

	tests := []struct {
		name         string
		method       string
		address      string
		requestBody  string
		expectedCode int
		emptyBody    bool
	}{
		{
			name:         "Good test",
			method:       http.MethodPost,
			address:      "localhost:8000",
			requestBody:  "https://yandex.ru/",
			expectedCode: http.StatusCreated,
			emptyBody:    false,
		},
		{
			name:         "Bad method",
			method:       http.MethodGet,
			address:      "localhost:8000",
			requestBody:  "https://yandex.ru/",
			expectedCode: http.StatusBadRequest,
			emptyBody:    true,
		},
		{
			name:         "Bad url",
			method:       http.MethodPost,
			address:      "localhost:8000",
			requestBody:  "Hello",
			expectedCode: http.StatusBadRequest,
			emptyBody:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serv.SetBaseSA(tt.address)
			r := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			serv.PostHandler(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
			if !tt.emptyBody {
				assert.Contains(t, w.Body.String(), tt.address,
					"Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}

func TestHSStorage_GetHandler(t *testing.T) {
	var serv struct {
		handlers.HStorage
	}
	serv.InitStorage()
	serv.SetStorage("GyuRe0", "https://yandex.ru/")

	tests := []struct {
		name         string
		method       string
		url          string
		expectedCode int
		expectedBody string
		emptyBody    bool
	}{
		{
			name:         "Good test",
			method:       http.MethodGet,
			url:          "/GyuRe0",
			expectedCode: http.StatusTemporaryRedirect,
			expectedBody: "https://yandex.ru/",
			emptyBody:    false,
		},
		{
			name:         "Bad method",
			method:       http.MethodPost,
			url:          "/GyuRe0",
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
			emptyBody:    true,
		},
		{
			name:         "Bad url",
			method:       http.MethodGet,
			url:          "/GyuAe0",
			expectedCode: http.StatusBadRequest,
			expectedBody: "",
			emptyBody:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			serv.GetHandler(w, r)

			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
			if !tt.emptyBody {
				assert.Equal(t, w.Header().Get("Location"), tt.expectedBody,
					"Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}

func TestHSStorage_JSONPostHandler(t *testing.T) {
	var serv struct {
		handlers.HStorage
	}
	serv.InitStorage()

	tests := []struct {
		name         string
		method       string
		address      string
		requestBody  string
		contentType  string
		expectedCode int
		emptyBody    bool
	}{
		{
			name:         "Good test",
			method:       http.MethodPost,
			address:      "localhost:8000",
			requestBody:  "{\"url\":\"https://practicum.yandex.ru\"}",
			contentType:  "application/json",
			expectedCode: http.StatusCreated,
			emptyBody:    false,
		},
		{
			name:         "Bad JSON test",
			method:       http.MethodPost,
			address:      "localhost:8000",
			requestBody:  "{\"url:\"https://practicum.yandex.ru\"}",
			contentType:  "application/json",
			expectedCode: http.StatusBadRequest,
			emptyBody:    true,
		},
		{
			name:         "Bad Content-Type test",
			method:       http.MethodPost,
			address:      "localhost:8000",
			requestBody:  "{\"url\":\"https://practicum.yandex.ru\"}",
			contentType:  "text/plain",
			expectedCode: http.StatusBadRequest,
			emptyBody:    true,
		},
		{
			name:         "Bad Method test",
			method:       http.MethodGet,
			address:      "localhost:8000",
			requestBody:  "{\"url\":\"https://practicum.yandex.ru\"}",
			contentType:  "application/json",
			expectedCode: http.StatusBadRequest,
			emptyBody:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serv.SetBaseSA(tt.address)
			r := httptest.NewRequest(tt.method, "/", strings.NewReader(tt.requestBody))
			w := httptest.NewRecorder()
			r.Header.Set("Content-Type", tt.contentType)
			serv.JSONPostHandler(w, r)
			assert.Equal(t, w.Code, tt.expectedCode, "Код ответа не совпадает с ожидаемым")
			if !tt.emptyBody {
				assert.Contains(t, w.Body.String(), tt.address,
					"Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}
