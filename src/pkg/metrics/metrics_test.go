package metrics

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMetrics(t *testing.T) {
	metrics1 := GetMetrics()
	metrics2 := GetMetrics()

	// Should return the same instance (singleton pattern)
	assert.Equal(t, metrics1, metrics2)
	assert.NotNil(t, metrics1)
	assert.False(t, metrics1.StartTime.IsZero())
}

func TestMetricsRecordEncrypt(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	duration := 100 * time.Millisecond
	metrics.RecordEncrypt(duration, true)

	assert.Equal(t, int64(1), metrics.EncryptCount)
	assert.Equal(t, duration, metrics.EncryptDuration)
	assert.Equal(t, duration, metrics.AverageEncryptTime)
	assert.Equal(t, int64(0), metrics.EncryptErrors)
}

func TestMetricsRecordEncryptFailure(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	duration := 100 * time.Millisecond
	metrics.RecordEncrypt(duration, false)

	assert.Equal(t, int64(1), metrics.EncryptCount)
	assert.Equal(t, duration, metrics.EncryptDuration)
	assert.Equal(t, int64(1), metrics.EncryptErrors)
}

func TestMetricsRecordDecrypt(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	duration := 150 * time.Millisecond
	metrics.RecordDecrypt(duration, true)

	assert.Equal(t, int64(1), metrics.DecryptCount)
	assert.Equal(t, duration, metrics.DecryptDuration)
	assert.Equal(t, duration, metrics.AverageDecryptTime)
	assert.Equal(t, int64(0), metrics.DecryptErrors)
}

func TestMetricsRecordDecryptFailure(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	duration := 150 * time.Millisecond
	metrics.RecordDecrypt(duration, false)

	assert.Equal(t, int64(1), metrics.DecryptCount)
	assert.Equal(t, duration, metrics.DecryptDuration)
	assert.Equal(t, int64(1), metrics.DecryptErrors)
}

func TestMetricsRecordIntermediate(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	duration := 200 * time.Millisecond
	metrics.RecordIntermediate(duration, true)

	assert.Equal(t, int64(1), metrics.IntermediateCount)
	assert.Equal(t, duration, metrics.IntermediateDuration)
	assert.Equal(t, duration, metrics.AverageIntermediateTime)
	assert.Equal(t, int64(0), metrics.IntermediateErrors)
}

func TestMetricsRecordIntermediateFailure(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	duration := 200 * time.Millisecond
	metrics.RecordIntermediate(duration, false)

	assert.Equal(t, int64(1), metrics.IntermediateCount)
	assert.Equal(t, duration, metrics.IntermediateDuration)
	assert.Equal(t, int64(1), metrics.IntermediateErrors)
}

func TestMetricsMultipleOperations(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	// Record multiple encrypt operations
	durations := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		300 * time.Millisecond,
	}
	for _, duration := range durations {
		metrics.RecordEncrypt(duration, true)
	}

	assert.Equal(t, int64(3), metrics.EncryptCount)
	assert.Equal(t, 600*time.Millisecond, metrics.EncryptDuration)
	assert.Equal(t, 200*time.Millisecond, metrics.AverageEncryptTime)
}

func TestMetricsUpdateUptime(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	metrics.UpdateUptime()

	assert.Greater(t, metrics.Uptime, time.Duration(0))
	assert.Less(t, metrics.Uptime, 100*time.Millisecond)
}

func TestMetricsGetSnapshot(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	metrics.RecordEncrypt(100*time.Millisecond, true)
	metrics.RecordDecrypt(150*time.Millisecond, false)

	snapshot := metrics.GetSnapshot()

	// Verify snapshot contains the data
	assert.Equal(t, int64(1), snapshot.EncryptCount)
	assert.Equal(t, int64(1), snapshot.DecryptCount)
	assert.Equal(t, int64(1), snapshot.DecryptErrors)
	assert.Greater(t, snapshot.Uptime, time.Duration(0))

	// Verify it's a copy (modifying original shouldn't affect snapshot)
	metrics.RecordEncrypt(50*time.Millisecond, true)
	assert.Equal(t, int64(1), snapshot.EncryptCount) // Should still be 1
}

func TestMetricsReset(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	// Add some data
	metrics.RecordEncrypt(100*time.Millisecond, true)
	metrics.RecordDecrypt(150*time.Millisecond, false)
	metrics.RecordIntermediate(200*time.Millisecond, true)

	// Reset
	metrics.Reset()

	// Verify all values are reset except StartTime
	assert.Equal(t, int64(0), metrics.EncryptCount)
	assert.Equal(t, int64(0), metrics.DecryptCount)
	assert.Equal(t, int64(0), metrics.IntermediateCount)
	assert.Equal(t, int64(0), metrics.EncryptErrors)
	assert.Equal(t, int64(0), metrics.DecryptErrors)
	assert.Equal(t, int64(0), metrics.IntermediateErrors)
	assert.Equal(t, time.Duration(0), metrics.EncryptDuration)
	assert.Equal(t, time.Duration(0), metrics.DecryptDuration)
	assert.Equal(t, time.Duration(0), metrics.IntermediateDuration)
	assert.False(t, metrics.StartTime.IsZero())
}

