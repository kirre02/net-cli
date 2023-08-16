package netcli

import (
	"context"
	"time"
)

type QueryClient interface {
	GetRequest(ctx context.Context, url string) (interface{}, error)
	DeleteRequest(ctx context.Context, url string) (interface{}, error)
	PostRequest(ctx context.Context, url string, payload []byte) (interface{}, error)
    PutRequest(ctx context.Context, url string, payload []byte) (interface{}, error)
}

type DefaultHTTPClient struct {
	Timeout   time.Duration // Maximum time for the HTTP request to complete
	UserAgent string        // User-Agent header value for the requests
}

func NewDefaultHTTPClient() *DefaultHTTPClient {
	return &DefaultHTTPClient{
		Timeout:   10 * time.Second,
		UserAgent: "netcli/0.1.0",
	}
}
