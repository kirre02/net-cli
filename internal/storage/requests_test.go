package storage

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kirre02/net-cli/internal/domain"
)

func TestRequestCRUD(t *testing.T) {
	store := newTestSQLite(t)
	ctx := context.Background()

	collection := domain.Collection{
		ID:   "collection-1",
		Name: "My collection",
	}

	if err := store.CreateCollection(ctx, collection); err != nil {
		t.Fatalf("CreateCollection() error = %v", err)
	}

	request := domain.Request{
		ID:           "request-1",
		CollectionID: collection.ID,
		Name:         "Get users",
		Method:       domain.MethodGet,
		URL:          "https://example.com/users",
		Headers: map[string]string{
			"Authorization": "Bearer token",
			"Accept":        "application/json",
		},
		Body: "",
	}

	// Create
	if err := store.CreateRequest(ctx, request); err != nil {
		t.Fatalf("CreateRequest() error = %v", err)
	}

	// Get
	got, err := store.GetRequest(ctx, request.ID)
	if err != nil {
		t.Fatalf("GetRequest() error = %v", err)
	}

	if !reflect.DeepEqual(got, request) {
		t.Errorf("GetRequest() = %+v, want %+v", got, request)
	}

	// Update
	request.Name = "Updated request"
	request.Method = domain.MethodPost
	request.URL = "https://example.com/users"
	request.Headers["Content-Type"] = "application/json"
	request.Body = `{"name":"John"}`

	if err := store.UpdateRequest(ctx, request); err != nil {
		t.Fatalf("UpdateRequest() error = %v", err)
	}

	got, err = store.GetRequest(ctx, request.ID)
	if err != nil {
		t.Fatalf("GetRequest() after update error = %v", err)
	}

	if !reflect.DeepEqual(got, request) {
		t.Errorf("GetRequest() after update = %+v, want %+v", got, request)
	}

	// Delete
	if err := store.DeleteRequest(ctx, request.ID); err != nil {
		t.Fatalf("DeleteRequest() error = %v", err)
	}

	_, err = store.GetRequest(ctx, request.ID)
	if !errors.Is(err, errorRequestNotFound) {
		t.Errorf("GetRequest() error = %v, want %v", err, errorRequestNotFound)
	}
}

func TestCreateRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request domain.Request
		wantErr error
	}{
		{
			name: "empty name",
			request: domain.Request{
				ID:           "request-1",
				CollectionID: "collection-1",
				Name:         "",
				Method:       domain.MethodGet,
				URL:          "https://example.com",
			},
			wantErr: domain.ErrEmptyRequestName,
		},
		{
			name: "whitespace name",
			request: domain.Request{
				ID:           "request-1",
				CollectionID: "collection-1",
				Name:         "   ",
				Method:       domain.MethodGet,
				URL:          "https://example.com",
			},
			wantErr: domain.ErrEmptyRequestName,
		},
		{
			name: "empty URL",
			request: domain.Request{
				ID:           "request-1",
				CollectionID: "collection-1",
				Name:         "Get users",
				Method:       domain.MethodGet,
				URL:          "",
			},
			wantErr: domain.ErrEmptyRequestURL,
		},
		{
			name: "whitespace URL",
			request: domain.Request{
				ID:           "request-1",
				CollectionID: "collection-1",
				Name:         "Get users",
				Method:       domain.MethodGet,
				URL:          "   ",
			},
			wantErr: domain.ErrEmptyRequestURL,
		},
		{
			name: "invalid method",
			request: domain.Request{
				ID:           "request-1",
				CollectionID: "collection-1",
				Name:         "Get users",
				Method:       "INVALID",
				URL:          "https://example.com",
			},
			wantErr: domain.ErrInvalidMethod,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newTestSQLite(t)
			ctx := context.Background()

			collection := domain.Collection{
				ID:   "collection-1",
				Name: "My collection",
			}

			if err := store.CreateCollection(ctx, collection); err != nil {
				t.Fatalf("CreateCollection() error = %v", err)
			}

			err := store.CreateRequest(ctx, tt.request)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateRequest() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetRequestsByCollection(t *testing.T) {
	store := newTestSQLite(t)
	ctx := context.Background()

	collection := domain.Collection{
		ID:   "collection-1",
		Name: "My collection",
	}

	if err := store.CreateCollection(ctx, collection); err != nil {
		t.Fatalf("CreateCollection() error = %v", err)
	}

	requests := []domain.Request{
		{
			ID:           "request-2",
			CollectionID: collection.ID,
			Name:         "Create user",
			Method:       domain.MethodPost,
			URL:          "https://example.com/users",
			Body:         `{"name":"John"}`,
		},
		{
			ID:           "request-1",
			CollectionID: collection.ID,
			Name:         "Get users",
			Method:       domain.MethodGet,
			URL:          "https://example.com/users",
		},
	}

	for _, request := range requests {
		if err := store.CreateRequest(ctx, request); err != nil {
			t.Fatalf("CreateRequest() error = %v", err)
		}
	}

	got, err := store.GetRequestsByCollection(ctx, collection.ID)
	if err != nil {
		t.Fatalf("GetRequestsByCollection() error = %v", err)
	}

	if len(got) != len(requests) {
		t.Fatalf(
			"GetRequestsByCollection() returned %d requests, want %d",
			len(got),
			len(requests),
		)
	}

	for i, request := range requests {
		if !reflect.DeepEqual(got[i], request) {
			t.Errorf(
				"GetRequestsByCollection()[%d] = %+v, want %+v",
				i,
				got[i],
				request,
			)
		}
	}
}

func TestCreateRequestRequiresExistingCollection(t *testing.T) {
	store := newTestSQLite(t)
	ctx := context.Background()

	request := domain.Request{
		ID:           "request-1",
		CollectionID: "nonexistent-collection",
		Name:         "Get users",
		Method:       domain.MethodGet,
		URL:          "https://example.com/users",
	}

	err := store.CreateRequest(ctx, request)
	if err == nil {
		t.Fatal("CreateRequest() expected an error for nonexistent collection")
	}
}
