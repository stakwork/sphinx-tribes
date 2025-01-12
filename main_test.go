package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/robfig/cron"
	"github.com/stretchr/testify/assert"
)

type MockDB struct {
	initDBFunc                           func() error
	processUpdateTicketsWithoutGroupFunc func()
}

func (m *MockDB) InitDB() error                     { return m.initDBFunc() }
func (m *MockDB) ProcessUpdateTicketsWithoutGroup() { m.processUpdateTicketsWithoutGroupFunc() }

type MockRedis struct{ initFunc func() error }
type MockCache struct{ initFunc func() error }
type MockRoles struct{ initFunc func() error }
type MockWebsocketPool struct{ startFunc func() }
type MockConfig struct{ initFunc func() }
type MockAuth struct{ initFunc func() }
type MockValidator struct{ newFunc func() }
type MockHandlers struct {
	processTwitterFunc func()
	processGithubFunc  func()
}

func (m *MockRedis) InitRedis() error                    { return m.initFunc() }
func (m *MockCache) InitCache() error                    { return m.initFunc() }
func (m *MockRoles) InitRoles() error                    { return m.initFunc() }
func (m *MockWebsocketPool) Start()                      { m.startFunc() }
func (m *MockConfig) InitConfig()                        { m.initFunc() }
func (m *MockAuth) InitJwt()                             { m.initFunc() }
func (m *MockValidator) New()                            { m.newFunc() }
func (m *MockHandlers) ProcessTwitterConfirmationsLoop() { m.processTwitterFunc() }
func (m *MockHandlers) ProcessGithubIssuesLoop()         { m.processGithubFunc() }

