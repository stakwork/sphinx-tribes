package logger

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
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

func TestClearRequestUUID(t *testing.T) {
	tests := []struct {
		name           string
		setup          func(*Logger)
		concurrent     bool
		concurrentNum  int
		validateBefore func(*testing.T, *Logger)
		validateAfter  func(*testing.T, *Logger)
	}{
		{
			name: "Clear existing UUID",
			setup: func(l *Logger) {
				l.SetRequestUUID("test-uuid")
			},
			validateBefore: func(t *testing.T, l *Logger) {
				assert.Equal(t, "test-uuid", l.requestUUID)
			},
			validateAfter: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name: "Clear empty UUID",
			setup: func(l *Logger) {
				l.requestUUID = ""
			},
			validateBefore: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
			validateAfter: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name: "Clear with concurrent access",
			setup: func(l *Logger) {
				l.SetRequestUUID("concurrent-test-uuid")
			},
			concurrent:    true,
			concurrentNum: 100,
			validateBefore: func(t *testing.T, l *Logger) {
				assert.Equal(t, "concurrent-test-uuid", l.requestUUID)
			},
			validateAfter: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name: "Clear with long UUID",
			setup: func(l *Logger) {
				l.SetRequestUUID(strings.Repeat("a", 1000))
			},
			validateBefore: func(t *testing.T, l *Logger) {
				assert.Equal(t, strings.Repeat("a", 1000), l.requestUUID)
			},
			validateAfter: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name: "Clear with special characters",
			setup: func(l *Logger) {
				l.SetRequestUUID("!@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`")
			},
			validateBefore: func(t *testing.T, l *Logger) {
				assert.Equal(t, "!@#$%^&*()_+{}[]|\\:;\"'<>,.?/~`", l.requestUUID)
			},
			validateAfter: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name: "Multiple clear calls",
			setup: func(l *Logger) {
				l.SetRequestUUID("multiple-clear-test")
			},
			validateBefore: func(t *testing.T, l *Logger) {
				assert.Equal(t, "multiple-clear-test", l.requestUUID)
			},
			validateAfter: func(t *testing.T, l *Logger) {
				l.ClearRequestUUID()
				l.ClearRequestUUID()
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name: "Clear with Unicode characters",
			setup: func(l *Logger) {
				l.SetRequestUUID("你好👋Привет")
			},
			validateBefore: func(t *testing.T, l *Logger) {
				assert.Equal(t, "你好👋Привет", l.requestUUID)
			},
			validateAfter: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name: "Basic Functionality",
			setup: func(l *Logger) {
				l.SetRequestUUID("basic-test-uuid")
			},
			validateBefore: func(t *testing.T, l *Logger) {
				assert.Equal(t, "basic-test-uuid", l.requestUUID)
			},
			validateAfter: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name: "Edge Case - Already Empty UUID",
			setup: func(l *Logger) {
				l.ClearRequestUUID()
			},
			validateBefore: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
			validateAfter: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name: "Concurrency Test",
			setup: func(l *Logger) {
				l.SetRequestUUID("concurrent-uuid")
			},
			concurrent:    true,
			concurrentNum: 1000,
			validateBefore: func(t *testing.T, l *Logger) {
				assert.Equal(t, "concurrent-uuid", l.requestUUID)
			},
			validateAfter: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name: "Performance Under Load",
			setup: func(l *Logger) {
				l.SetRequestUUID("performance-test-uuid")
			},
			validateBefore: func(t *testing.T, l *Logger) {
				start := time.Now()
				for i := 0; i < 10000; i++ {
					l.ClearRequestUUID()
					l.SetRequestUUID("performance-test-uuid")
				}
				duration := time.Since(start)
				assert.Less(t, duration, 1*time.Second, "Performance test took too long")
			},
			validateAfter: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name: "Logger Initialization",
			setup: func(l *Logger) {
				newLogger := &Logger{
					mu: sync.Mutex{},
				}
				newLogger.ClearRequestUUID()
				assert.Empty(t, newLogger.requestUUID)
			},
			validateAfter: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name: "Large UUID String",
			setup: func(l *Logger) {
				largeUUID := strings.Repeat("abcdef0123456789", 1000)
				l.SetRequestUUID(largeUUID)
			},
			validateBefore: func(t *testing.T, l *Logger) {
				assert.Equal(t, strings.Repeat("abcdef0123456789", 1000), l.requestUUID)
			},
			validateAfter: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name: "Repeated Calls",
			setup: func(l *Logger) {
				l.SetRequestUUID("repeated-test")
			},
			validateBefore: func(t *testing.T, l *Logger) {
				for i := 0; i < 100; i++ {
					l.ClearRequestUUID()
					assert.Empty(t, l.requestUUID)
					l.SetRequestUUID("repeated-test")
					assert.Equal(t, "repeated-test", l.requestUUID)
				}
			},
			validateAfter: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			logger := &Logger{
				mu: sync.Mutex{},
			}

			if tt.setup != nil {
				tt.setup(logger)
			}

			if tt.validateBefore != nil {
				tt.validateBefore(t, logger)
			}

			if tt.concurrent {

				var wg sync.WaitGroup
				for i := 0; i < tt.concurrentNum; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						logger.ClearRequestUUID()
					}()
				}
				wg.Wait()
			} else {
				logger.ClearRequestUUID()
			}

			if tt.validateAfter != nil {
				tt.validateAfter(t, logger)
			}

			logger.mu.Lock()
			logger.mu.Unlock()
		})
	}
}

