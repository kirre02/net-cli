package httpclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	rawUrl string,
	body []byte,
	headers http.Header,
) (*Response, error) {
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid URL %q: %w", rawUrl, err)
	}

	if (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") ||
		parsedURL.Host == "" {
		return nil, fmt.Errorf(
			"invalid URL %q: provide a complete HTTP or HTTPS URL",
			rawUrl,
		)
	}

	req, err := http.NewRequestWithContext(ctx, method, rawUrl, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating HTTP request: %w", err)
	}

	req.Header.Set("User-Agent", c.userAgent)

	// Apply any headers provided by the caller.
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"executing %s request to %q: %w",
			method,
			rawUrl,
			err,
		)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"reading response body from %s request to %q: %w",
			method,
			rawUrl,
			err,
		)
	}

	return &Response{
		StatusCode: res.StatusCode,
		Header:     res.Header.Clone(),
		Body:       resBody,
	}, nil
}

func (c *Client) GetRequest(
	ctx context.Context,
	rawUrl string,
) (*Response, error) {
	return c.do(ctx, http.MethodGet, rawUrl, nil, nil)
}

func (c *Client) PostRequest(
	ctx context.Context,
	rawUrl string,
	body []byte,
) (*Response, error) {
	return c.do(ctx, http.MethodPost, rawUrl, body, nil)
}

func (c *Client) PostRequestWithHeaders(
	ctx context.Context,
	rawURL string,
	body []byte,
	headers http.Header,
) (*Response, error) {
	return c.do(ctx, http.MethodPost, rawURL, body, headers)
}

func (c *Client) PutRequest(
	ctx context.Context,
	rawUrl string,
	body []byte,
) (*Response, error) {
	return c.do(ctx, http.MethodPut, rawUrl, body, nil)
}

func (c *Client) PutRequestWithHeaders(
	context context.Context,
	rawURL string,
	payload []byte,
	headers http.Header) (*Response, error) {
	return c.do(context, http.MethodPut, rawURL, payload, headers)
}

func (c *Client) DeleteRequest(
	ctx context.Context,
	rawUrl string,
) (*Response, error) {
	return c.do(ctx, http.MethodDelete, rawUrl, nil, nil)
}

// NewClient returns a Client with the supplied request timeout and default
// User-Agent header.
//
// A timeout of zero or less uses a default of 10 seconds. The timeout includes
// the entire request, including reading the response body.
func NewClient(timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &Client{
		client: &http.Client{
			Timeout: timeout,
		},
		userAgent: "netcli/0.1.0",
	}
}
