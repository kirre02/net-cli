package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client    *http.Client
	userAgent string
}

type Response struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

func (c *Client) do(
	ctx context.Context,
	method string,
	url string,
	body []byte,
) (*Response, error) {
	// Build the request

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	// Execute it using c.client.
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Return status code, headers, and body.
	return &Response{
		StatusCode: res.StatusCode,
		Header:     res.Header.Clone(),
		Body:       resBody,
	}, nil
}

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
