package domain

import "testing"

func TestRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		request Request
		wantErr error
	}{
		{
			name: "valid request",
			request: Request{
				ID:           "request-1",
				CollectionID: "collection-1",
				Name:         "Get users",
				Method:       MethodGet,
				URL:          "https://example.com/users",
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			request: Request{
				Name:   "",
				Method: MethodGet,
				URL:    "https://example.com",
			},
			wantErr: ErrEmptyRequestName,
		},
		{
			name: "whitespace name",
			request: Request{
				Name:   "   ",
				Method: MethodGet,
				URL:    "https://example.com",
			},
			wantErr: ErrEmptyRequestName,
		},
		{
			name: "empty URL",
			request: Request{
				Name:   "Get users",
				Method: MethodGet,
				URL:    "",
			},
			wantErr: ErrEmptyRequestURL,
		},
		{
			name: "whitespace URL",
			request: Request{
				Name:   "Get users",
				Method: MethodGet,
				URL:    "   ",
			},
			wantErr: ErrEmptyRequestURL,
		},
		{
			name: "invalid method",
			request: Request{
				Name:   "Get users",
				Method: "INVALID",
				URL:    "https://example.com",
			},
			wantErr: ErrInvalidMethod,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()

			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
