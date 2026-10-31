package domain

import (
	"errors"
	"net/http"
	"strings"
)

type HTTPMethod string

const (
	MethodGet     HTTPMethod = http.MethodGet
	MethodPost    HTTPMethod = http.MethodPost
	MethodPut     HTTPMethod = http.MethodPut
	MethodPatch   HTTPMethod = http.MethodPatch
	MethodDelete  HTTPMethod = http.MethodDelete
	MethodHead    HTTPMethod = http.MethodHead
	MethodOptions HTTPMethod = http.MethodOptions
)

var (
	ErrEmptyRequestName = errors.New("request name cannot be empty")
	ErrEmptyRequestURL  = errors.New("request URL cannot be empty")
	ErrInvalidMethod    = errors.New("invalid HTTP method")
)

type Request struct {
	ID           string
	CollectionID string
	Name         string
	Method       HTTPMethod
	URL          string
	Headers      map[string]string
	Body         string
}

func (r Request) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return ErrEmptyRequestName
	}

	if strings.TrimSpace(r.URL) == "" {
		return ErrEmptyRequestURL
	}

	switch r.Method {
	case MethodGet, MethodPost, MethodPut, MethodPatch,
		MethodDelete, MethodHead, MethodOptions:
		return nil
	default:
		return ErrInvalidMethod
	}
}
