package domain

import (
	"time"
)

type Quote struct {
	ID        int       `json:"id" db:"id"`
	Author    string    `json:"author" db:"author"`
	Quote     string    `json:"quote" db:"quote"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateQuoteRequest struct {
	Author string `json:"author" db:"author"`
	Quote  string `json:"quote" db:"quote"`
}

func (r *CreateQuoteRequest) Validate() error {
	if r.Author == "" {
		return ErrInvalidAuthor
	}
	if r.Quote == "" {
		return ErrInvalidQuote
	}
	return nil
}
