package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestLoggingResponseWriter(t *testing.T) {
	w := httptest.NewRecorder()
	lrw := newLoggingResponseWriter(w)

	header := lrw.Header()
	if header == nil {
		t.Error("Header() returned nil")
	}

	lrw.WriteHeader(http.StatusCreated)
	if lrw.statusCode != http.StatusCreated {
		t.Errorf("expected status code %d, got %d", http.StatusCreated, lrw.statusCode)
	}

	testData := []byte("test data")
	n, err := lrw.Write(testData)
	if err != nil {
		t.Errorf("Write() returned error: %v", err)
	}
	if n != len(testData) {
		t.Errorf("expected %d bytes written, got %d", len(testData), n)
	}
	if lrw.bytes != len(testData) {
		t.Errorf("expected %d bytes counted, got %d", len(testData), lrw.bytes)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	var logBuffer bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&logBuffer)
	defer log.SetOutput(originalOutput)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	wrappedHandler := LoggingMiddleware(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "GET") {
		t.Error("log output should contain HTTP method")
	}
	if !strings.Contains(logOutput, "/test") {
		t.Error("log output should contain request path")
	}
	if !strings.Contains(logOutput, "200") {
		t.Error("log output should contain status code")
	}
	if !strings.Contains(logOutput, "13B") { // "test response" is 13 bytes
		t.Error("log output should contain response size")
	}
}

func TestLoggingMiddleware_DifferentStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
	}{
		{
			name:       "OK status",
			statusCode: http.StatusOK,
			body:       "OK",
		},
		{
			name:       "Not Found status",
			statusCode: http.StatusNotFound,
			body:       "Not Found",
		},
		{
			name:       "Internal Server Error",
			statusCode: http.StatusInternalServerError,
			body:       "Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var logBuffer bytes.Buffer
			originalOutput := log.Writer()
			log.SetOutput(&logBuffer)
			defer log.SetOutput(originalOutput)

			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			})

			wrappedHandler := LoggingMiddleware(testHandler)

			req := httptest.NewRequest(http.MethodPost, "/test", nil)
			rec := httptest.NewRecorder()
			wrappedHandler.ServeHTTP(rec, req)

			if rec.Code != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, rec.Code)
			}

			logOutput := logBuffer.String()
			statusStr := strconv.Itoa(tt.statusCode)
			if !strings.Contains(logOutput, statusStr) {
				t.Errorf("log output should contain status code %d", tt.statusCode)
			}
		})
	}
}

func TestLoggingMiddleware_DefaultStatusCode(t *testing.T) {
	var logBuffer bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&logBuffer)
	defer log.SetOutput(originalOutput)

	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't call WriteHeader, should default to 200
		w.Write([]byte("test"))
	})

	wrappedHandler := LoggingMiddleware(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rec, req)

	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, "200") {
		t.Error("log output should contain default status code 200")
	}
}
