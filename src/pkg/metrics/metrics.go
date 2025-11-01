package metrics

import (
	"context"
	"math"
	"sync"
	"time"
)

// Metrics holds basic performance metrics for the application
type Metrics struct {
	mu sync.RWMutex

	// Operation metrics
	EncryptCount      int64 `json:"encrypt_count"`
	DecryptCount      int64 `json:"decrypt_count"`
	IntermediateCount int64 `json:"intermediate_count"`

	// Timing metrics
	EncryptDuration      time.Duration `json:"encrypt_duration_total"`
	DecryptDuration      time.Duration `json:"decrypt_duration_total"`
	IntermediateDuration time.Duration `json:"intermediate_duration_total"`

	// Error metrics
	EncryptErrors      int64 `json:"encrypt_errors"`
	DecryptErrors      int64 `json:"decrypt_errors"`
	IntermediateErrors int64 `json:"intermediate_errors"`

	// Performance metrics
	AverageEncryptTime      time.Duration `json:"average_encrypt_time"`
	AverageDecryptTime      time.Duration `json:"average_decrypt_time"`
	AverageIntermediateTime time.Duration `json:"average_intermediate_time"`

	// System metrics
	StartTime time.Time     `json:"start_time"`
	Uptime    time.Duration `json:"uptime"`
}

var (
	// Global metrics instance
	defaultMetrics *Metrics
	once           sync.Once
)

// GetMetrics returns the global metrics instance
func GetMetrics() *Metrics {
	once.Do(func() {
		defaultMetrics = &Metrics{
			StartTime: time.Now(),
		}
	})
	return defaultMetrics
}

// RecordEncrypt records an encryption operation
func (m *Metrics) RecordEncrypt(duration time.Duration, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.EncryptCount++
	m.EncryptDuration += duration

	if success {
		// Update average
		m.AverageEncryptTime = m.EncryptDuration / time.Duration(m.EncryptCount)
	} else {
		m.EncryptErrors++
	}
}

// RecordDecrypt records a decryption operation
func (m *Metrics) RecordDecrypt(duration time.Duration, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.DecryptCount++
	m.DecryptDuration += duration

	if success {
		// Update average
		m.AverageDecryptTime = m.DecryptDuration / time.Duration(m.DecryptCount)
	} else {
		m.DecryptErrors++
	}
}

// RecordIntermediate records an intermediate code generation operation
func (m *Metrics) RecordIntermediate(duration time.Duration, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.IntermediateCount++
	m.IntermediateDuration += duration

	if success {
		// Update average
		m.AverageIntermediateTime = m.IntermediateDuration / time.Duration(m.IntermediateCount)
	} else {
		m.IntermediateErrors++
	}
}

// UpdateUptime updates the uptime duration
func (m *Metrics) UpdateUptime() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Uptime = time.Since(m.StartTime)
}

// GetSnapshot returns a copy of current metrics
func (m *Metrics) GetSnapshot() Metrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Uptime = time.Since(m.StartTime)

	snapshot := *m
	snapshot.mu = sync.RWMutex{}
	return snapshot
}

// Reset resets all metrics to zero
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.EncryptCount = 0
	m.DecryptCount = 0
	m.IntermediateCount = 0

	m.EncryptDuration = 0
	m.DecryptDuration = 0
	m.IntermediateDuration = 0

	m.EncryptErrors = 0
	m.DecryptErrors = 0
	m.IntermediateErrors = 0

	m.AverageEncryptTime = 0
	m.AverageDecryptTime = 0
	m.AverageIntermediateTime = 0

	m.StartTime = time.Now()
	m.Uptime = 0
}

// SuccessRate returns the success rate for encrypt operations (0-100)
func (m *Metrics) EncryptSuccessRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.EncryptCount == 0 {
		return 0
	}
	return roundToTwoDecimals(float64(m.EncryptCount-m.EncryptErrors) /
		float64(m.EncryptCount) * 100)
}

// SuccessRate returns the success rate for decrypt operations (0-100)
func (m *Metrics) DecryptSuccessRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.DecryptCount == 0 {
		return 0
	}
	return roundToTwoDecimals(float64(m.DecryptCount-m.DecryptErrors) /
		float64(m.DecryptCount) * 100)
}

// SuccessRate returns the success rate for intermediate operations (0-100)
func (m *Metrics) IntermediateSuccessRate() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.IntermediateCount == 0 {
		return 0
	}
	return roundToTwoDecimals(float64(m.IntermediateCount-m.IntermediateErrors) /
		float64(m.IntermediateCount) * 100)
}

func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

// Timer is a helper for timing operations
type Timer struct {
	start     time.Time
	metrics   *Metrics
	operation string
}

// NewTimer creates a new timer for the given operation
func NewTimer(operation string) *Timer {
	return &Timer{
		start:     time.Now(),
		metrics:   GetMetrics(),
		operation: operation,
	}
}

// Stop stops the timer and records the metric
func (t *Timer) Stop(success bool) {
	duration := time.Since(t.start)

	switch t.operation {
	case "encrypt":
		t.metrics.RecordEncrypt(duration, success)
	case "decrypt":
		t.metrics.RecordDecrypt(duration, success)
	case "intermediate":
		t.metrics.RecordIntermediate(duration, success)
	}
}

// ContextTimer is a timer that works with context cancellation
type ContextTimer struct {
	*Timer
	ctx context.Context
}

// NewContextTimer creates a new timer with context
func NewContextTimer(ctx context.Context, operation string) *ContextTimer {
	return &ContextTimer{
		Timer: NewTimer(operation),
		ctx:   ctx,
	}
}

// Stop stops the timer and records the metric if context hasn't been cancelled
func (t *ContextTimer) Stop(success bool) {
	select {
	case <-t.ctx.Done():
		// Context was cancelled, don't record the metric
		return
	default:
		t.Timer.Stop(success)
	}
}

// Helper functions for global metrics instance

// RecordEncrypt records an encryption operation using the global metrics instance
func RecordEncrypt(duration time.Duration, success bool) {
	GetMetrics().RecordEncrypt(duration, success)
}

// RecordDecrypt records a decryption operation using the global metrics instance
func RecordDecrypt(duration time.Duration, success bool) {
	GetMetrics().RecordDecrypt(duration, success)
}

// RecordIntermediate records an intermediate operation using the global metrics instance
func RecordIntermediate(duration time.Duration, success bool) {
	GetMetrics().RecordIntermediate(duration, success)
}

// GetSnapshot returns a snapshot of the global metrics
func GetSnapshot() Metrics {
	return GetMetrics().GetSnapshot()
}

// Reset resets the global metrics
func Reset() {
	GetMetrics().Reset()
}
