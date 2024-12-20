package utils

import (
	"github.com/google/uuid"
	"log"
	"os"
	"strings"
	"sync"
)

type Logger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
	machineLogger *log.Logger
	logLevel      string
	mu            sync.Mutex
	requestUUID   string
}

var Log = Logger{
	infoLogger:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
	warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
	errorLogger:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	debugLogger:   log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	machineLogger: log.New(os.Stdout, "MACHINE: ", log.Ldate|log.Ltime|log.Lshortfile),
	logLevel:      strings.ToUpper(os.Getenv("LOG_LEVEL")),
}

func (l *Logger) SetRequestUUID() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.requestUUID = uuid.NewString()
}

func (l *Logger) ClearRequestUUID() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.requestUUID = ""
}

func (l *Logger) logWithPrefix(logger *log.Logger, format string, v ...interface{}) {
	l.mu.Lock()

	requestUUID := l.requestUUID
	l.mu.Unlock()
	logger.Printf("["+requestUUID+"] "+format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.SetRequestUUID()
	defer l.ClearRequestUUID()
	if l.logLevel == "MACHINE" || l.logLevel == "DEBUG" || l.logLevel == "INFO" {
		l.logWithPrefix(l.infoLogger, format, v...)
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.SetRequestUUID()
	defer l.ClearRequestUUID()
	if l.logLevel == "MACHINE" || l.logLevel == "DEBUG" {
		l.logWithPrefix(l.debugLogger, format, v...)
	}
}

func (l *Logger) Warning(format string, v ...interface{}) {
	l.SetRequestUUID()
	defer l.ClearRequestUUID()
	if l.logLevel == "MACHINE" || l.logLevel == "DEBUG" || l.logLevel == "INFO" || l.logLevel == "WARNING" {
		l.logWithPrefix(l.warningLogger, format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.SetRequestUUID()
	defer l.ClearRequestUUID()
	if l.logLevel == "MACHINE" || l.logLevel == "DEBUG" || l.logLevel == "INFO" || l.logLevel == "WARNING" || l.logLevel == "ERROR" {
		l.logWithPrefix(l.errorLogger, format, v...)
	}
}