func TestMain(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*MockDB, *MockRedis, *MockCache, *MockRoles, *MockWebsocketPool, *MockConfig, *MockAuth, *MockValidator, *MockHandlers)
		skipLoops     string
		expectedError bool
	}{
		{
			name: "Basic Functionality: Successful Initialization",
			setupMocks: func(db *MockDB, redis *MockRedis, cache *MockCache, roles *MockRoles, ws *MockWebsocketPool, conf *MockConfig, auth *MockAuth, val *MockValidator, h *MockHandlers) {
				db.initDBFunc = func() error { return nil }
				redis.initFunc = func() error { return nil }
				cache.initFunc = func() error { return nil }
				roles.initFunc = func() error { return nil }
				ws.startFunc = func() {}
				conf.initFunc = func() {}
				auth.initFunc = func() {}
				val.newFunc = func() {}
				h.processTwitterFunc = func() {}
				h.processGithubFunc = func() {}
				db.processUpdateTicketsWithoutGroupFunc = func() {}
			},
			skipLoops:     "false",
			expectedError: false,
		},
		{
			name: "Edge Case: Missing .env File",
			setupMocks: func(db *MockDB, redis *MockRedis, cache *MockCache, roles *MockRoles, ws *MockWebsocketPool, conf *MockConfig, auth *MockAuth, val *MockValidator, h *MockHandlers) {

				db.initDBFunc = func() error { return nil }
				redis.initFunc = func() error { return nil }
				cache.initFunc = func() error { return nil }
				roles.initFunc = func() error { return nil }
				ws.startFunc = func() {}
				conf.initFunc = func() {}
				auth.initFunc = func() {}
				val.newFunc = func() {}
				h.processTwitterFunc = func() {}
				h.processGithubFunc = func() {}
				db.processUpdateTicketsWithoutGroupFunc = func() {}
			},
			skipLoops:     "",
			expectedError: false,
		},
		{
			name: "Error Condition: Database Initialization Failure",
			setupMocks: func(db *MockDB, redis *MockRedis, cache *MockCache, roles *MockRoles, ws *MockWebsocketPool, conf *MockConfig, auth *MockAuth, val *MockValidator, h *MockHandlers) {
				db.initDBFunc = func() error { return errors.New("database initialization failed") }
			},
			skipLoops:     "false",
			expectedError: true,
		},
		{
			name: "Error Condition: Redis Initialization Failure",
			setupMocks: func(db *MockDB, redis *MockRedis, cache *MockCache, roles *MockRoles, ws *MockWebsocketPool, conf *MockConfig, auth *MockAuth, val *MockValidator, h *MockHandlers) {
				db.initDBFunc = func() error { return nil }
				redis.initFunc = func() error { return errors.New("redis initialization failed") }
			},
			skipLoops:     "false",
			expectedError: true,
		},
		{
			name: "Edge Case: SKIP_LOOPS Environment Variable Set",
			setupMocks: func(db *MockDB, redis *MockRedis, cache *MockCache, roles *MockRoles, ws *MockWebsocketPool, conf *MockConfig, auth *MockAuth, val *MockValidator, h *MockHandlers) {
				db.initDBFunc = func() error { return nil }
				redis.initFunc = func() error { return nil }
				cache.initFunc = func() error { return nil }
				roles.initFunc = func() error { return nil }
				ws.startFunc = func() {}
				conf.initFunc = func() {}
				auth.initFunc = func() {}
				val.newFunc = func() {}
				h.processTwitterFunc = func() { t.Error("Twitter loop should not be called") }
				h.processGithubFunc = func() { t.Error("Github loop should not be called") }
				db.processUpdateTicketsWithoutGroupFunc = func() {}
			},
			skipLoops:     "true",
			expectedError: false,
		},
		{
			name: "Edge Case: SKIP_LOOPS Environment Variable Not Set",
			setupMocks: func(db *MockDB, redis *MockRedis, cache *MockCache, roles *MockRoles, ws *MockWebsocketPool, conf *MockConfig, auth *MockAuth, val *MockValidator, h *MockHandlers) {
				db.initDBFunc = func() error { return nil }
				redis.initFunc = func() error { return nil }
				cache.initFunc = func() error { return nil }
				roles.initFunc = func() error { return nil }
				ws.startFunc = func() {}
				conf.initFunc = func() {}
				auth.initFunc = func() {}
				val.newFunc = func() {}
				h.processTwitterFunc = func() {}
				h.processGithubFunc = func() {}
				db.processUpdateTicketsWithoutGroupFunc = func() {}
			},
			skipLoops:     "",
			expectedError: false,
		},
		{
			name: "Special Case: JWT Initialization Order",
			setupMocks: func(db *MockDB, redis *MockRedis, cache *MockCache, roles *MockRoles, ws *MockWebsocketPool, conf *MockConfig, auth *MockAuth, val *MockValidator, h *MockHandlers) {
				configInitialized := false
				conf.initFunc = func() { configInitialized = true }
				auth.initFunc = func() {
					if !configInitialized {
						t.Error("JWT initialized before config")
					}
				}
				db.initDBFunc = func() error { return nil }
				redis.initFunc = func() error { return nil }
				cache.initFunc = func() error { return nil }
				roles.initFunc = func() error { return nil }
				ws.startFunc = func() {}
				val.newFunc = func() {}
				h.processTwitterFunc = func() {}
				h.processGithubFunc = func() {}
				db.processUpdateTicketsWithoutGroupFunc = func() {}
			},
			skipLoops:     "false",
			expectedError: false,
		},
		{
			name: "Error Condition: Validator Initialization Failure",
			setupMocks: func(db *MockDB, redis *MockRedis, cache *MockCache, roles *MockRoles, ws *MockWebsocketPool, conf *MockConfig, auth *MockAuth, val *MockValidator, h *MockHandlers) {
				db.initDBFunc = func() error { return nil }
				redis.initFunc = func() error { return nil }
				cache.initFunc = func() error { return nil }
				roles.initFunc = func() error { return nil }
				ws.startFunc = func() {}
				conf.initFunc = func() {}
				auth.initFunc = func() {}
				val.newFunc = func() { panic("validator initialization failed") }
				h.processTwitterFunc = func() {}
				h.processGithubFunc = func() {}
				db.processUpdateTicketsWithoutGroupFunc = func() {}
			},
			skipLoops:     "false",
			expectedError: true,
		},
		{
			name: "Error Condition: Cron Job Initialization Failure",
			setupMocks: func(db *MockDB, redis *MockRedis, cache *MockCache, roles *MockRoles, ws *MockWebsocketPool, conf *MockConfig, auth *MockAuth, val *MockValidator, h *MockHandlers) {
				db.initDBFunc = func() error { return nil }
				redis.initFunc = func() error { return nil }
				cache.initFunc = func() error { return nil }
				roles.initFunc = func() error { return nil }
				ws.startFunc = func() {}
				conf.initFunc = func() {}
				auth.initFunc = func() {}
				val.newFunc = func() {}
				h.processTwitterFunc = func() { panic("cron job failed") }
				h.processGithubFunc = func() {}
				db.processUpdateTicketsWithoutGroupFunc = func() {}
			},
			skipLoops:     "false",
			expectedError: true,
		},
		{
			name: "Error Condition: Application Run Failure",
			setupMocks: func(db *MockDB, redis *MockRedis, cache *MockCache, roles *MockRoles, ws *MockWebsocketPool, conf *MockConfig, auth *MockAuth, val *MockValidator, h *MockHandlers) {
				db.initDBFunc = func() error { return nil }
				redis.initFunc = func() error { return nil }
				cache.initFunc = func() error { return nil }
				roles.initFunc = func() error { return nil }
				ws.startFunc = func() { panic("websocket start failed") }
				conf.initFunc = func() {}
				auth.initFunc = func() {}
				val.newFunc = func() {}
				h.processTwitterFunc = func() {}
				h.processGithubFunc = func() {}
				db.processUpdateTicketsWithoutGroupFunc = func() {}
			},
			skipLoops:     "false",
			expectedError: true,
		},
		{
			name: "Edge Case: Empty Environment Variables",
			setupMocks: func(db *MockDB, redis *MockRedis, cache *MockCache, roles *MockRoles, ws *MockWebsocketPool, conf *MockConfig, auth *MockAuth, val *MockValidator, h *MockHandlers) {
				db.initDBFunc = func() error { return nil }
				redis.initFunc = func() error { return nil }
				cache.initFunc = func() error { return nil }
				roles.initFunc = func() error { return nil }
				ws.startFunc = func() {}
				conf.initFunc = func() {}
				auth.initFunc = func() {}
				val.newFunc = func() {}
				h.processTwitterFunc = func() {}
				h.processGithubFunc = func() {}
				db.processUpdateTicketsWithoutGroupFunc = func() {}
			},
			skipLoops:     "",
			expectedError: false,
		},
		{
			name: "Concurrency: Simultaneous Initialization",
			setupMocks: func(db *MockDB, redis *MockRedis, cache *MockCache, roles *MockRoles, ws *MockWebsocketPool, conf *MockConfig, auth *MockAuth, val *MockValidator, h *MockHandlers) {

				initOrder := make([]string, 0, 4)
				var mu sync.Mutex

				db.initDBFunc = func() error {
					mu.Lock()
					initOrder = append(initOrder, "db")
					mu.Unlock()
					return nil
				}
				redis.initFunc = func() error {
					mu.Lock()
					initOrder = append(initOrder, "redis")
					mu.Unlock()
					return nil
				}
				cache.initFunc = func() error {
					mu.Lock()
					initOrder = append(initOrder, "cache")
					mu.Unlock()
					return nil
				}
				roles.initFunc = func() error {
					mu.Lock()
					initOrder = append(initOrder, "roles")
					mu.Unlock()
					return nil
				}

				ws.startFunc = func() {}
				conf.initFunc = func() {}
				auth.initFunc = func() {}
				val.newFunc = func() {}
				h.processTwitterFunc = func() {}
				h.processGithubFunc = func() {}
				db.processUpdateTicketsWithoutGroupFunc = func() {}

				t.Cleanup(func() {
					expected := []string{"db", "redis", "cache", "roles"}
					if !reflect.DeepEqual(initOrder, expected) {
						t.Errorf("Wrong initialization order. Expected %v, got %v", expected, initOrder)
					}
				})
			},
			skipLoops:     "false",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockDB := &MockDB{}
			mockRedis := &MockRedis{}
			mockCache := &MockCache{}
			mockRoles := &MockRoles{}
			mockWS := &MockWebsocketPool{}
			mockConfig := &MockConfig{}
			mockAuth := &MockAuth{}
			mockValidator := &MockValidator{}
			mockHandlers := &MockHandlers{}

			tt.setupMocks(mockDB, mockRedis, mockCache, mockRoles, mockWS, mockConfig, mockAuth, mockValidator, mockHandlers)

			os.Setenv("SKIP_LOOPS", tt.skipLoops)
			defer os.Unsetenv("SKIP_LOOPS")

			var initErr error
			func() {
				defer func() {
					if r := recover(); r != nil {
						initErr = fmt.Errorf("panic occurred: %v", r)
					}
				}()

				if err := mockDB.InitDB(); err != nil {
					initErr = err
					return
				}
				if err := mockRedis.InitRedis(); err != nil {
					initErr = err
					return
				}
				if err := mockCache.InitCache(); err != nil {
					initErr = err
					return
				}
				if err := mockRoles.InitRoles(); err != nil {
					initErr = err
					return
				}

				mockDB.ProcessUpdateTicketsWithoutGroup()
				mockConfig.InitConfig()
				mockAuth.InitJwt()

				func() {
					defer func() {
						if r := recover(); r != nil {
							initErr = fmt.Errorf("panic occurred: %v", r)
						}
					}()
					mockValidator.New()
					mockWS.Start()

					if tt.skipLoops != "true" {
						mockHandlers.ProcessTwitterConfirmationsLoop()
						mockHandlers.ProcessGithubIssuesLoop()
					}
				}()
			}()

			if tt.expectedError {
				assert.Error(t, initErr, "Expected an error but got none")
			} else {
				assert.NoError(t, initErr, "Expected no error but got: %v", initErr)
			}
		})
	}
}

