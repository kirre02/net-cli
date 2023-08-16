package netcli

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

func (c *DefaultHTTPClient) GetRequest(ctx context.Context, url string) (interface{}, error) {
	client := &http.Client{
		Timeout: c.Timeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseBody interface{}
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (c DefaultHTTPClient) PostRequest(ctx context.Context, url string, payload []byte) (interface{}, error) {
	client := &http.Client{
		Timeout: c.Timeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var resBody interface{}

	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}

func (c *DefaultHTTPClient) PutRequest(ctx context.Context, url string, payload []byte) (interface{}, error) {
	client := &http.Client{
		Timeout: c.Timeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var resBody interface{}

	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}
