package netcli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetRequest(url string) (*HTTPResponse, error) {
    response, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("GET request failed: %v", err)
    }
    
    defer response.Body.Close()

    var responseBody HTTPResponse

    err = json.NewDecoder(response.Body).Decode(&responseBody)
    if err != nil {
        return nil, fmt.Errorf("Error decoding JSON: %v", err)
    }

    return &responseBody, nil
}

func PostRequest(url string, payload []byte) (*HTTPResponse, error) {
    resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
    if err != nil {
        return nil, err
    }

    defer resp.Body.Close()

    var response HTTPResponse

    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, err
    }

    if response.Error != nil {
        return nil, fmt.Errorf("Error in response: %d - %s", response.Error.Status, response.Error.Message)
    }

    return &response, nil
}