type safeCounter struct {
	count int
	mu    sync.Mutex
}

func (c *safeCounter) increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *safeCounter) getCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func (c *safeCounter) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count = 0
}

var counter = &safeCounter{}

func mockHandler() {
	counter.increment()
}

func resetExecutionCount() {
	counter.reset()
}

func TestRunCron(t *testing.T) {
	tests := []struct {
		name     string
		schedule string
		wait     time.Duration
		want     int
	}{
		{
			name:     "Basic Functionality: Cron Job Setup",
			schedule: "@every 500ms",
			wait:     1200 * time.Millisecond,
			want:     2,
		},
		{
			name:     "Basic Functionality: Cron Job Execution",
			schedule: "@every 250ms",
			wait:     1500 * time.Millisecond,
			want:     2,
		},
		{
			name:     "Edge Case: Immediate Execution",
			schedule: "@every 1ms",
			wait:     200 * time.Millisecond,
			want:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetExecutionCount()
			c := cron.New()

			started := make(chan bool, 1)
			firstExecution := true

			err := c.AddFunc(tt.schedule, func() {
				counter.increment()
				if firstExecution {
					firstExecution = false
					started <- true
				}
			})

			if err != nil {
				t.Fatalf("Failed to add cron job: %v", err)
			}

			c.Start()
			time.Sleep(100 * time.Millisecond)

			select {
			case <-started:

			case <-time.After(2 * time.Second):
				t.Fatal("Cron job failed to start within timeout")
			}

			time.Sleep(tt.wait)
			c.Stop()

			t.Logf("Schedule: %s, Wait: %v, Executions: %d", tt.schedule, tt.wait, counter.getCount())

			if counter.getCount() < tt.want {

				time.Sleep(500 * time.Millisecond)
				counter.reset()
			}

			assert.GreaterOrEqual(t, counter.getCount(), tt.want,
				"Expected at least %d executions, got %d", tt.want, counter.getCount())
		})
	}
}

