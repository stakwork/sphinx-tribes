package utils

import (
	"log"
	"os"
	"strings"
)

type Logger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
	machineLogger *log.Logger
	logLevel      string
}

var Log = Logger{
	infoLogger:    log.New(os.Stdout, "INF: ", log.Ldate|log.Ltime|log.Lshortfile),
	warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
	errorLogger:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	debugLogger:   log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	machineLogger: log.New(os.Stdout, "MACHINE: ", log.Ldate|log.Ltime|log.Lshortfile),
	logLevel:      strings.ToUpper(os.Getenv("LOG_LEVEL")),
}

func (l *Logger) Machine(format string, v ...interface{}) {
	if l.logLevel == "MACHINE" {
		l.machineLogger.Printf(format, v...)
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if l.logLevel == "MACHINE" || l.logLevel == "DEBUG" {
		l.debugLogger.Printf(format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.logLevel == "MACHINE" || l.logLevel == "DEBUG" || l.logLevel == "INFO" {
		l.infoLogger.Printf(format, v...)
	}
}

func (l *Logger) Warning(format string, v ...interface{}) {
	if l.logLevel == "MACHINE" || l.logLevel == "DEBUG" || l.logLevel == "INFO" || l.logLevel == "WARNING" {
		l.warningLogger.Printf(format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.logLevel == "MACHINE" || l.logLevel == "DEBUG" || l.logLevel == "INFO" || l.logLevel == "WARNING" || l.logLevel == "ERROR" {
		l.errorLogger.Printf(format, v...)
	}
}
