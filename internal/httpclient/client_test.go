package httpclient

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestClient_PostRequest(t *testing.T) {

	const requestBody = `{"name":"test","message":"hello"}`

	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %q, want %q", r.Method, http.MethodPost)
		}

		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want %q", got, "application/json")
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading request body: %v", err)
			return
		}

		if string(body) != requestBody {
			t.Errorf("Request body = %q, want %q", body, requestBody)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))

	defer server.Close()

	tests := []struct {
		name       string
		url        string
		body       []byte
		wantStatus int
		wantBody   string
		wantErr    bool
	}{
		{
			name:       "successful request",
			url:        server.URL,
			body:       []byte(requestBody),
			wantStatus: http.StatusOK,
			wantBody:   "success",
		},
		{
			name:    "missing scheme",
			url:     "example.com",
			body:    []byte(requestBody),
			wantErr: true,
		},
		{
			name:    "unsupported scheme",
			url:     "ftp://example.com",
			body:    []byte(requestBody),
			wantErr: true,
		},
	}

	client := NewClient(2 * time.Second)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			headers.Set("Content-Type", "application/json")

			got, err := client.PostRequestWithHeaders(
				context.Background(),
				tt.url,
				tt.body,
				headers,
			)

			if (err != nil) != tt.wantErr {
				t.Fatalf("PostRequest() error = %v, wantErr %v", err, tt.wantErr)
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

func TestClient_PutRequest(t *testing.T) {
	const requestBody = `{"name":"test","message":"hello"}`

	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %q, want %q", r.Method, http.MethodPut)
		}

		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Errorf("Content-Type = %q, want %q", got, "application/json")
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("reading request body: %v", err)
			return
		}

		if string(body) != requestBody {
			t.Errorf("Request body = %q, want %q", body, requestBody)
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))
	defer server.Close()

	tests := []struct {
		name       string
		url        string
		body       []byte
		wantStatus int
		wantBody   string
		wantErr    bool
	}{
		{
			name:       "successful request",
			url:        server.URL,
			body:       []byte(requestBody),
			wantStatus: http.StatusOK,
			wantBody:   "success",
		},
		{
			name:    "missing scheme",
			url:     "example.com",
			body:    []byte(requestBody),
			wantErr: true,
		},
		{
			name:    "unsupported scheme",
			url:     "ftp://example.com",
			body:    []byte(requestBody),
			wantErr: true,
		},
	}

	client := NewClient(2 * time.Second)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			headers.Set("Content-Type", "application/json")

			got, err := client.PutRequestWithHeaders(
				context.Background(),
				tt.url,
				tt.body,
				headers,
			)

			if (err != nil) != tt.wantErr {
				t.Fatalf("PutRequestWithHeaders() error = %v, wantErr %v", err, tt.wantErr)
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

func TestClient_DeleteRequest(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		statusCode int
		wantErr    bool
	}{
		{
			name:       "successful request",
			statusCode: http.StatusOK,
		},
		{
			name:    "missing scheme",
			url:     "localhost:8080/resource",
			wantErr: true,
		},
		{
			name:    "unsupported scheme",
			url:     "ftp://example.com/resource",
			wantErr: true,
		},
		{
			name:       "not found",
			statusCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					if r.Method != http.MethodDelete {
						t.Errorf("expected method DELETE, got %s", r.Method)
					}

					if r.Header.Get("User-Agent") != "netcli/0.1.0" {
						t.Errorf(
							"expected User-Agent netcli/0.1.0, got %q",
							r.Header.Get("User-Agent"),
						)
					}

					if tt.statusCode != 0 {
						w.WriteHeader(tt.statusCode)
					}

					_, _ = w.Write([]byte(`{"message":"test response"}`))
				},
			))
			defer server.Close()

			url := tt.url
			if url == "" {
				url = server.URL
			}

			client := NewClient(0)
			res, err := client.DeleteRequest(context.Background(), url)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if res == nil {
				t.Fatal("expected a response, got nil")
			}

			expectedStatus := tt.statusCode
			if expectedStatus == 0 {
				expectedStatus = http.StatusOK
			}

			if res.StatusCode != expectedStatus {
				t.Errorf(
					"expected status code %d, got %d",
					expectedStatus,
					res.StatusCode,
				)
			}

			if !strings.Contains(
				string(res.Body),
				"test response",
			) {
				t.Errorf("unexpected response body: %s", res.Body)
			}
		})
	}
}