func TestMetricsSuccessRates(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	// Test with no operations
	assert.Equal(t, float64(0), metrics.EncryptSuccessRate())
	assert.Equal(t, float64(0), metrics.DecryptSuccessRate())
	assert.Equal(t, float64(0), metrics.IntermediateSuccessRate())

	// Add some successful operations
	metrics.RecordEncrypt(100*time.Millisecond, true)
	metrics.RecordEncrypt(100*time.Millisecond, true)
	metrics.RecordEncrypt(100*time.Millisecond, false) // One failure

	assert.Equal(t, float64(66.67), metrics.EncryptSuccessRate())

	// Test decrypt success rate
	metrics.RecordDecrypt(150*time.Millisecond, true)
	assert.Equal(t, float64(100), metrics.DecryptSuccessRate())

	// Test intermediate success rate
	metrics.RecordIntermediate(200*time.Millisecond, true)
	metrics.RecordIntermediate(200*time.Millisecond, false)
	assert.Equal(t, float64(50), metrics.IntermediateSuccessRate())
}

func TestNewTimer(t *testing.T) {
	timer := NewTimer("encrypt")

	assert.NotNil(t, timer)
	assert.Equal(t, "encrypt", timer.operation)
	assert.NotNil(t, timer.metrics)
	assert.False(t, timer.start.IsZero())
}

func TestTimerStop(t *testing.T) {
	// Reset global metrics
	Reset()

	timer := NewTimer("encrypt")
	time.Sleep(10 * time.Millisecond)
	timer.Stop(true)

	snapshot := GetSnapshot()
	assert.Equal(t, int64(1), snapshot.EncryptCount)
	assert.Equal(t, int64(0), snapshot.EncryptErrors)
	assert.Greater(t, snapshot.AverageEncryptTime, time.Duration(0))
}

func TestTimerStopFailure(t *testing.T) {
	// Reset global metrics
	Reset()

	timer := NewTimer("decrypt")
	time.Sleep(10 * time.Millisecond)
	timer.Stop(false)

	snapshot := GetSnapshot()
	assert.Equal(t, int64(1), snapshot.DecryptCount)
	assert.Equal(t, int64(1), snapshot.DecryptErrors)
}

func TestNewContextTimer(t *testing.T) {
	ctx := context.Background()
	timer := NewContextTimer(ctx, "encrypt")

	assert.NotNil(t, timer)
	assert.NotNil(t, timer.ctx)
	assert.Equal(t, "encrypt", timer.operation)
}

func TestContextTimerStop(t *testing.T) {
	// Reset global metrics
	Reset()

	ctx := context.Background()
	timer := NewContextTimer(ctx, "intermediate")
	time.Sleep(10 * time.Millisecond)
	timer.Stop(true)

	snapshot := GetSnapshot()
	assert.Equal(t, int64(1), snapshot.IntermediateCount)
	assert.Equal(t, int64(0), snapshot.IntermediateErrors)
}

func TestContextTimerStopCancelled(t *testing.T) {
	// Reset global metrics
	Reset()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	timer := NewContextTimer(ctx, "encrypt")
	time.Sleep(10 * time.Millisecond)
	timer.Stop(true) // Should not record due to cancellation

	snapshot := GetSnapshot()
	assert.Equal(t, int64(0), snapshot.EncryptCount) // Should not be recorded
}

func TestMetricsConcurrency(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	var wg sync.WaitGroup
	numGoroutines := 50
	operationsPerGoroutine := 10

	// Test concurrent operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(3)

		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				metrics.RecordEncrypt(time.Duration(j)*time.Millisecond, true)
			}
		}()

		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				metrics.RecordDecrypt(time.Duration(j)*time.Millisecond, false)
			}
		}()

		go func() {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				metrics.RecordIntermediate(time.Duration(j)*time.Millisecond, true)
			}
		}()
	}

	wg.Wait()

	expectedCount := int64(numGoroutines * operationsPerGoroutine)
	assert.Equal(t, expectedCount, metrics.EncryptCount)
	assert.Equal(t, expectedCount, metrics.DecryptCount)
	assert.Equal(t, expectedCount, metrics.IntermediateCount)
	assert.Equal(t, expectedCount, metrics.DecryptErrors) // All decrypt operations failed
	assert.Equal(t, int64(0), metrics.EncryptErrors)      // All encrypt operations succeeded
	assert.Equal(t, int64(0), metrics.IntermediateErrors) // All intermediate operations succeeded
}

