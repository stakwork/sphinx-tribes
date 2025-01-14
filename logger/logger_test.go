package logger

import (
	"bytes"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stretchr/testify/assert"
)

func normalizeLogLevel(level string) string {
	return strings.ToUpper(strings.TrimSpace(level))
}

func TestLogger_Warning(t *testing.T) {

	originalStdout := os.Stdout
	_, w, _ := os.Pipe()
	os.Stdout = w

	originalLogLevel := config.LogLevel
	defer func() {
		config.LogLevel = originalLogLevel
		os.Stdout = originalStdout
	}()

	tests := []struct {
		name           string
		logLevel       string
		message        string
		args           []interface{}
		setupLogger    func() *Logger
		expectedOutput string
		concurrent     bool
	}{
		{
			name:     "Basic Warning with WARNING Level",
			logLevel: "WARNING",
			message:  "Test warning message",
			args:     []interface{}{},
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "WARNING: ",
		},
		{
			name:     "Warning with Format Arguments",
			logLevel: "WARNING",
			message:  "Test warning %s %d",
			args:     []interface{}{"message", 123},
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "Test warning message 123",
		},
		{
			name:     "Warning with RequestUUID",
			logLevel: "WARNING",
			message:  "Test with UUID",
			setupLogger: func() *Logger {
				l := &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
				l.SetRequestUUID("test-uuid")
				return l
			},
			expectedOutput: "[test-uuid]",
		},
		{
			name:     "Warning with ERROR Level (Should Not Log)",
			logLevel: "ERROR",
			message:  "Should not appear",
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "",
		},
		{
			name:     "Warning with Newline and Tab Characters",
			logLevel: "WARNING",
			message:  "Line1\nLine2\tTabbed",
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "Line1\nLine2\tTabbed",
		},
		{
			name:       "Concurrent Warning Calls",
			logLevel:   "WARNING",
			message:    "Concurrent message %d",
			concurrent: true,
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "Concurrent message",
		},
		{
			name:     "Empty Message",
			logLevel: "WARNING",
			message:  "",
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "WARNING: ",
		},
		{
			name:     "Very Long Message",
			logLevel: "WARNING",
			message:  strings.Repeat("a", 1000),
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: strings.Repeat("a", 1000),
		},
		{
			name:     "Unicode Characters",
			logLevel: "WARNING",
			message:  "测试警告 テスト警告",
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "测试警告 テスト警告",
		},
		{
			name:     "Multiple UUID Changes",
			logLevel: "WARNING",
			message:  "Test message",
			setupLogger: func() *Logger {
				l := &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
				l.SetRequestUUID("uuid-1")
				l.SetRequestUUID("uuid-2")
				return l
			},
			expectedOutput: "[uuid-2]",
		},
		{
			name:     "Log Level is MACHINE",
			logLevel: "MACHINE",
			message:  "Machine level message",
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "Machine level message",
		},
		{
			name:     "Log Level is DEBUG",
			logLevel: "DEBUG",
			message:  "Debug level message",
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "Debug level message",
		},
		{
			name:     "Log Level is INFO",
			logLevel: "INFO",
			message:  "Info level message",
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "Info level message",
		},
		{
			name:     "Empty Format String",
			logLevel: "WARNING",
			message:  "",
			args:     []interface{}{"ignored", "args"},
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "WARNING: ",
		},
		{
			name:     "No Variadic Arguments",
			logLevel: "WARNING",
			message:  "Simple message without args",
			args:     []interface{}{},
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "Simple message without args",
		},
		{
			name:     "Large Number of Variadic Arguments",
			logLevel: "WARNING",
			message:  strings.Repeat("%v ", 100),
			setupLogger: func() *Logger {
				args := make([]interface{}, 100)
				for i := range args {
					args[i] = i
				}
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "WARNING: ",
		},
		{
			name:     "Log Level is Case Insensitive",
			logLevel: "WaRnInG",
			message:  "Case insensitive test",
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "Case insensitive test",
		},
		{
			name:     "Log Level is Empty String",
			logLevel: "",
			message:  "Empty log level test",
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "",
		},
		{
			name:     "Format String with Multiple Lines",
			logLevel: "WARNING",
			message:  "Line1\nLine2\nLine3",
			setupLogger: func() *Logger {
				return &Logger{
					warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
				}
			},
			expectedOutput: "Line1\nLine2\nLine3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config.LogLevel = normalizeLogLevel(tt.logLevel)
			logger := tt.setupLogger()

			buf := &bytes.Buffer{}
			logger.warningLogger.SetOutput(buf)

			if tt.concurrent {
				var wg sync.WaitGroup
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func(i int) {
						defer wg.Done()
						logger.Warning(tt.message, i)
					}(i)
				}
				wg.Wait()
			} else {
				logger.Warning(tt.message, tt.args...)
			}

			time.Sleep(10 * time.Millisecond)
			output := buf.String()

			if tt.logLevel == "ERROR" {
				assert.Empty(t, output)
			} else if normalizeLogLevel(tt.logLevel) == "WARNING" ||
				normalizeLogLevel(tt.logLevel) == "INFO" ||
				normalizeLogLevel(tt.logLevel) == "DEBUG" ||
				normalizeLogLevel(tt.logLevel) == "MACHINE" {
				if tt.expectedOutput != "" {
					containsExpected := strings.Contains(
						strings.TrimSpace(output),
						tt.expectedOutput,
					)
					assert.True(t, containsExpected,
						"Expected output '%s' not found in actual output '%s'",
						tt.expectedOutput, output)
				}
				assert.Contains(t, output, "WARNING: ")
				assert.Contains(t, output, "[logger_test.go:")
			}

			buf.Reset()
		})
	}
}

func TestLoggerConcurrency(t *testing.T) {
	logger := &Logger{
		warningLogger: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime),
		mu:            sync.Mutex{},
	}

	var wg sync.WaitGroup
	iterations := 100

	for i := 0; i < iterations; i++ {
		wg.Add(3)

		go func() {
			defer wg.Done()
			logger.SetRequestUUID("uuid-1")
		}()

		go func() {
			defer wg.Done()
			logger.Warning("Test message")
		}()

		go func() {
			defer wg.Done()
			logger.ClearRequestUUID()
		}()
	}

	wg.Wait()
}
