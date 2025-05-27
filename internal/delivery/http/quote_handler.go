package handler

import (
	"encoding/json"
	"github.com/shoksin/quotes-service/internal/domain"
	"net/http"
	"strconv"
	"strings"
)

type QuoteUseCase interface {
	CreateQuote(req *domain.CreateQuoteRequest) (*domain.Quote, error)
	GetAllQuotes() ([]*domain.Quote, error)
	GetQuotesByAuthor(author string) ([]*domain.Quote, error)
	GetRandomQuote() (*domain.Quote, error)
	DeleteQuote(id int) error
}

type QuoteHandler struct {
	quoteUseCase QuoteUseCase
}

func NewQuoteHandler(quoteUseCase QuoteUseCase) *QuoteHandler {
	return &QuoteHandler{
		quoteUseCase: quoteUseCase,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *QuoteHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *QuoteHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, ErrorResponse{Error: message})
}

// CreateQuote POST /quotes
func (h *QuoteHandler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateQuoteRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, domain.MsgInvalidJSON)
		return
	}

	quote, err := h.quoteUseCase.CreateQuote(&req)
	if err != nil {
		switch err {
		case domain.ErrInvalidAuthor, domain.ErrInvalidQuote:
			h.writeError(w, http.StatusBadRequest, err.Error())
		default:
			h.writeError(w, http.StatusInternalServerError, domain.MsgFailedCreateQuote)
		}
		return
	}

	h.writeJSON(w, http.StatusCreated, quote)
}

// GetQuotes GET /quotes с опциональным фильтром ?author=Name
func (h *QuoteHandler) GetQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	var quotes []*domain.Quote
	var err error

	if author != "" {
		quotes, err = h.quoteUseCase.GetQuotesByAuthor(author)
	} else {
		quotes, err = h.quoteUseCase.GetAllQuotes()
	}

	if err != nil {
		switch err {
		case domain.ErrInvalidAuthor:
			h.writeError(w, http.StatusBadRequest, err.Error())
		default:
			h.writeError(w, http.StatusInternalServerError, domain.MsgFailedGetQuotes)
		}
		return
	}
	if len(quotes) == 0 {
		h.writeJSON(w, http.StatusOK, struct{}{})
	} else {
		h.writeJSON(w, http.StatusOK, quotes)
	}
}

// GetRandomQuote GET /quotes/random
func (h *QuoteHandler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := h.quoteUseCase.GetRandomQuote()
	if err != nil {
		switch err {
		case domain.ErrNoQuotesFound:
			h.writeError(w, http.StatusNotFound, domain.MsgQuotesNotFound)
		default:
			h.writeError(w, http.StatusInternalServerError, domain.MsgFailedGetRandomQuote)
		}
		return
	}

	h.writeJSON(w, http.StatusOK, quote)
}

// DeleteQuote DELETE /quotes/{id}
func (h *QuoteHandler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/quotes/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, domain.MsgInvalidQuoteID)
		return
	}

	err = h.quoteUseCase.DeleteQuote(id)
	if err != nil {
		switch err {
		case domain.ErrInvalidID:
			h.writeError(w, http.StatusBadRequest, err.Error())
		case domain.ErrQuoteNotFound:
			h.writeError(w, http.StatusNotFound, domain.MsgQuotesNotFound)
		default:
			h.writeError(w, http.StatusInternalServerError, domain.MsgFailedDeleteQuote)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// HealthCheck handles health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status":  "healthy",
		"service": "quotes-service",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RegisterRoutes register all handler for QuoteHandler
func (h *QuoteHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", HealthCheck)

	mux.HandleFunc("/quotes/random", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetRandomQuote(w, r)
		} else {
			h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/quotes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateQuote(w, r)
		case http.MethodGet:
			h.GetQuotes(w, r)
		default:
			h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	mux.HandleFunc("/quotes/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			h.DeleteQuote(w, r)
		} else {
			h.writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})
}