func TestSetRequestUUID(t *testing.T) {
	tests := []struct {
		name          string
		uuidString    string
		setup         func(*Logger)
		concurrent    bool
		concurrentNum int
		validate      func(*testing.T, *Logger)
	}{
		{
			name:       "Set Standard UUID",
			uuidString: "123e4567-e89b-12d3-a456-426614174000",
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, "123e4567-e89b-12d3-a456-426614174000", l.requestUUID)
			},
		},
		{
			name:       "Set Empty UUID",
			uuidString: "",
			validate: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name:       "Set UUID with Special Characters",
			uuidString: "test!@#$%^&*()_+-=[]{}|;:,.<>?",
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, "test!@#$%^&*()_+-=[]{}|;:,.<>?", l.requestUUID)
			},
		},
		{
			name:       "Set Very Long UUID",
			uuidString: strings.Repeat("a", 10000),
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, strings.Repeat("a", 10000), l.requestUUID)
			},
		},
		{
			name:       "Set Unicode UUID",
			uuidString: "你好👋Привет",
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, "你好👋Привет", l.requestUUID)
			},
		},
		{
			name:          "Concurrent UUID Setting",
			uuidString:    "concurrent-test",
			concurrent:    true,
			concurrentNum: 100,
			validate: func(t *testing.T, l *Logger) {

				assert.Regexp(t, `^concurrent-test-\d+$`, l.requestUUID)
			},
		},
		{
			name: "Overwrite Existing UUID",
			setup: func(l *Logger) {
				l.SetRequestUUID("existing-uuid")
			},
			uuidString: "new-uuid",
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, "new-uuid", l.requestUUID)
			},
		},
		{
			name:       "Set Null Characters",
			uuidString: "test\x00uuid\x00",
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, "test\x00uuid\x00", l.requestUUID)
			},
		},
		{
			name:       "Set Multiple Times",
			uuidString: "final-uuid",
			setup: func(l *Logger) {
				for i := 0; i < 1000; i++ {
					l.SetRequestUUID(fmt.Sprintf("uuid-%d", i))
				}
			},
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, "final-uuid", l.requestUUID)
			},
		},
		{
			name:       "Performance Test",
			uuidString: "performance-test",
			setup: func(l *Logger) {
				start := time.Now()
				for i := 0; i < 10000; i++ {
					l.SetRequestUUID(fmt.Sprintf("perf-uuid-%d", i))
				}
				duration := time.Since(start)
				assert.Less(t, duration, 1*time.Second, "Performance test took too long")
			},
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, "performance-test", l.requestUUID)
			},
		},
		{
			name:       "Standard UUID Input",
			uuidString: "550e8400-e29b-41d4-a716-446655440000",
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", l.requestUUID)
			},
		},
		{
			name:       "Empty String Input",
			uuidString: "",
			setup: func(l *Logger) {
				l.SetRequestUUID("previous-uuid")
			},
			validate: func(t *testing.T, l *Logger) {
				assert.Empty(t, l.requestUUID)
			},
		},
		{
			name:       "Maximum Length String",
			uuidString: strings.Repeat("x", 65536),
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, strings.Repeat("x", 65536), l.requestUUID)
			},
		},
		{
			name:       "Special Characters in UUID",
			uuidString: "~!@#$%^&*()_+`-={}[]|\\:;\"'<>,.?/",
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, "~!@#$%^&*()_+`-={}[]|\\:;\"'<>,.?/", l.requestUUID)
			},
		},
		{
			name:       "Repeated Calls with Same UUID",
			uuidString: "repeat-uuid",
			setup: func(l *Logger) {
				for i := 0; i < 100; i++ {
					l.SetRequestUUID("repeat-uuid")
				}
			},
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, "repeat-uuid", l.requestUUID)
			},
		},
		{
			name:       "UUID with Leading and Trailing Spaces",
			uuidString: "   space-uuid   ",
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, "   space-uuid   ", l.requestUUID)
			},
		},
		{
			name:       "Unicode Characters in UUID",
			uuidString: "🌟星🌙月☀️日⭐",
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, "🌟星🌙月☀️日⭐", l.requestUUID)
			},
		},
		{
			name:       "Null Character in UUID",
			uuidString: "before\x00after",
			validate: func(t *testing.T, l *Logger) {
				assert.Equal(t, "before\x00after", l.requestUUID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &Logger{
				mu: sync.Mutex{},
			}

			if tt.setup != nil {
				tt.setup(logger)
			}

			if tt.concurrent {
				var wg sync.WaitGroup
				var counter int32
				for i := 0; i < tt.concurrentNum; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						index := atomic.AddInt32(&counter, 1) - 1
						logger.SetRequestUUID(fmt.Sprintf("%s-%d", tt.uuidString, index))
					}()
				}
				wg.Wait()
			} else {
				logger.SetRequestUUID(tt.uuidString)
			}

			if tt.validate != nil {
				tt.validate(t, logger)
			}

			logger.mu.Lock()
			logger.mu.Unlock()
		})
	}
}

func TestRouteBasedUUIDMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		setupHandler   func() http.Handler
		validateUUID   bool
		expectedStatus int
	}{
		{
			name: "Standard GET Request",
			setupRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/test", nil)
			},
			setupHandler: func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.NotEmpty(t, Log.requestUUID)
					w.WriteHeader(http.StatusOK)
				})
			},
			validateUUID:   true,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Request with Existing UUID Header",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/test", nil)
				req.Header.Set("X-Request-ID", "existing-uuid")
				return req
			},
			setupHandler: func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.NotEmpty(t, Log.requestUUID)
					w.WriteHeader(http.StatusOK)
				})
			},
			validateUUID:   true,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Request with Large Body",
			setupRequest: func() *http.Request {
				body := strings.NewReader(strings.Repeat("a", 1024*1024))
				return httptest.NewRequest(http.MethodPost, "/test", body)
			},
			setupHandler: func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.NotEmpty(t, Log.requestUUID)
					w.WriteHeader(http.StatusOK)
				})
			},
			validateUUID:   true,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Request with Special Characters in URL",
			setupRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/test!@#$%^&*()", nil)
			},
			setupHandler: func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.NotEmpty(t, Log.requestUUID)
					w.WriteHeader(http.StatusOK)
				})
			},
			validateUUID:   true,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Long Running Request",
			setupRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/test", nil)
			},
			setupHandler: func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.NotEmpty(t, Log.requestUUID)
					time.Sleep(100 * time.Millisecond)
					w.WriteHeader(http.StatusOK)
				})
			},
			validateUUID:   true,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Standard Request Handling",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
				req.Header.Set("Accept", "application/json")
				return req
			},
			setupHandler: func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.NotEmpty(t, Log.requestUUID)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
				})
			},
			validateUUID:   true,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Empty Request",
			setupRequest: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "", nil)
			},
			setupHandler: func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.NotEmpty(t, Log.requestUUID)
					w.WriteHeader(http.StatusOK)
				})
			},
			validateUUID:   true,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Request with Headers",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/api/test", nil)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer test-token")
				req.Header.Set("X-Custom-Header", "test-value")
				req.Header.Set("Accept-Language", "en-US")
				return req
			},
			setupHandler: func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.NotEmpty(t, Log.requestUUID)
					assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
					assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
					assert.Equal(t, "test-value", r.Header.Get("X-Custom-Header"))
					w.WriteHeader(http.StatusOK)
				})
			},
			validateUUID:   true,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Request with Body",
			setupRequest: func() *http.Request {
				body := strings.NewReader(`{"key":"value"}`)
				req := httptest.NewRequest(http.MethodPost, "/api/test", body)
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			setupHandler: func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.NotEmpty(t, Log.requestUUID)

					bodyBytes, err := io.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.Equal(t, `{"key":"value"}`, string(bodyBytes))

					w.WriteHeader(http.StatusOK)
				})
			},
			validateUUID:   true,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Request with Query Parameters",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest(http.MethodGet, "/api/test?param1=value1&param2=value2", nil)
				return req
			},
			setupHandler: func() http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					assert.NotEmpty(t, Log.requestUUID)
					assert.Equal(t, "value1", r.URL.Query().Get("param1"))
					assert.Equal(t, "value2", r.URL.Query().Get("param2"))
					w.WriteHeader(http.StatusOK)
				})
			},
			validateUUID:   true,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			recorder := httptest.NewRecorder()
			handler := RouteBasedUUIDMiddleware(tt.setupHandler())
			done := make(chan bool)

			go func() {
				defer func() {
					if r := recover(); r != nil {
						t.Logf("Panic recovered in test: %v", r)
					}
					done <- true
				}()

				handler.ServeHTTP(recorder, tt.setupRequest())
			}()

			<-done

			assert.Equal(t, tt.expectedStatus, recorder.Code)

			assert.Empty(t, Log.requestUUID, "UUID should be cleared after request")

			if tt.name == "Multiple Concurrent Requests" {
				var wg sync.WaitGroup
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						rec := httptest.NewRecorder()
						handler.ServeHTTP(rec, tt.setupRequest())
						assert.Equal(t, tt.expectedStatus, rec.Code)
					}()
				}
				wg.Wait()
			}
		})
	}
}
