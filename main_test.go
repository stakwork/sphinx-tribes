package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/robfig/cron"
	"github.com/stretchr/testify/assert"
)

var executionCount = 0
var mutex sync.Mutex

func mockHandler() {
	mutex.Lock()
	executionCount++
	mutex.Unlock()
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
				mutex.Lock()
				defer mutex.Unlock()
				executionCount++
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

			mutex.Lock()
			count := executionCount
			mutex.Unlock()

			t.Logf("Schedule: %s, Wait: %v, Executions: %d", tt.schedule, tt.wait, count)

			if count < tt.want {

				time.Sleep(500 * time.Millisecond)
				mutex.Lock()
				count = executionCount
				mutex.Unlock()
			}

			assert.GreaterOrEqual(t, count, tt.want,
				"Expected at least %d executions, got %d", tt.want, count)
		})
	}
}

func resetExecutionCount() {
	mutex.Lock()
	defer mutex.Unlock()
	executionCount = 0
}

func TestCronStopAndRestart(t *testing.T) {
	resetExecutionCount()
	c := cron.New()

	started := make(chan bool, 1)
	firstExecution := true

	err := c.AddFunc("@every 100ms", func() {
		mutex.Lock()
		executionCount++
		if firstExecution {
			firstExecution = false
			started <- true
		}
		mutex.Unlock()
	})

	if err != nil {
		t.Fatalf("Failed to add cron job: %v", err)
	}

	c.Start()

	select {
	case <-started:

	case <-time.After(1 * time.Second):
		t.Fatal("Cron job failed to start within timeout")
	}

	time.Sleep(300 * time.Millisecond)

	c.Stop()
	mutex.Lock()
	initialCount := executionCount
	mutex.Unlock()

	t.Logf("Initial count after first run: %d", initialCount)
	assert.Greater(t, initialCount, 0, "Should have at least one execution before stopping")

	time.Sleep(300 * time.Millisecond)

	mutex.Lock()
	countAfterStop := executionCount
	mutex.Unlock()

	assert.Equal(t, initialCount, countAfterStop, "Cron should not execute after stopping")

	firstExecution = true

	c.Start()

	time.Sleep(500 * time.Millisecond)

	mutex.Lock()
	finalCount := executionCount
	mutex.Unlock()

	t.Logf("Final count after restart: %d (initial: %d)", finalCount, initialCount)
	assert.Greater(t, finalCount, initialCount,
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
	initialCount := executionCount
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, initialCount, executionCount, "Cron should stop after context cancellation")
}
