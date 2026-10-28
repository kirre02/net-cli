package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

// Client executes HTTP requests with a configured timeout and User-Agent header.
type Client struct {
	client    *http.Client
	userAgent string
}

// Response contains the data returned by an HTTP request.
type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// do executes an HTTP request using ctx. A canceled context or expired client
// timeout causes an error. Transport, request-construction, and body-reading
// failures also return an error. Non-2xx HTTP status codes return a Response
// with a nil error.
func (c *Client) do(
	ctx context.Context,
	method string,
	url string,
	body []byte,
) (*Response, error) {

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: res.StatusCode,
		Header:     res.Header.Clone(),
		Body:       resBody,
	}, nil
}

// NewClient returns a Client with the supplied request timeout and User-Agent header.
//
// A timeout of zero or less uses a default of 10 seconds. An empty userAgent uses
// "netcli/0.1.0". The timeout includes the entire request, including reading the
// response body.
func NewClient(timeout time.Duration, userAgent string) *Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	if userAgent == "" {
		userAgent = "netcli/0.1.0"
	}

	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
		userAgent: userAgent,
	}
}
