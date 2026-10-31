package storage

import (
	"context"
	"errors"
	"testing"

	"github.com/kirre02/net-cli/internal/domain"
)

func newTestSQLite(t *testing.T) *SQLite {
	t.Helper()

	dbPath := t.TempDir() + "/test.db"

	db, err := openSQLite(dbPath)
	if err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	})

	return db
}

func TestCollectionCRUD(t *testing.T) {
	store := newTestSQLite(t)
	ctx := context.Background()

	collection := domain.Collection{
		ID:   "collection-1",
		Name: "My collection",
	}

	// Create
	if err := store.CreateCollection(ctx, collection); err != nil {
		t.Fatalf("CreateCollection() error = %v", err)
	}

	// Get
	got, err := store.GetCollection(ctx, collection.ID)
	if err != nil {
		t.Fatalf("GetCollection() error = %v", err)
	}

	if got != collection {
		t.Errorf("GetCollection() = %+v, want %+v", got, collection)
	}

	// Update
	collection.Name = "Updated collection"

	if err := store.UpdateCollection(ctx, collection); err != nil {
		t.Fatalf("UpdateCollection() error = %v", err)
	}

	got, err = store.GetCollection(ctx, collection.ID)
	if err != nil {
		t.Fatalf("GetCollection() after update error = %v", err)
	}

	if got != collection {
		t.Errorf("GetCollection() after update = %+v, want %+v", got, collection)
	}

	// Delete
	if err := store.DeleteCollection(ctx, collection.ID); err != nil {
		t.Fatalf("DeleteCollection() error = %v", err)
	}

	_, err = store.GetCollection(ctx, collection.ID)
	if !errors.Is(err, errorCollectionNotFound) {
		t.Errorf("GetCollection() error = %v, want %v", err, errorCollectionNotFound)
	}
}

func TestCreateCollectionValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.Collection
		wantErr error
	}{
		{
			name: "empty name",
			input: domain.Collection{
				ID:   "collection-1",
				Name: "",
			},
			wantErr: domain.ErrEmptyCollectionName,
		},
		{
			name: "whitespace name",
			input: domain.Collection{
				ID:   "collection-1",
				Name: "   ",
			},
			wantErr: domain.ErrEmptyCollectionName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := newTestSQLite(t)
			ctx := context.Background()

			err := store.CreateCollection(ctx, tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateCollection() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateCollectionDuplicateName(t *testing.T) {
	store := newTestSQLite(t)
	ctx := context.Background()

	first := domain.Collection{
		ID:   "collection-1",
		Name: "My collection",
	}

	second := domain.Collection{
		ID:   "collection-2",
		Name: "My collection",
	}

	if err := store.CreateCollection(ctx, first); err != nil {
		t.Fatalf("CreateCollection() first error = %v", err)
	}

	err := store.CreateCollection(ctx, second)
	if err == nil {
		t.Fatal("CreateCollection() expected an error for duplicate name")
	}

	collections, err := store.GetCollections(ctx)
	if err != nil {
		t.Fatalf("GetCollections() error = %v", err)
	}

	if len(collections) != 1 {
		t.Errorf("GetCollections() returned %d collections, want 1", len(collections))
	}
}

func TestDeleteCollectionCascadesRequests(t *testing.T) {
	store := newTestSQLite(t)
	ctx := context.Background()

	collection := domain.Collection{
		ID:   "collection-1",
		Name: "My collection",
	}

	request := domain.Request{
		ID:           "request-1",
		CollectionID: collection.ID,
		Name:         "My request",
		Method:       domain.MethodGet,
		URL:          "https://example.com",
	}

	if err := store.CreateCollection(ctx, collection); err != nil {
		t.Fatalf("CreateCollection() error = %v", err)
	}

	if err := store.CreateRequest(ctx, request); err != nil {
		t.Fatalf("CreateRequest() error = %v", err)
	}

	if err := store.DeleteCollection(ctx, collection.ID); err != nil {
		t.Fatalf("DeleteCollection() error = %v", err)
	}

	_, err := store.GetRequest(ctx, request.ID)
	if !errors.Is(err, errorRequestNotFound) {
		t.Errorf("GetRequest() error = %v, want %v", err, errorRequestNotFound)
	}
}
