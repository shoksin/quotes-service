package usecase

import (
	"github.com/shoksin/quotes-service/internal/domain"
	"math/rand"
	"strings"
	"time"
)

type QuoteRepository interface {
	Create(quote *domain.Quote) (*domain.Quote, error)
	GetAll() ([]*domain.Quote, error)
	GetByAuthor(author string) ([]*domain.Quote, error)
	GetRandom() (*domain.Quote, error)
	Delete(id int) error
	GetByID(id int) (*domain.Quote, error)
}

type QuoteUseCase struct {
	quoteRepository QuoteRepository
}

func NewQuoteUseCase(quoteRepository QuoteRepository) *QuoteUseCase {
	return &QuoteUseCase{
		quoteRepository: quoteRepository,
	}
}

func (uc *QuoteUseCase) CreateQuote(req *domain.CreateQuoteRequest) (*domain.Quote, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Trim whitespace and normalize the input
	req.Author = strings.TrimSpace(req.Author)
	req.Quote = strings.TrimSpace(req.Quote)

	quote := &domain.Quote{
		Author:    req.Author,
		Quote:     req.Quote,
		CreatedAt: time.Now(),
	}

	return uc.quoteRepository.Create(quote)
}

func (uc *QuoteUseCase) GetAllQuotes() ([]*domain.Quote, error) {
	return uc.quoteRepository.GetAll()
}

func (uc *QuoteUseCase) GetQuotesByAuthor(author string) ([]*domain.Quote, error) {
	if author == "" {
		return nil, domain.ErrInvalidAuthor
	}

	author = strings.TrimSpace(author)
	return uc.quoteRepository.GetByAuthor(author)
}

func (uc *QuoteUseCase) GetRandomQuote() (*domain.Quote, error) {
	quotes, err := uc.quoteRepository.GetAll()
	if err != nil {
		return nil, err
	}

	if len(quotes) == 0 {
		return nil, domain.ErrNoQuotesFound
	}

	randomIndex := rand.Intn(len(quotes))
	return quotes[randomIndex], nil
}

// DeleteQuote deletes a quote by ID
func (uc *QuoteUseCase) DeleteQuote(id int) error {
	if id <= 0 {
		return domain.ErrInvalidID
	}

	return uc.quoteRepository.Delete(id)
}
