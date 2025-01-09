package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/robfig/cron"
	"github.com/stretchr/testify/assert"
)

type safeCounter struct {
	count          int
	firstExecution bool
	mu             sync.Mutex
}

func (c *safeCounter) increment(started chan<- bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
	if c.firstExecution && started != nil {
		c.firstExecution = false
		started <- true
	}
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
	c.firstExecution = true
}

var counter = &safeCounter{
	firstExecution: true,
}

func mockHandler() {
	counter.increment(nil)
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
			counter.reset()
			c := cron.New()

			started := make(chan bool, 1)

			err := c.AddFunc(tt.schedule, func() {
				counter.increment(started)
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

			count := counter.getCount()
			t.Logf("Schedule: %s, Wait: %v, Executions: %d", tt.schedule, tt.wait, count)

			if count < tt.want {
				time.Sleep(500 * time.Millisecond)
				count = counter.getCount()
			}

			assert.GreaterOrEqual(t, count, tt.want,
				"Expected at least %d executions, got %d", tt.want, count)
		})
	}
}

func TestCronStopAndRestart(t *testing.T) {
	counter.reset()
	c := cron.New()

	started := make(chan bool, 1)

	err := c.AddFunc("@every 100ms", func() {
		counter.increment(started)
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
	initialCount := counter.getCount()

	t.Logf("Initial count after first run: %d", initialCount)
	assert.Greater(t, initialCount, 0, "Should have at least one execution before stopping")

	time.Sleep(300 * time.Millisecond)

	countAfterStop := counter.getCount()
	assert.Equal(t, initialCount, countAfterStop, "Cron should not execute after stopping")

	counter.reset()
	started = make(chan bool, 1)

	c.Start()

	time.Sleep(500 * time.Millisecond)

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
	counter.reset()
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	c := cron.New()
	started := make(chan bool, 1)

	c.AddFunc("@every 50ms", func() {
		counter.increment(started)
	})

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
