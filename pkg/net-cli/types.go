package netcli

import "fmt"

type HTTPResponse struct {
    Error *HTTPError `json:"error"`
    Status *int `json:"status"`
    Data string `json:"Data"`
}

type HTTPError struct {
    Status  int `json:"error"` // HTTP status code
    Message string `json:"message"` // a user friendly error message
    Cause   error `json:"-"` // The underlying cause of error
}

func (e HTTPError) Error() string {
    return fmt.Sprintf("HTTP Error: %d - %s. cause %v", e.Status, e.Message, e.Cause)
}
