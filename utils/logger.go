package utils

import (
	"context"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type contextKey string

type Logger struct {
	infoLogger     *log.Logger
	warningLogger  *log.Logger
	errorLogger    *log.Logger
	debugLogger    *log.Logger
	machineLogger  *log.Logger
	logLevel       string
	mu             sync.Mutex
	requestUUIDKey contextKey
}

var Log = Logger{
	infoLogger:     log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
	warningLogger:  log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
	errorLogger:    log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	debugLogger:    log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	machineLogger:  log.New(os.Stdout, "MACHINE: ", log.Ldate|log.Ltime|log.Lshortfile),
	logLevel:       strings.ToUpper(os.Getenv("LOG_LEVEL")),
	requestUUIDKey: contextKey("requestUUID"),
}

func (l *Logger) RequestUUIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestUUID := uuid.NewString()
		ctx := context.WithValue(r.Context(), l.requestUUIDKey, requestUUID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (l *Logger) getRequestUUID(ctx context.Context) string {
	if uuid, ok := ctx.Value(l.requestUUIDKey).(string); ok {
		return uuid
	}
	return "no request uuid"
}

func (l *Logger) logWithPrefix(ctx context.Context, logger *log.Logger, format string, v ...interface{}) {
	requestUUID := l.getRequestUUID(ctx)
	logger.Printf("["+requestUUID+"] "+format, v...)
}

func (l *Logger) Info(ctx context.Context, format string, v ...interface{}) {
	if l.logLevel == "MACHINE" || l.logLevel == "DEBUG" || l.logLevel == "INFO" {
		l.logWithPrefix(ctx, l.infoLogger, format, v...)
	}
}

func (l *Logger) Debug(ctx context.Context, format string, v ...interface{}) {
	if l.logLevel == "MACHINE" || l.logLevel == "DEBUG" {
		l.logWithPrefix(ctx, l.debugLogger, format, v...)
	}
}

func (l *Logger) Warning(ctx context.Context, format string, v ...interface{}) {
	if l.logLevel == "MACHINE" || l.logLevel == "DEBUG" || l.logLevel == "INFO" || l.logLevel == "WARNING" {
		l.logWithPrefix(ctx, l.warningLogger, format, v...)
	}
}

func (l *Logger) Error(ctx context.Context, format string, v ...interface{}) {
	if l.logLevel == "MACHINE" || l.logLevel == "DEBUG" || l.logLevel == "INFO" || l.logLevel == "WARNING" || l.logLevel == "ERROR" {
		l.logWithPrefix(ctx, l.errorLogger, format, v...)
	}
}
