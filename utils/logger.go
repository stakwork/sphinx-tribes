package utils

import (
	"github.com/stakwork/sphinx-tribes/config"
	"log"
	"os"
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
	infoLogger:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
	warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
	errorLogger:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	debugLogger:   log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	machineLogger: log.New(os.Stdout, "MACHINE: ", log.Ldate|log.Ltime|log.Lshortfile),
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

func (l *Logger) logWithPrefix(logger *log.Logger, format string, v ...interface{}) {
	l.mu.Lock()

	requestUUID := l.requestUUID
	l.mu.Unlock()

	if requestUUID == "" {
		logger.Printf(format, v...)
	} else {
		logger.Printf("["+requestUUID+"] "+format, v...)
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
