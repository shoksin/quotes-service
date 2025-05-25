package handler

type QuoteService interface {
	CreateQuote()
}

type QuoteHandler struct {
	quoteService QuoteService
}

func NewQuoteHandler(quoteService QuoteService) *QuoteHandler {
	return &QuoteHandler{quoteService: quoteService}
}

func (h *QuoteHandler) CreateQuote() {}