func TestGlobalHelperFunctions(t *testing.T) {
	// Reset global metrics
	Reset()

	// Test global helper functions
	RecordEncrypt(100*time.Millisecond, true)
	RecordDecrypt(150*time.Millisecond, false)
	RecordIntermediate(200*time.Millisecond, true)

	snapshot := GetSnapshot()
	assert.Equal(t, int64(1), snapshot.EncryptCount)
	assert.Equal(t, int64(1), snapshot.DecryptCount)
	assert.Equal(t, int64(1), snapshot.IntermediateCount)
	assert.Equal(t, int64(1), snapshot.DecryptErrors)
}

func TestMetricsEdgeCases(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	// Test success rate with zero operations
	assert.Equal(t, float64(0), metrics.EncryptSuccessRate())
	assert.Equal(t, float64(0), metrics.DecryptSuccessRate())
	assert.Equal(t, float64(0), metrics.IntermediateSuccessRate())

	// Test average calculation with single operation
	metrics.RecordEncrypt(100*time.Millisecond, true)
	assert.Equal(t, 100*time.Millisecond, metrics.AverageEncryptTime)
}

func TestMetricsPerformance(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	// Test performance with many operations
	numOperations := 1000

	start := time.Now()
	for i := 0; i < numOperations; i++ {
		metrics.RecordEncrypt(time.Duration(i)*time.Microsecond, true)
	}
	duration := time.Since(start)

	// Should complete quickly
	assert.Less(t, duration, 100*time.Millisecond)
	assert.Equal(t, int64(numOperations), metrics.EncryptCount)
}

func TestMetricsThreadSafety(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	// Test that the metrics struct is properly protected
	var wg sync.WaitGroup
	done := make(chan bool, 1)

	// Start a goroutine that continuously reads metrics
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				_ = metrics.GetSnapshot()
				_ = metrics.EncryptSuccessRate()
				_ = metrics.DecryptSuccessRate()
				_ = metrics.IntermediateSuccessRate()
			}
		}
	}()

	// Start multiple goroutines that write metrics
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				metrics.RecordEncrypt(time.Duration(j)*time.Millisecond, true)
				metrics.RecordDecrypt(time.Duration(j)*time.Millisecond, false)
				metrics.RecordIntermediate(time.Duration(j)*time.Millisecond, true)
			}
		}()
	}

	wg.Wait()
	close(done)

	// If we reach here without race conditions, the test passes
	assert.True(t, true)
}

func TestMetricsValidation(t *testing.T) {
	metrics := &Metrics{StartTime: time.Now()}

	require.NotNil(t, metrics)

	// Test that all methods return sensible values
	snapshot := metrics.GetSnapshot()
	assert.GreaterOrEqual(t, snapshot.EncryptCount, int64(0))
	assert.GreaterOrEqual(t, snapshot.DecryptCount, int64(0))
	assert.GreaterOrEqual(t, snapshot.IntermediateCount, int64(0))
	assert.GreaterOrEqual(t, snapshot.EncryptErrors, int64(0))
	assert.GreaterOrEqual(t, snapshot.DecryptErrors, int64(0))
	assert.GreaterOrEqual(t, snapshot.IntermediateErrors, int64(0))
	assert.GreaterOrEqual(t, snapshot.EncryptDuration, time.Duration(0))
	assert.GreaterOrEqual(t, snapshot.DecryptDuration, time.Duration(0))
	assert.GreaterOrEqual(t, snapshot.IntermediateDuration, time.Duration(0))
	assert.GreaterOrEqual(t, snapshot.AverageEncryptTime, time.Duration(0))
	assert.GreaterOrEqual(t, snapshot.AverageDecryptTime, time.Duration(0))
	assert.GreaterOrEqual(t, snapshot.AverageIntermediateTime, time.Duration(0))
	assert.GreaterOrEqual(t, snapshot.Uptime, time.Duration(0))
}

func TestRealWorldScenario(t *testing.T) {
	// Reset global metrics
	Reset()

	// Simulate a real-world scenario
	timer := NewTimer("encrypt")
	time.Sleep(50 * time.Millisecond)
	timer.Stop(true)

	timer = NewTimer("decrypt")
	time.Sleep(75 * time.Millisecond)
	timer.Stop(false)

	timer = NewTimer("intermediate")
	time.Sleep(100 * time.Millisecond)
	timer.Stop(true)

	snapshot := GetSnapshot()
	assert.Equal(t, int64(1), snapshot.EncryptCount)
	assert.Equal(t, int64(1), snapshot.DecryptCount)
	assert.Equal(t, int64(1), snapshot.IntermediateCount)
	assert.Equal(t, int64(1), snapshot.DecryptErrors)
	assert.Equal(t, float64(100), snapshot.EncryptSuccessRate())
	assert.Equal(t, float64(0), snapshot.DecryptSuccessRate())
	assert.Equal(t, float64(100), snapshot.IntermediateSuccessRate())
}
