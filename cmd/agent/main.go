package main

import (
	"bytes"
	"fmt"
	"strconv"
	"net/http"
	"runtime"
	"time"
	"math"
	"flag"
	"os"
)

// PollCount is the number of pollings done
var (
	metrics        map[string]uint64
	serverAddress  string
	pollInterval   time.Duration
	reportInterval time.Duration
	PollCount      uint64 = 0
	RandomValue    uint64 = 0
)

func main() {
	var (
		err error
		serverAddressFlag  string
		pollIntervalFlag   int
		reportIntervalFlag int
	)
	flag.StringVar(&serverAddressFlag,  "a", "localhost:8080", "The value for the -a flag")
	flag.IntVar(&pollIntervalFlag,   "p",  2, "The value for the -p flag")
	flag.IntVar(&reportIntervalFlag, "r", 10, "The value for the -r flag")
	flag.Parse()


	environmetAddress, exists := os.LookupEnv("ADDRESS")
    if exists { serverAddressFlag = environmetAddress }

	environmetPollInterval, exists := os.LookupEnv("POLL_INTERVAL")
    if exists {
		if pollIntervalFlag, err = strconv.Atoi(environmetPollInterval); err != nil {
			panic(err)
		}
	}

	environmetReportInterval, exists := os.LookupEnv("REPORT_INTERVAL")
    if exists {
		if reportIntervalFlag, err = strconv.Atoi(environmetReportInterval); err != nil{
            panic(err)
        }
    }

fmt.Printf("Running with following parameters:\n" +
	"Server Name:      %s\n" +
	"Poll Interval:    %d\n" +
	"Report Interval:  %d\n",
	serverAddressFlag, pollIntervalFlag, reportIntervalFlag)

	serverAddress  = fmt.Sprintf("http://%s/update/", serverAddressFlag)
	pollInterval   = time.Duration(pollIntervalFlag) * time.Second
	reportInterval = time.Duration(reportIntervalFlag) * time.Second


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
	PollCount++
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
