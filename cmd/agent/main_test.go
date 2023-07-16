package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSendMetric(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", req.Method)
		}
		if req.URL.String() != serverAddress+"gauge/metricName/123" {
			t.Errorf("Unexpected URL: %s", req.URL.String())
		}
		body := make([]byte, req.ContentLength)
		_, _ = req.Body.Read(body)
		payload := string(body)
		expectedPayload := ""
		if !strings.Contains(payload, expectedPayload) {
			t.Errorf("Unexpected payload: %s", payload)
		}
	}))
	defer server.Close()
	sendMetric("gauge", "metricName", 123)
}