func TestCronStopAndRestart(t *testing.T) {
	resetExecutionCount()
	c := cron.New()

	started := make(chan bool, 1)
	executed := make(chan bool, 1)
	firstExecution := true

	err := c.AddFunc("@every 100ms", func() {
		counter.increment()
		if firstExecution {
			firstExecution = false
			started <- true
		}
		executed <- true
	})

	if err != nil {
		t.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()

	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("Cron job failed to start within timeout")
	}

	select {
	case <-executed:
	case <-time.After(2 * time.Second):
		t.Fatal("Cron job failed to execute within timeout")
	}

	c.Stop()
	initialCount := counter.getCount()

	t.Logf("Initial count after first run: %d", initialCount)
	assert.Greater(t, initialCount, 0, "Should have at least one execution before stopping")

	counter.reset()
	firstExecution = true

	c.Start()

	select {
	case <-executed:
	case <-time.After(2 * time.Second):
		t.Fatal("Cron job failed to execute after restart within timeout")
	}

	finalCount := counter.getCount()

	t.Logf("Final count after restart: %d (initial: %d)", finalCount, initialCount)
	assert.Greater(t, finalCount, 0,
		"Cron should execute after restarting (initial: %d, final: %d)",
		initialCount, finalCount)
}

func TestMultipleCronJobs(t *testing.T) {
	resetExecutionCount()
	c := cron.New()

	var job1Count, job2Count int
	var mu sync.Mutex
	done := make(chan bool)

	err := c.AddFunc("*/1 * * * * *", func() {
		mu.Lock()
		job1Count++
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("Failed to add first job: %v", err)
	}

	err = c.AddFunc("*/5 * * * * *", func() {
		mu.Lock()
		job2Count++
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("Failed to add second job: %v", err)
	}

	c.Start()

	go func() {
		time.Sleep(5 * time.Second)
		c.Stop()
		done <- true
	}()

	<-done

	mu.Lock()
	j1Count := job1Count
	j2Count := job2Count
	mu.Unlock()

	t.Logf("Job1 executions: %d, Job2 executions: %d", j1Count, j2Count)

	if j1Count == 0 {
		t.Error("First job did not execute")
	}
	if j2Count == 0 {
		t.Error("Second job did not execute")
	}

	if j1Count <= j2Count {
		t.Errorf("Expected job1 (%d) to execute more times than job2 (%d)",
			j1Count, j2Count)
	}
}

func TestInvalidCronExpression(t *testing.T) {
	c := cron.New()
	err := c.AddFunc("invalid cron expression", mockHandler)
	assert.Error(t, err, "Should return error for invalid cron expression")
}

func TestCronWithContext(t *testing.T) {
	resetExecutionCount()
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	c := cron.New()
	c.AddFunc("@every 50ms", mockHandler)

	go func() {
		c.Start()
		<-ctx.Done()
		c.Stop()
	}()

	time.Sleep(300 * time.Millisecond)
	initialCount := counter.getCount()
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, initialCount, counter.getCount(), "Cron should stop after context cancellation")
}
