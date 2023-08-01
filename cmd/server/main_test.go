package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetricsHandler(t *testing.T) {
	storage := make(MemStorage)
	handler := http.HandlerFunc(storage.MetricsHandler)

	testCases := []struct {
		name       string
		method     string
		path       string
		body       string
		statusCode int
	}{
		{
			name:       "ValidGaugeRequest",
			method:     http.MethodPost,
			path:       "/update/gauge/metric1/10",
			body:       "",
			statusCode: http.StatusOK,
		},
		{
			name:       "ValidCounterRequest",
			method:     http.MethodPost,
			path:       "/update/counter/metric2/5",
			body:       "",
			statusCode: http.StatusOK,
		},
		{
			name:       "InvalidMethod",
			method:     http.MethodGet,
			path:       "/update/gauge/metric1/10",
			body:       "",
			statusCode: http.StatusMethodNotAllowed,
		},
		{
			name:       "InvalidPath",
			method:     http.MethodPost,
			path:       "/invalid/path",
			body:       "",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "InvalidMetricValue",
			method:     http.MethodPost,
			path:       "/update/gauge/metric1/invalid",
			body:       "",
			statusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			req, err := http.NewRequest(testCase.method, testCase.path, strings.NewReader(testCase.body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			if recorder.Code != testCase.statusCode {
				t.Errorf("Expected status code %d, but got %d", testCase.statusCode, recorder.Code)
			}
		})
	}
}
