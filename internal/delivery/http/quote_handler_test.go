package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shoksin/quotes-service/internal/domain"
)

type MockQuoteUseCase struct {
	CreateQuoteFunc       func(req *domain.CreateQuoteRequest) (*domain.Quote, error)
	GetAllQuotesFunc      func() ([]*domain.Quote, error)
	GetQuotesByAuthorFunc func(author string) ([]*domain.Quote, error)
	GetRandomQuoteFunc    func() (*domain.Quote, error)
	DeleteQuoteFunc       func(id int) error
}

func (m *MockQuoteUseCase) CreateQuote(req *domain.CreateQuoteRequest) (*domain.Quote, error) {
	if m.CreateQuoteFunc != nil {
		return m.CreateQuoteFunc(req)
	}
	return nil, nil
}

func (m *MockQuoteUseCase) GetAllQuotes() ([]*domain.Quote, error) {
	if m.GetAllQuotesFunc != nil {
		return m.GetAllQuotesFunc()
	}
	return nil, nil
}

func (m *MockQuoteUseCase) GetQuotesByAuthor(author string) ([]*domain.Quote, error) {
	if m.GetQuotesByAuthorFunc != nil {
		return m.GetQuotesByAuthorFunc(author)
	}
	return nil, nil
}

func (m *MockQuoteUseCase) GetRandomQuote() (*domain.Quote, error) {
	if m.GetRandomQuoteFunc != nil {
		return m.GetRandomQuoteFunc()
	}
	return nil, nil
}

func (m *MockQuoteUseCase) DeleteQuote(id int) error {
	if m.DeleteQuoteFunc != nil {
		return m.DeleteQuoteFunc(id)
	}
	return nil
}

func TestQuoteHandler_CreateQuote(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(req *domain.CreateQuoteRequest) (*domain.Quote, error)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "successful quote creation",
			requestBody: map[string]string{
				"author": "Test Author",
				"quote":  "Test Quote",
			},
			mockFunc: func(req *domain.CreateQuoteRequest) (*domain.Quote, error) {
				return &domain.Quote{
					ID:        1,
					Author:    "Test Author",
					Quote:     "Test Quote",
					CreatedAt: time.Now(),
				}, nil
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   nil,
		},
		{
			name:           "invalid JSON",
			requestBody:    "invalid json",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": domain.MsgInvalidJSON},
		},
		{
			name: "invalid author error",
			requestBody: map[string]string{
				"author": "",
				"quote":  "Test Quote",
			},
			mockFunc: func(req *domain.CreateQuoteRequest) (*domain.Quote, error) {
				return nil, domain.ErrInvalidAuthor
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": domain.ErrInvalidAuthor.Error()},
		},
		{
			name: "internal server error",
			requestBody: map[string]string{
				"author": "Test Author",
				"quote":  "Test Quote",
			},
			mockFunc: func(req *domain.CreateQuoteRequest) (*domain.Quote, error) {
				return nil, errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": domain.MsgFailedCreateQuote},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockQuoteUseCase{
				CreateQuoteFunc: tt.mockFunc,
			}
			handler := NewQuoteHandler(mockUseCase)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler.CreateQuote(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if tt.expectedBody != nil {
				var response interface{}
				json.NewDecoder(rec.Body).Decode(&response)

				expectedJSON, _ := json.Marshal(tt.expectedBody)
				actualJSON, _ := json.Marshal(response)

				if string(expectedJSON) != string(actualJSON) {
					t.Errorf("expected body %s, got %s", expectedJSON, actualJSON)
				}
			} else if tt.name == "successful quote creation" {
				// For successful creation, just verify it's a valid quote response
				var response domain.Quote
				err := json.NewDecoder(rec.Body).Decode(&response)
				if err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if response.ID == 0 || response.Author == "" || response.Quote == "" {
					t.Error("response missing required fields")
				}
			}
		})
	}
}

