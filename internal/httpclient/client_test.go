package httpclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_GetRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))
	defer server.Close()

	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantBody   string
		wantErr    bool
	}{
		{
			name:       "successful request",
			url:        server.URL,
			wantStatus: http.StatusOK,
			wantBody:   "success",
		},
		{
			name:    "missing scheme",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "unsupported scheme",
			url:     "ftp://example.com",
			wantErr: true,
		},
		{
			name:       "not found",
			url:        server.URL + "/missing",
			wantStatus: http.StatusOK,
			wantBody:   "success",
		},
	}

	client := NewClient(2 * time.Second)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetRequest(context.Background(), tt.url)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetRequest() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			if got.StatusCode != tt.wantStatus {
				t.Errorf("StatusCode = %d, want %d", got.StatusCode, tt.wantStatus)
			}

			if string(got.Body) != tt.wantBody {
				t.Errorf("Body = %q, want %q", got.Body, tt.wantBody)
			}
		})
	}
}
