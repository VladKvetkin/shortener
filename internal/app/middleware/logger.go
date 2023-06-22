package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func Logger(next http.Handler) http.Handler {
	logFunc := func(resp http.ResponseWriter, req *http.Request) {
		start := time.Now()

		responseData := &responseData{}

		loggingResponse := loggingResponseWriter{
			ResponseWriter: resp,
			responseData:   responseData,
		}
		next.ServeHTTP(&loggingResponse, req)

		duration := time.Since(start)

		zap.L().Info(
			"HTTP request",
			zap.String("uri", req.RequestURI),
			zap.String("method", req.Method),
			zap.Duration("duration", duration),
		)

		zap.L().Info(
			"HTTP response",
			zap.Int("status", responseData.status),
			zap.Int("size", responseData.size),
		)
	}

	return http.HandlerFunc(logFunc)
}
