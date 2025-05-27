package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/shoksin/quotes-service/internal/domain"
)

// MockQuoteRepository is a mock implementation of QuoteRepository interface
type MockQuoteRepository struct {
	CreateFunc      func(quote *domain.Quote) (*domain.Quote, error)
	GetAllFunc      func() ([]*domain.Quote, error)
	GetByAuthorFunc func(author string) ([]*domain.Quote, error)
	GetRandomFunc   func() (*domain.Quote, error)
	DeleteFunc      func(id int) error
	GetByIDFunc     func(id int) (*domain.Quote, error)
}

func (m *MockQuoteRepository) Create(quote *domain.Quote) (*domain.Quote, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(quote)
	}
	return nil, nil
}

func (m *MockQuoteRepository) GetAll() ([]*domain.Quote, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil, nil
}

func (m *MockQuoteRepository) GetByAuthor(author string) ([]*domain.Quote, error) {
	if m.GetByAuthorFunc != nil {
		return m.GetByAuthorFunc(author)
	}
	return nil, nil
}

func (m *MockQuoteRepository) GetRandom() (*domain.Quote, error) {
	if m.GetRandomFunc != nil {
		return m.GetRandomFunc()
	}
	return nil, nil
}

func (m *MockQuoteRepository) Delete(id int) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockQuoteRepository) GetByID(id int) (*domain.Quote, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func TestQuoteUseCase_CreateQuote(t *testing.T) {
	tests := []struct {
		name          string
		request       *domain.CreateQuoteRequest
		mockFunc      func(quote *domain.Quote) (*domain.Quote, error)
		expectedError error
		checkResult   func(t *testing.T, result *domain.Quote)
	}{
		{
			name: "successful quote creation",
			request: &domain.CreateQuoteRequest{
				Author: "  Test Author  ",
				Quote:  "  Test Quote  ",
			},
			mockFunc: func(quote *domain.Quote) (*domain.Quote, error) {
				// Check that input was trimmed
				if quote.Author != "Test Author" || quote.Quote != "Test Quote" {
					t.Errorf("expected trimmed values, got Author: '%s', Quote: '%s'", quote.Author, quote.Quote)
				}
				quote.ID = 1
				return quote, nil
			},
			checkResult: func(t *testing.T, result *domain.Quote) {
				if result.ID != 1 {
					t.Errorf("expected ID 1, got %d", result.ID)
				}
				if result.Author != "Test Author" {
					t.Errorf("expected Author 'Test Author', got '%s'", result.Author)
				}
				if result.Quote != "Test Quote" {
					t.Errorf("expected Quote 'Test Quote', got '%s'", result.Quote)
				}
			},
		},
		{
			name: "validation error - empty author",
			request: &domain.CreateQuoteRequest{
				Author: "",
				Quote:  "Test Quote",
			},
			expectedError: domain.ErrInvalidAuthor,
		},
		{
			name: "validation error - empty quote",
			request: &domain.CreateQuoteRequest{
				Author: "Test Author",
				Quote:  "",
			},
			expectedError: domain.ErrInvalidQuote,
		},
		{
			name: "repository error",
			request: &domain.CreateQuoteRequest{
				Author: "Test Author",
				Quote:  "Test Quote",
			},
			mockFunc: func(quote *domain.Quote) (*domain.Quote, error) {
				return nil, errors.New("database error")
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockQuoteRepository{
				CreateFunc: tt.mockFunc,
			}
			useCase := NewQuoteUseCase(mockRepo)

			result, err := useCase.CreateQuote(tt.request)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

func TestQuoteUseCase_GetAllQuotes(t *testing.T) {
	tests := []struct {
		name          string
		mockFunc      func() ([]*domain.Quote, error)
		expectedError error
		expectedLen   int
	}{
		{
			name: "successful get all quotes",
			mockFunc: func() ([]*domain.Quote, error) {
				return []*domain.Quote{
					{ID: 1, Author: "Author1", Quote: "Quote1"},
					{ID: 2, Author: "Author2", Quote: "Quote2"},
				}, nil
			},
			expectedLen: 2,
		},
		{
			name: "empty result",
			mockFunc: func() ([]*domain.Quote, error) {
				return []*domain.Quote{}, nil
			},
			expectedLen: 0,
		},
		{
			name: "repository error",
			mockFunc: func() ([]*domain.Quote, error) {
				return nil, errors.New("database error")
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockQuoteRepository{
				GetAllFunc: tt.mockFunc,
			}
			useCase := NewQuoteUseCase(mockRepo)

			result, err := useCase.GetAllQuotes()

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(result) != tt.expectedLen {
				t.Errorf("expected %d quotes, got %d", tt.expectedLen, len(result))
			}
		})
	}
}

func TestQuoteUseCase_GetQuotesByAuthor(t *testing.T) {
	tests := []struct {
		name          string
		author        string
		mockFunc      func(author string) ([]*domain.Quote, error)
		expectedError error
		expectedLen   int
	}{
		{
			name:   "successful get quotes by author",
			author: "  Test Author  ",
			mockFunc: func(author string) ([]*domain.Quote, error) {
				// Check that author was trimmed
				if author != "Test Author" {
					t.Errorf("expected trimmed author 'Test Author', got '%s'", author)
				}
				return []*domain.Quote{
					{ID: 1, Author: "Test Author", Quote: "Quote1"},
					{ID: 2, Author: "Test Author", Quote: "Quote2"},
				}, nil
			},
			expectedLen: 2,
		},
		{
			name:          "empty author",
			author:        "",
			expectedError: domain.ErrInvalidAuthor,
		},
		{
			name:   "whitespace only author",
			author: "   ",
			mockFunc: func(author string) ([]*domain.Quote, error) {
				// Current implementation trims after checking for empty,
				// so whitespace-only strings pass through
				if author != "" {
					t.Errorf("expected empty string after trim, got '%s'", author)
				}
				return []*domain.Quote{}, nil
			},
			expectedLen: 0,
		},
		{
			name:   "repository error",
			author: "Test Author",
			mockFunc: func(author string) ([]*domain.Quote, error) {
				return nil, errors.New("database error")
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockQuoteRepository{
				GetByAuthorFunc: tt.mockFunc,
			}
			useCase := NewQuoteUseCase(mockRepo)

			result, err := useCase.GetQuotesByAuthor(tt.author)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(result) != tt.expectedLen {
				t.Errorf("expected %d quotes, got %d", tt.expectedLen, len(result))
			}
		})
	}
}

func TestQuoteUseCase_GetRandomQuote(t *testing.T) {
	tests := []struct {
		name          string
		mockFunc      func() ([]*domain.Quote, error)
		expectedError error
		checkResult   func(t *testing.T, result *domain.Quote)
	}{
		{
			name: "successful get random quote",
			mockFunc: func() ([]*domain.Quote, error) {
				return []*domain.Quote{
					{ID: 1, Author: "Author1", Quote: "Quote1"},
					{ID: 2, Author: "Author2", Quote: "Quote2"},
					{ID: 3, Author: "Author3", Quote: "Quote3"},
				}, nil
			},
			checkResult: func(t *testing.T, result *domain.Quote) {
				if result == nil {
					t.Error("expected non-nil quote")
				}
				if result.ID < 1 || result.ID > 3 {
					t.Errorf("unexpected quote ID: %d", result.ID)
				}
			},
		},
		{
			name: "no quotes found",
			mockFunc: func() ([]*domain.Quote, error) {
				return []*domain.Quote{}, nil
			},
			expectedError: domain.ErrNoQuotesFound,
		},
		{
			name: "repository error",
			mockFunc: func() ([]*domain.Quote, error) {
				return nil, errors.New("database error")
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockQuoteRepository{
				GetAllFunc: tt.mockFunc,
			}
			useCase := NewQuoteUseCase(mockRepo)

			result, err := useCase.GetRandomQuote()

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

func TestQuoteUseCase_DeleteQuote(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		mockFunc      func(id int) error
		expectedError error
	}{
		{
			name: "successful deletion",
			id:   1,
			mockFunc: func(id int) error {
				if id != 1 {
					t.Errorf("expected id 1, got %d", id)
				}
				return nil
			},
		},
		{
			name:          "invalid ID - zero",
			id:            0,
			expectedError: domain.ErrInvalidID,
		},
		{
			name:          "invalid ID - negative",
			id:            -1,
			expectedError: domain.ErrInvalidID,
		},
		{
			name: "quote not found",
			id:   999,
			mockFunc: func(id int) error {
				return domain.ErrQuoteNotFound
			},
			expectedError: domain.ErrQuoteNotFound,
		},
		{
			name: "repository error",
			id:   1,
			mockFunc: func(id int) error {
				return errors.New("database error")
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockQuoteRepository{
				DeleteFunc: tt.mockFunc,
			}
			useCase := NewQuoteUseCase(mockRepo)

			err := useCase.DeleteQuote(tt.id)

			if tt.expectedError != nil {
				if err == nil || err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestQuoteUseCase_CreateQuote_CreatedAtSet tests that CreatedAt is set when creating a quote
func TestQuoteUseCase_CreateQuote_CreatedAtSet(t *testing.T) {
	beforeTest := time.Now()

	mockRepo := &MockQuoteRepository{
		CreateFunc: func(quote *domain.Quote) (*domain.Quote, error) {
			// Verify CreatedAt was set
			if quote.CreatedAt.IsZero() {
				t.Error("CreatedAt should be set")
			}
			if quote.CreatedAt.Before(beforeTest) {
				t.Error("CreatedAt should be after test start time")
			}
			quote.ID = 1
			return quote, nil
		},
	}

	useCase := NewQuoteUseCase(mockRepo)

	request := &domain.CreateQuoteRequest{
		Author: "Test Author",
		Quote:  "Test Quote",
	}

	result, err := useCase.CreateQuote(request)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.CreatedAt.IsZero() {
		t.Error("result CreatedAt should be set")
	}
}
