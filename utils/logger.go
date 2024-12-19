package utils

import (
	"log"
	"os"
	"github.com/stakwork/sphinx-tribes/config"
)

type Logger struct {
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
	debugLogger   *log.Logger
	machineLogger *log.Logger
}

var Log = Logger{
	infoLogger:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
	warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
	errorLogger:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	debugLogger:   log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	machineLogger: log.New(os.Stdout, "MACHINE: ", log.Ldate|log.Ltime|log.Lshortfile),
}

func (l *Logger) Machine(format string, v ...interface{}) {
	if config.LogLevel == "MACHINE" {
		l.machineLogger.Printf(format, v...)
	}
}

func (l *Logger) Debug(format string, v ...interface{}) {
	if config.LogLevel == "MACHINE" || config.LogLevel == "DEBUG" {
		l.debugLogger.Printf(format, v...)
	}
}

func (l *Logger) Info(format string, v ...interface{}) {
	if config.LogLevel == "MACHINE" || config.LogLevel == "DEBUG" || config.LogLevel == "INFO" {
		l.infoLogger.Printf(format, v...)
	}
}

func (l *Logger) Warning(format string, v ...interface{}) {
	if config.LogLevel == "MACHINE" || config.LogLevel == "DEBUG" || config.LogLevel == "INFO" || config.LogLevel == "WARNING" {
		l.warningLogger.Printf(format, v...)
	}
}

func (l *Logger) Error(format string, v ...interface{}) {
	if config.LogLevel == "MACHINE" || config.LogLevel == "DEBUG" || config.LogLevel == "INFO" || config.LogLevel == "WARNING" || config.LogLevel == "ERROR" {
		l.errorLogger.Printf(format, v...)
	}
}
