package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime"
	"time"
	"math"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
	serverAddress  = "http://localhost:8080/update/"
)

var (
	metrics     map[string]uint64
	PollCount   uint64 = 0
	RandomValue uint64 = 0
)

func main() {
	ticker  := time.NewTicker(reportInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sendMetrics()
		default:
			collectMetrics()
			time.Sleep(pollInterval)
		}
	}
}

func collectMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics = map[string]uint64{
		"Alloc":            memStats.Alloc,
		"BuckHashSys":      memStats.BuckHashSys,
		"Frees":            memStats.Frees,
		"GCCPUFraction":    uint64(math.Round(memStats.GCCPUFraction)),
		"GCSys":            memStats.GCSys,
		"HeapAlloc":        memStats.HeapAlloc,
		"HeapIdle":         memStats.HeapIdle,
		"HeapInuse":        memStats.HeapInuse,
		"HeapObjects":      memStats.HeapObjects,
		"HeapReleased":     memStats.HeapReleased,
		"HeapSys":          memStats.HeapSys,
		"LastGC":           memStats.LastGC,
		"Lookups":          memStats.Lookups,
		"MCacheInuse":      memStats.MCacheInuse,
		"MCacheSys":        memStats.MCacheSys,
		"MSpanInuse":       memStats.MSpanInuse,
		"MSpanSys":         memStats.MSpanSys,
		"Mallocs":          memStats.Mallocs,
		"NextGC":           memStats.NextGC,
		"NumForcedGC":      uint64(memStats.NumForcedGC),
		"NumGC":            uint64(memStats.NumGC),
		"OtherSys":         memStats.OtherSys,
		"PauseTotalNs":     memStats.PauseTotalNs,
		"StackInuse":       memStats.StackInuse,
		"StackSys":         memStats.StackSys,
		"Sys":              memStats.Sys,
		"TotalAlloc":       memStats.TotalAlloc,
	}
	PollCount   += 1
    RandomValue = uint64(time.Now().Nanosecond())
}

func sendMetric(metricType string, metricName string, metricValue uint64) {
	url := fmt.Sprintf("%s%s/%s/%d", serverAddress, metricType, metricName, metricValue)
	response, err := http.Post(url, "text/plain", bytes.NewBuffer(nil))
	if err != nil {
		fmt.Printf("Failed to send metric: %v\n", err)
		return
	}
	defer response.Body.Close()
}


func sendMetrics() {
	for metricName, metricValue := range metrics {
		sendMetric("gauge", metricName, metricValue)
	}
	sendMetric("counter", "PollCount", PollCount)
	sendMetric("gauge", "RandomValue", RandomValue)
}
