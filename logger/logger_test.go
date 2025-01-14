package logger

import (
	"bytes"
	"fmt"
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

func TestMachine(t *testing.T) {

	originalLogLevel := config.LogLevel
	defer func() {
		config.LogLevel = originalLogLevel
	}()

	tests := []struct {
		name         string
		logLevel     string
		format       string
		args         []interface{}
		expectedLogs bool
		setup        func(*Logger)
	}{
		{
			name:         "Machine level enabled with simple message",
			logLevel:     "MACHINE",
			format:       "Test message",
			args:         []interface{}{},
			expectedLogs: true,
		},
		{
			name:         "Machine level disabled",
			logLevel:     "INFO",
			format:       "Test message",
			args:         []interface{}{},
			expectedLogs: false,
		},
		{
			name:         "Machine level with format parameters",
			logLevel:     "MACHINE",
			format:       "Test message %s %d",
			args:         []interface{}{"param", 123},
			expectedLogs: true,
		},
		{
			name:         "Machine level with special characters",
			logLevel:     "MACHINE",
			format:       "Test !@#$%^&*()",
			args:         []interface{}{},
			expectedLogs: true,
		},
		{
			name:         "Machine level with empty message",
			logLevel:     "MACHINE",
			format:       "",
			args:         []interface{}{},
			expectedLogs: true,
		},
		{
			name:         "Machine level with nil arguments",
			logLevel:     "MACHINE",
			format:       "Test with nil: %v",
			args:         []interface{}{nil},
			expectedLogs: true,
		},
		{
			name:         "Machine level with multiple arguments",
			logLevel:     "MACHINE",
			format:       "%v %v %v %v",
			args:         []interface{}{1, "two", true, 4.5},
			expectedLogs: true,
		},
		{
			name:         "Machine level with Unicode characters",
			logLevel:     "MACHINE",
			format:       "Unicode test: %s",
			args:         []interface{}{"你好 👋 Привет"},
			expectedLogs: true,
		},
		{
			name:         "Machine level with large message",
			logLevel:     "MACHINE",
			format:       "Large message: %s",
			args:         []interface{}{strings.Repeat("a", 1000)},
			expectedLogs: true,
		},
		{
			name:         "Machine level with invalid format specifier",
			logLevel:     "MACHINE",
			format:       "Invalid format %z",
			args:         []interface{}{"test"},
			expectedLogs: true,
		},
		{
			name:         "Standard Logging with MACHINE LogLevel",
			logLevel:     "MACHINE",
			format:       "Standard log message: %s",
			args:         []interface{}{"test"},
			expectedLogs: true,
		},
		{
			name:         "LogLevel Case Sensitivity",
			logLevel:     "machine",
			format:       "Case sensitivity test",
			args:         []interface{}{},
			expectedLogs: false,
		},
		{
			name:         "Large Number of Arguments",
			logLevel:     "MACHINE",
			format:       "%v %v %v %v %v %v %v %v %v %v",
			args:         []interface{}{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expectedLogs: true,
		},
		{
			name:         "Long Format String",
			logLevel:     "MACHINE",
			format:       strings.Repeat("Very long message with placeholder %s ", 100),
			args:         []interface{}{strings.Repeat("test", 100)},
			expectedLogs: true,
		},
		{
			name:         "Format String with Special Characters",
			logLevel:     "MACHINE",
			format:       "Special chars: %s\n\t\r\b\f%s",
			args:         []interface{}{"test1", "test2"},
			expectedLogs: true,
		},
		{
			name:         "Format String with Incorrect Placeholders",
			logLevel:     "MACHINE",
			format:       "Incorrect placeholders: %d %s",
			args:         []interface{}{"string", 123},
			expectedLogs: true,
		},
		{
			name:         "Invalid LogLevel",
			logLevel:     "INVALID_LEVEL",
			format:       "Test message",
			args:         []interface{}{},
			expectedLogs: false,
		},
		{
			name:         "Mixed Type Arguments",
			logLevel:     "MACHINE",
			format:       "%v %v %v %v %v",
			args:         []interface{}{123, "string", true, 45.67, struct{ Name string }{"test"}},
			expectedLogs: true,
		},
		{
			name:         "Format String with Unicode Placeholders",
			logLevel:     "MACHINE",
			format:       "Unicode: %s 你好 %s Привет %s",
			args:         []interface{}{"Hello", "World", "!"},
			expectedLogs: true,
		},
		{
			name:         "Zero Values Arguments",
			logLevel:     "MACHINE",
			format:       "Zero values: %v %v %v %v",
			args:         []interface{}{0, "", false, nil},
			expectedLogs: true,
		},
		{
			name:         "Escaped Percent Signs",
			logLevel:     "MACHINE",
			format:       "Escaped %%s %%d test %s",
			args:         []interface{}{"actual"},
			expectedLogs: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var buf bytes.Buffer

			logger := &Logger{
				machineLogger: log.New(&buf, "", 0),
			}

			config.LogLevel = tt.logLevel

			if tt.setup != nil {
				tt.setup(logger)
			}

			logger.Machine(tt.format, tt.args...)

			output := buf.String()

			if tt.expectedLogs {
				assert.NotEmpty(t, output, "Expected log output but got none")

				expectedOutput := fmt.Sprintf(tt.format, tt.args...)
				assert.Contains(t, output, expectedOutput,
					"Log output doesn't contain expected message")

			} else {
				assert.Empty(t, output, "Expected no log output but got some")
			}
		})
	}
}
