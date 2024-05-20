package stats

import (
	"testing"
	"time"

	"github.com/shirou/gopsutil/disk"
	"github.com/stretchr/testify/assert"
)

type MockCPUUsageProvider struct{}

func (m *MockCPUUsageProvider) Percent(interval time.Duration, percpu bool) ([]float64, error) {
	return []float64{10.5, 20.3, 30.7}, nil
}

type MockDiskUsageProvider struct{}

func (m *MockDiskUsageProvider) Usage(path string) (*disk.UsageStat, error) {
	return &disk.UsageStat{
		Total: 500 * 1024 * 1024,
		Free:  200 * 1024 * 1024,
	}, nil
}

var (
	cpuProvider  = &MockCPUUsageProvider{}
	diskProvider = &MockDiskUsageProvider{}
)

func TestStats(t *testing.T) {
	t.Run("Should create a new Stats instance", func(t *testing.T) {
		stats := New(cpuProvider, diskProvider)

		assert.NotNil(t, stats)
		assert.Equal(t, 0, stats.pollCount)
	})

	t.Run("Should read memory stats and increment poll count", func(t *testing.T) {
		stats := New(cpuProvider, diskProvider)
		metrics := stats.Read()

		assert.NotNil(t, metrics)
		assert.Equal(t, 1.0, metrics["PollCount"])
		assert.Contains(t, metrics, "Alloc")
		assert.Contains(t, metrics, "TotalAlloc")
		assert.Contains(t, metrics, "RandomValue")
	})

	t.Run("Should read CPU and disk stats using gopsutil", func(t *testing.T) {
		stats := New(cpuProvider, diskProvider)
		metrics, err := stats.ReadGopsUtil()

		assert.NoError(t, err)
		assert.NotNil(t, metrics)
		assert.Equal(t, float64(500*1024*1024), metrics["TotalMemory"])
		assert.Equal(t, float64(200*1024*1024), metrics["FreeMemory"])
		assert.Equal(t, 10.5, metrics["CPUutilization0"])
		assert.Equal(t, 20.3, metrics["CPUutilization1"])
		assert.Equal(t, 30.7, metrics["CPUutilization2"])
	})
}

func TestStats_RealCPUUsageProvider(t *testing.T) {
	t.Run("Should return data about CPU usage", func(t *testing.T) {
		var data interface{}
		data, err := (&RealCPUUsageProvider{}).Percent(0, true)

		assert.NoError(t, err)
		_, ok := data.([]float64)

		assert.True(t, ok)
	})
}

func TestStats_RealDiskUsageProvider(t *testing.T) {
	t.Run("Should return data about disk usage", func(t *testing.T) {
		var data interface{}
		data, err := (&RealDiskUsageProvider{}).Usage("/")

		assert.NoError(t, err)
		_, ok := data.(*disk.UsageStat)

		assert.True(t, ok)
	})
}
