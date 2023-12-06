package middleware

import (
	"net/http"
	"time"

	"loyalty-system/internal/auth"
	"loyalty-system/internal/logger"
)

type loggerResponseWriter struct {
	rw     http.ResponseWriter
	status int
	size   int
}

func WithRequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		respWriter := loggerResponseWriter{rw: w}

		logger.WithFields(logger.Fields{
			"size":     respWriter.size,
			"status":   respWriter.status,
			"duration": time.Since(start),
			"method":   r.Method,
			"uri":      r.RequestURI,
			"user":     r.Header.Get(auth.AuthHeaderName),
		}).Info("incoming request")

		h.ServeHTTP(&respWriter, r)
	})
}

func (l *loggerResponseWriter) Write(b []byte) (int, error) {
	size, error := l.rw.Write(b)
	l.size = size

	return size, error
}

func (l *loggerResponseWriter) WriteHeader(statusCode int) {
	l.status = statusCode
	l.rw.WriteHeader(statusCode)
}

func (l *loggerResponseWriter) Header() http.Header {
	return l.rw.Header()
}
