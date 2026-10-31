package domain

import "testing"

func TestCollectionValidate(t *testing.T) {
	tests := []struct {
		name       string
		collection Collection
		wantErr    error
	}{
		{
			name: "valid collection",
			collection: Collection{
				ID:   "collection-1",
				Name: "My collection",
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			collection: Collection{
				ID:   "collection-1",
				Name: "",
			},
			wantErr: ErrEmptyCollectionName,
		},
		{
			name: "whitespace name",
			collection: Collection{
				ID:   "collection-1",
				Name: "   ",
			},
			wantErr: ErrEmptyCollectionName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.collection.Validate()

			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
