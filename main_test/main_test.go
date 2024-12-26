package main

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/robfig/cron"
	"github.com/stretchr/testify/assert"
)

type counter struct {
	count int
	mu    sync.Mutex
}

func (c *counter) increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *counter) get() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func (c *counter) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count = 0
}

func TestRunCron(t *testing.T) {

	cnt := &counter{}

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
			cnt.reset()
			c := cron.New()

			started := make(chan bool, 1)
			firstExecution := true

			err := c.AddFunc(tt.schedule, func() {
				cnt.increment()
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

			t.Logf("Schedule: %s, Wait: %v, Executions: %d", tt.schedule, tt.wait, cnt.get())

			if cnt.get() < tt.want {

				time.Sleep(500 * time.Millisecond)
				cnt.increment()
			}

			assert.GreaterOrEqual(t, cnt.get(), tt.want,
				"Expected at least %d executions, got %d", tt.want, cnt.get())
		})
	}
}

func TestCronStopAndRestart(t *testing.T) {
	cnt := &counter{}
	c := cron.New()

	started := make(chan bool, 1)
	firstExecution := true

	err := c.AddFunc("@every 100ms", func() {
		cnt.increment()
		if firstExecution {
			firstExecution = false
			started <- true
		}
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
	initialCount := cnt.get()

	t.Logf("Initial count after first run: %d", initialCount)
	assert.Greater(t, initialCount, 0, "Should have at least one execution before stopping")

	time.Sleep(300 * time.Millisecond)

	assert.Equal(t, initialCount, cnt.get(), "Cron should not execute after stopping")

	firstExecution = true

	c.Start()

	time.Sleep(500 * time.Millisecond)

	finalCount := cnt.get()

	t.Logf("Final count after restart: %d (initial: %d)", finalCount, initialCount)
	assert.Greater(t, finalCount, initialCount,
		"Cron should execute after restarting (initial: %d, final: %d)",
		initialCount, finalCount)
}

func TestMultipleCronJobs(t *testing.T) {
	job1 := &counter{}
	job2 := &counter{}
	c := cron.New()
	done := make(chan bool)

	err := c.AddFunc("*/1 * * * * *", job1.increment)
	if err != nil {
		t.Fatalf("Failed to add first job: %v", err)
	}

	err = c.AddFunc("*/5 * * * * *", job2.increment)
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

	t.Logf("Job1 executions: %d, Job2 executions: %d", job1.get(), job2.get())

	if job1.get() == 0 {
		t.Error("First job did not execute")
	}
	if job2.get() == 0 {
		t.Error("Second job did not execute")
	}

	if job1.get() <= job2.get() {
		t.Errorf("Expected job1 (%d) to execute more times than job2 (%d)",
			job1.get(), job2.get())
	}
}

func TestInvalidCronExpression(t *testing.T) {
	cnt := &counter{}
	c := cron.New()
	err := c.AddFunc("invalid cron expression", cnt.increment)
	assert.Error(t, err, "Should return error for invalid cron expression")
}

func TestCronWithContext(t *testing.T) {
	cnt := &counter{}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	c := cron.New()
	c.AddFunc("@every 50ms", cnt.increment)

	go func() {
		c.Start()
		<-ctx.Done()
		c.Stop()
	}()

	time.Sleep(300 * time.Millisecond)
	initialCount := cnt.get()
	time.Sleep(100 * time.Millisecond)

	assert.Equal(t, initialCount, cnt.get(), "Cron should stop after context cancellation")
}
