package service

import (
	"context"
	"github.com/shoksin/quotes-service/internal/models"
)

type QuoteRepository interface {
	CreateQuote(ctx context.Context, quote *models.Quote) (*models.Quote, error)
	GetQuoteByAuthor(author string) (*models.Quote, error)
	GetAllQuotes() ([]*models.Quote, error)
	GetRandomQuote() (*models.Quote, error)
	DeleteQuoteByID(id int) (bool, error)
}
type QuoteService struct {
	repository QuoteRepository
}

func NewService(repository QuoteRepository) *QuoteService {
	return &QuoteService{
		repository: repository,
	}
}

func (s *QuoteService) GetQuoteByAuthor(author string) (models.Quote, error) {}
func (s *QuoteService) CreateQuote()                                         {}
