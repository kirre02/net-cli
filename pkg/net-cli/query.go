package netcli

import (
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
