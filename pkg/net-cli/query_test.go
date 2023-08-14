package netcli_test

import (
	"net/http"
	"net/http/httptest"
    "testing"
	netcli "github.com/kirre02/net-cli/pkg/net-cli"
)


var (
    testHandler http.Handler
    testMockServer *httptest.Server
    mockServerURL string 
)

func setupMockServer(handler http.Handler) {
    testMockServer = httptest.NewServer(handler)
    mockServerURL = testMockServer.URL
}

func teardownMockServer() {
    testMockServer.Close()
}

func TestSendPostRequest(t *testing.T) {
    testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status": 200, "Data": "Response data"}`))
    })
    setupMockServer(testHandler)
    defer teardownMockServer()

    response, err := netcli.PostRequest(mockServerURL, []byte(`{"key": "value"}`))
    if err != nil {
        t.Fatalf("Error sending POST request: %v", err)
    }

    // Verify the response fields
    if response.Status == nil || *response.Status != 200 {
        t.Errorf("Expected status code 200, got: %v", response.Status)
    }
    if response.Data != "Response data" {
        t.Errorf("Expected data 'Response data', got: %v", response.Data)
    }

}

func TestGetRequest(t *testing.T) {
    testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status": 200, "Data": "Response data"}`))
    })
    setupMockServer(testHandler)
    defer teardownMockServer()

    response, err := netcli.GetRequest(mockServerURL)
    if err != nil {
        t.Fatalf("Error sending GET request: %v", err)
    }
    // Verify the response fields
    if response.Status == nil || *response.Status != 200 {
        t.Errorf("Expected status code 200, got: %v", response.Status)
    }
    if response.Data != "Response data" {
        t.Errorf("Expected data 'Response data', got: %v", response.Data)
    }
}
