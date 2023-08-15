package netcli_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	netcli "github.com/kirre02/net-cli/pkg/net-cli"
)

var (
	testHandler    http.Handler
	testMockServer *httptest.Server
	mockServerURL  string
)

func setupMockServer(handler http.Handler) {
	testMockServer = httptest.NewServer(handler)
	mockServerURL = testMockServer.URL
}

func teardownMockServer() {
	testMockServer.Close()
}

func TestGetRequest(t *testing.T) {
	testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": 200, "body": {"key": "value"}}`))
	})
	setupMockServer(testHandler)
	defer teardownMockServer()

	client := netcli.NewDefaultHTTPClient()

	response, err := client.GetRequest(context.Background(), mockServerURL)
	if err != nil {
		t.Fatalf("Error sending GET request: %v", err)
	}

	// Validate the response type
	responseObject, ok := response.(map[string]interface{})
	if !ok {
		t.Fatalf("Unexpected response type: %T", response)
	}

	// Validate specific fields in the response object
	if value, exists := responseObject["key"]; !exists {
		t.Errorf("Expected key 'key' in response body")
	} else if value != "value" {
		t.Errorf("Expected value 'value', got: %v", value)
	}
}

func TestPostRequest(t *testing.T) {
	testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": 200, "body": {"key": "value"}}`))
	})
	setupMockServer(testHandler)
	defer teardownMockServer()

	client := netcli.NewDefaultHTTPClient()

	payload := []byte(`{"key": "value"}`)
	response, err := client.PostRequest(context.Background(), mockServerURL, payload)
	if err != nil {
		t.Fatalf("Error sending POST request: %v", err)
	}

	// Validate the response type
	responseObject, ok := response.(map[string]interface{})
	if !ok {
		t.Fatalf("Unexpected response type: %T", response)
	}

	// Validate specific fields in the response object
	if value, exists := responseObject["key"]; !exists {
		t.Errorf("Expected key 'key' in response body")
	} else if value != "value" {
		t.Errorf("Expected value 'value', got: %v", value)
	}
}
