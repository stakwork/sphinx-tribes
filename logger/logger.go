package logger

import (
	"github.com/google/uuid"
	"github.com/stakwork/sphinx-tribes/config"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
)

type Logger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
	machineLogger *log.Logger
	mu            sync.Mutex
	requestUUID   string
}

var Log = Logger{
	infoLogger:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
	warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
	errorLogger:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime),
	debugLogger:   log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime),
	machineLogger: log.New(os.Stdout, "MACHINE: ", log.Ldate|log.Ltime),
}

func (l *Logger) SetRequestUUID(uuidString string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.requestUUID = uuidString
}

func (l *Logger) ClearRequestUUID() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.requestUUID = ""
}

func RouteBasedUUIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := uuid.NewString()
		Log.SetRequestUUID(uuid)

		defer Log.ClearRequestUUID()

		next.ServeHTTP(w, r)
	})
}

func (l *Logger) logWithPrefix(logger *log.Logger, format string, v ...interface{}) {
	l.mu.Lock()

	requestUUID := l.requestUUID
	l.mu.Unlock()

	var file string
	var line int
	var ok bool

	// Use runtime.Caller with skip 3 to go to the caller of the method that called logWithPrefix
	_, file, line, ok = runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	shortFile := filepath.Base(file)
	line_str := strconv.Itoa(line)

	if requestUUID == "" {
		logger.Printf("["+shortFile+":"+line_str+"] "+format, v...)
	} else {
		logger.Printf("["+shortFile+":"+line_str+"] ["+requestUUID+"] "+format, v...)
	}
}

func (l *Logger) Machine(format string, v ...interface{}) {
	if config.LogLevel == "MACHINE" {
		l.logWithPrefix(l.machineLogger, format, v...)
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if config.LogLevel == "MACHINE" || config.LogLevel == "DEBUG" {
		l.logWithPrefix(l.debugLogger, format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if config.LogLevel == "MACHINE" || config.LogLevel == "DEBUG" || config.LogLevel == "INFO" {
		l.logWithPrefix(l.infoLogger, format, v...)
	}
}

func (l *Logger) Warning(format string, v ...interface{}) {
	if config.LogLevel == "MACHINE" || config.LogLevel == "DEBUG" || config.LogLevel == "INFO" || config.LogLevel == "WARNING" {
		l.logWithPrefix(l.warningLogger, format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if config.LogLevel == "MACHINE" || config.LogLevel == "DEBUG" || config.LogLevel == "INFO" || config.LogLevel == "WARNING" || config.LogLevel == "ERROR" {
		l.logWithPrefix(l.errorLogger, format, v...)
	}
}
