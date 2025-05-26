package repository

import (
	"database/sql"
	"fmt"
	"github.com/shoksin/quotes-service/internal/domain"

	_ "github.com/lib/pq"
)

type QuoteRepository struct {
	db *sql.DB
}

func NewQuoteRepository(db *sql.DB) *QuoteRepository {
	return &QuoteRepository{db: db}
}

// Create creates a new quote
func (r *QuoteRepository) Create(quote *domain.Quote) (*domain.Quote, error) {
	query := `INSERT INTO quotes (author, quote) VALUES ($1, $2) RETURNING id, created_at`

	row := r.db.QueryRow(query, quote.Author, quote.Quote)

	err := row.Scan(&quote.ID, &quote.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create quote: %w", err)
	}

	return quote, nil
}

// GetAll returns all quotes
func (r *QuoteRepository) GetAll() ([]*domain.Quote, error) {
	query := `SELECT id, author, quote, created_at FROM quotes ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes: %w", err)
	}
	defer rows.Close()

	var quotes []*domain.Quote
	for rows.Next() {
		quote := &domain.Quote{}
		err = rows.Scan(&quote.ID, &quote.Author, &quote.Quote, &quote.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, quote)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return quotes, nil
}

// GetByAuthor returns quotes by author
func (r *QuoteRepository) GetByAuthor(author string) ([]*domain.Quote, error) {
	query := `SELECT id, author, quote, created_at FROM quotes WHERE author = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, author)
	if err != nil {
		return nil, fmt.Errorf("failed to get quotes by author: %w", err)
	}
	defer rows.Close()

	var quotes []*domain.Quote
	for rows.Next() {
		quote := &domain.Quote{}
		err = rows.Scan(&quote.ID, &quote.Author, &quote.Quote, &quote.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}
		quotes = append(quotes, quote)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return quotes, nil
}

// GetRandom returns a random quote
func (r *QuoteRepository) GetRandom() (*domain.Quote, error) {
	query := `SELECT id, author, quote, created_at FROM quotes ORDER BY RANDOM() LIMIT 1`

	quote := &domain.Quote{}
	row := r.db.QueryRow(query)

	err := row.Scan(&quote.ID, &quote.Author, &quote.Quote, &quote.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNoQuotesFound
		}
		return nil, fmt.Errorf("failed to get random quote: %w", err)
	}

	return quote, nil
}

// GetByID returns a quote by ID
func (r *QuoteRepository) GetByID(id int) (*domain.Quote, error) {
	query := `SELECT id, author, quote, created_at FROM quotes WHERE id = $1`

	quote := &domain.Quote{}
	row := r.db.QueryRow(query, id)

	err := row.Scan(&quote.ID, &quote.Author, &quote.Quote, &quote.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrQuoteNotFound
		}
		return nil, fmt.Errorf("failed to get quote by ID: %w", err)
	}

	return quote, nil
}

// Delete deletes a quote by ID
func (r *QuoteRepository) Delete(id int) error {
	query := `DELETE FROM quotes WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete quote: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrQuoteNotFound
	}

	return nil
}
