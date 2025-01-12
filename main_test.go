package main

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"

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
