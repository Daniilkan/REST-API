package logger

import (
	"bytes"
	"context"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"net/http"
)

const (
	LoggerKey = "logger"
)

type Logger struct {
	l *zap.Logger
}

func New(ctx context.Context) (context.Context, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, LoggerKey, &Logger{logger})

	return ctx, nil
}

func GetLoggerFromContext(ctx context.Context) *Logger {
	return ctx.Value(LoggerKey).(*Logger)
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	if middleware.GetReqID(ctx) != "" {
		fields = append(fields, zap.String("RequestID", middleware.GetReqID(ctx)))
	}
	l.l.Info(msg, fields...)
}
func (l *Logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	if middleware.GetReqID(ctx) != "" {
		fields = append(fields, zap.String("RequestID", middleware.GetReqID(ctx)))
	}
	l.l.Fatal(msg, fields...)
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	rr.body.Write(b)
	return rr.ResponseWriter.Write(b)
}

func Middleware(baseCtx context.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := GetLoggerFromContext(baseCtx)
			if logger != nil {
				logger.Info(r.Context(), "Incoming request", zap.String("method", r.Method), zap.String("url", r.URL.String()), zap.Any("params", r.URL.Query()))
			}

			recorder := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK, body: &bytes.Buffer{}}

			ctx := context.WithValue(r.Context(), LoggerKey, logger)
			next.ServeHTTP(recorder, r.WithContext(ctx))

			if logger != nil {
				logger.Info(ctx, "Response sent", zap.Int("status", recorder.statusCode), zap.String("response", recorder.body.String()))
			}
		})
	}
}