func TestQuoteHandler_GetQuotes(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockAllFunc    func() ([]*domain.Quote, error)
		mockAuthorFunc func(author string) ([]*domain.Quote, error)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "get all quotes successfully",
			mockAllFunc: func() ([]*domain.Quote, error) {
				return []*domain.Quote{
					{ID: 1, Author: "Author1", Quote: "Quote1"},
					{ID: 2, Author: "Author2", Quote: "Quote2"},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "get quotes by author",
			queryParams: "?author=Author1",
			mockAuthorFunc: func(author string) ([]*domain.Quote, error) {
				return []*domain.Quote{
					{ID: 1, Author: "Author1", Quote: "Quote1"},
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "empty quotes list",
			mockAllFunc: func() ([]*domain.Quote, error) {
				return []*domain.Quote{}, nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "No quotes found",
		},
		{
			name:        "invalid ",
			queryParams: "?author=",
			mockAuthorFunc: func(author string) ([]*domain.Quote, error) {
				return []*domain.Quote{}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "internal server error",
			mockAllFunc: func() ([]*domain.Quote, error) {
				return nil, errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockQuoteUseCase{
				GetAllQuotesFunc:      tt.mockAllFunc,
				GetQuotesByAuthorFunc: tt.mockAuthorFunc,
			}
			handler := NewQuoteHandler(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, "/quotes"+tt.queryParams, nil)
			rec := httptest.NewRecorder()

			handler.GetQuotes(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestQuoteHandler_GetRandomQuote(t *testing.T) {
	tests := []struct {
		name           string
		mockFunc       func() (*domain.Quote, error)
		expectedStatus int
	}{
		{
			name: "successful random quote",
			mockFunc: func() (*domain.Quote, error) {
				return &domain.Quote{
					ID:     1,
					Author: "Random Author",
					Quote:  "Random Quote",
				}, nil
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "no quotes found",
			mockFunc: func() (*domain.Quote, error) {
				return nil, domain.ErrNoQuotesFound
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "internal server error",
			mockFunc: func() (*domain.Quote, error) {
				return nil, errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockQuoteUseCase{
				GetRandomQuoteFunc: tt.mockFunc,
			}
			handler := NewQuoteHandler(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
			rec := httptest.NewRecorder()

			handler.GetRandomQuote(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestQuoteHandler_DeleteQuote(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockFunc       func(id int) error
		expectedStatus int
	}{
		{
			name: "successful deletion",
			url:  "/quotes/1",
			mockFunc: func(id int) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "invalid quote ID",
			url:            "/quotes/invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "invalid ID error",
			url:  "/quotes/0",
			mockFunc: func(id int) error {
				return domain.ErrInvalidID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "quote not found",
			url:  "/quotes/999",
			mockFunc: func(id int) error {
				return domain.ErrQuoteNotFound
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "internal server error",
			url:  "/quotes/1",
			mockFunc: func(id int) error {
				return errors.New("database error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockQuoteUseCase{
				DeleteQuoteFunc: tt.mockFunc,
			}
			handler := NewQuoteHandler(mockUseCase)

			req := httptest.NewRequest(http.MethodDelete, tt.url, nil)
			rec := httptest.NewRecorder()

			handler.DeleteQuote(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	HealthCheck(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response map[string]string
	json.NewDecoder(rec.Body).Decode(&response)

	if response["status"] != "healthy" {
		t.Errorf("expected status 'healthy', got '%s'", response["status"])
	}

	if response["service"] != "quotes-service" {
		t.Errorf("expected service 'quotes-service', got '%s'", response["service"])
	}
}

func TestQuoteHandler_RegisterRoutes(t *testing.T) {
	mockUseCase := &MockQuoteUseCase{}
	handler := NewQuoteHandler(mockUseCase)
	mux := http.NewServeMux()

	handler.RegisterRoutes(mux)

	// Test that routes are registered by making requests
	tests := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/health"},
		{http.MethodPost, "/quotes"},
		{http.MethodGet, "/quotes"},
		{http.MethodGet, "/quotes/random"},
		{http.MethodDelete, "/quotes/1"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)

		if rec.Code == http.StatusNotFound {
			t.Errorf("route %s %s not registered", tt.method, tt.path)
		}
	}
}
