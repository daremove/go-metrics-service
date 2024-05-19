// Package stats предоставляет функции для чтения системной статистики,
// такой как использование CPU, диск и статистика памяти.
package stats

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"

	"github.com/shirou/gopsutil/disk"
)

// CPUUsageProvider определяет интерфейс для получения статистики CPU.
type CPUUsageProvider interface {
	Percent(interval time.Duration, percpu bool) ([]float64, error)
}

// DiskUsageProvider определяет интерфейс для получения статистики диска.
type DiskUsageProvider interface {
	Usage(path string) (*disk.UsageStat, error)
}

// RealCPUUsageProvider предоставляет реальные данные использования CPU.
type RealCPUUsageProvider struct{}

func (r *RealCPUUsageProvider) Percent(interval time.Duration, percpu bool) ([]float64, error) {
	return cpu.Percent(interval, percpu)
}

// RealDiskUsageProvider предоставляет реальные данные использования диска.
type RealDiskUsageProvider struct{}

func (r *RealDiskUsageProvider) Usage(path string) (*disk.UsageStat, error) {
	return disk.Usage(path)
}

// Stats содержит внутренние данные для отслеживания количества запросов.
type Stats struct {
	pollCount int // Количество запросов к статистике.
	cpu       CPUUsageProvider
	disk      DiskUsageProvider
}

// New создает и возвращает новый экземпляр Stats.
func New(cpu CPUUsageProvider, disk DiskUsageProvider) *Stats {
	return &Stats{pollCount: 0, cpu: cpu, disk: disk}
}

// Read возвращает основные метрики памяти и системные параметры.
func (s *Stats) Read() map[string]float64 {
	stats := runtime.MemStats{}

	runtime.ReadMemStats(&stats)
	s.pollCount += 1

	return map[string]float64{
		"Alloc":         float64(stats.Alloc),
		"BuckHashSys":   float64(stats.BuckHashSys),
		"Frees":         float64(stats.Frees),
		"GCCPUFraction": stats.GCCPUFraction,
		"GCSys":         float64(stats.GCSys),
		"HeapAlloc":     float64(stats.HeapAlloc),
		"HeapIdle":      float64(stats.HeapIdle),
		"HeapInuse":     float64(stats.HeapInuse),
		"HeapObjects":   float64(stats.HeapObjects),
		"HeapReleased":  float64(stats.HeapReleased),
		"HeapSys":       float64(stats.HeapSys),
		"LastGC":        float64(stats.LastGC),
		"Lookups":       float64(stats.Lookups),
		"MCacheInuse":   float64(stats.MCacheInuse),
		"MCacheSys":     float64(stats.MCacheSys),
		"MSpanInuse":    float64(stats.MSpanInuse),
		"MSpanSys":      float64(stats.MSpanSys),
		"Mallocs":       float64(stats.Mallocs),
		"NextGC":        float64(stats.NextGC),
		"NumForcedGC":   float64(stats.NumForcedGC),
		"NumGC":         float64(stats.NumGC),
		"OtherSys":      float64(stats.OtherSys),
		"PauseTotalNs":  float64(stats.PauseTotalNs),
		"StackInuse":    float64(stats.StackInuse),
		"StackSys":      float64(stats.StackSys),
		"Sys":           float64(stats.Sys),
		"TotalAlloc":    float64(stats.TotalAlloc),
		"RandomValue":   rand.Float64(),
		"PollCount":     float64(s.pollCount),
	}
}

// ReadGopsUtil читает и возвращает статистику использования CPU и диска через библиотеку gopsutil.
func (s *Stats) ReadGopsUtil() (map[string]float64, error) {
	cpuPercents, err := s.cpu.Percent(0, true)

	if err != nil {
		return nil, err
	}

	usageStat, err := s.disk.Usage("/")

	if err != nil {
		return nil, err
	}

	stats := map[string]float64{
		"TotalMemory": float64(usageStat.Total),
		"FreeMemory":  float64(usageStat.Free),
	}

	for i, cpuPercent := range cpuPercents {
		stats[fmt.Sprintf("CPUutilization%d", i)] = cpuPercent
	}

	return stats, nil
}
