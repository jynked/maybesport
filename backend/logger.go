package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func initLogger() {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func logError(msg string, err error, attrs ...any) {
	args := append([]any{"error", err}, attrs...)
	slog.Error(msg, args...)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lw, r)
		duration := time.Since(start)

		var userID int
		if uid, ok := r.Context().Value("userID").(int); ok {
			userID = uid
		}

		if lw.statusCode >= 400 {
			slog.Warn("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", lw.statusCode),
				slog.Duration("duration", duration),
				slog.String("ip", r.RemoteAddr),
				slog.Int("user_id", userID),
			)
		} else {
			slog.Info("request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", lw.statusCode),
				slog.Duration("duration", duration),
			)
		}
	})
}

func authContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, err := getUserIDFromToken(r)
		if err == nil && userID != 0 {
			ctx := context.WithValue(r.Context(), "userID", userID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
