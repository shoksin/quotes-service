package middleware

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	w          http.ResponseWriter
	statusCode int
	bytes      int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w: w, statusCode: http.StatusOK}
}

func (lrw *loggingResponseWriter) Header() http.Header {
	return lrw.w.Header()
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.w.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.w.Write(b)
	lrw.bytes += n
	return n, err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)

		log.Printf("%s %s %s %d %dB %s", r.RemoteAddr, r.Method, r.URL.Path, lrw.statusCode, lrw.bytes, duration)
	})
}
