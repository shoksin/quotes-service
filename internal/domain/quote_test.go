package domain

import (
	"errors"
	"testing"
)

func TestCreateQuoteRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateQuoteRequest
		wantErr error
	}{
		{
			name:    "valid request",
			req:     CreateQuoteRequest{Author: "Albert Einstein", Quote: "Life is like riding a bicycle."},
			wantErr: nil,
		},
		{
			name:    "empty author",
			req:     CreateQuoteRequest{Author: "", Quote: "Some quote"},
			wantErr: ErrInvalidAuthor,
		},
		{
			name:    "empty quote",
			req:     CreateQuoteRequest{Author: "Author", Quote: ""},
			wantErr: ErrInvalidQuote,
		},
		{
			name:    "both empty",
			req:     CreateQuoteRequest{Author: "", Quote: ""},
			wantErr: ErrInvalidAuthor, //before ErrInvalidQuote
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
