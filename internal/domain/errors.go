package domain

import "errors"

var (
	ErrQuoteNotFound = errors.New("quote not found")
	ErrInvalidAuthor = errors.New("invalid author")
	ErrInvalidQuote  = errors.New("invalid quote")
	ErrInvalidID     = errors.New("invalid quote ID")
	ErrNoQuotesFound = errors.New("no quotes found")
)
