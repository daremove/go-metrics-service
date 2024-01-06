package stats

import (
	"math/rand"
	"runtime"
)

type stats struct {
	pollCount int
}

type ReadResult struct {
	Alloc         float64
	BuckHashSys   float64
	Frees         float64
	GCCPUFraction float64
	GCSys         float64
	HeapAlloc     float64
	HeapIdle      float64
	HeapInuse     float64
	HeapObjects   float64
	HeapReleased  float64
	HeapSys       float64
	LastGC        float64
	Lookups       float64
	MCacheInuse   float64
	MCacheSys     float64
	MSpanInuse    float64
	MSpanSys      float64
	Mallocs       float64
	NextGC        float64
	NumForcedGC   float64
	NumGC         float64
	OtherSys      float64
	PauseTotalNs  float64
	StackInuse    float64
	StackSys      float64
	Sys           float64
	TotalAlloc    float64
	RandomValue   float64
	PollCount     int
}

type Stats interface {
	Read() ReadResult
}

func New() Stats {
	return &stats{pollCount: 0}
}

func (s *stats) Read() ReadResult {
	stats := runtime.MemStats{}

	runtime.ReadMemStats(&stats)
	s.pollCount += 1

	return ReadResult{
		Alloc:         float64(stats.Alloc),
		BuckHashSys:   float64(stats.BuckHashSys),
		Frees:         float64(stats.Frees),
		GCCPUFraction: stats.GCCPUFraction,
		GCSys:         float64(stats.GCSys),
		HeapAlloc:     float64(stats.HeapAlloc),
		HeapIdle:      float64(stats.HeapIdle),
		HeapInuse:     float64(stats.HeapInuse),
		HeapObjects:   float64(stats.HeapObjects),
		HeapReleased:  float64(stats.HeapReleased),
		HeapSys:       float64(stats.HeapSys),
		LastGC:        float64(stats.LastGC),
		Lookups:       float64(stats.Lookups),
		MCacheInuse:   float64(stats.MCacheInuse),
		MCacheSys:     float64(stats.MCacheSys),
		MSpanInuse:    float64(stats.MSpanInuse),
		MSpanSys:      float64(stats.MSpanSys),
		Mallocs:       float64(stats.Mallocs),
		NextGC:        float64(stats.NextGC),
		NumForcedGC:   float64(stats.NumForcedGC),
		NumGC:         float64(stats.NumGC),
		OtherSys:      float64(stats.OtherSys),
		PauseTotalNs:  float64(stats.PauseTotalNs),
		StackInuse:    float64(stats.StackInuse),
		StackSys:      float64(stats.StackSys),
		Sys:           float64(stats.Sys),
		TotalAlloc:    float64(stats.TotalAlloc),
		RandomValue:   rand.Float64(),
		PollCount:     s.pollCount,
	}
}
