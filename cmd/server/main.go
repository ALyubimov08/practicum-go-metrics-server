package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Индекс метрики
type MetricType int

const (
	Gauge MetricType = iota // 0
	Counter					// 1
)

// Струкрутра для хранения определнной метрики
type Metric struct {
	Type  MetricType
	Value interface{}
}


type MemStorage map[string]Metric

// Строка -> Метрика
func GetMetricTypeFromString(metricType string) MetricType {
    types := map[string]MetricType{
        "gauge":   Gauge,
        "counter": Counter,
    }
    return types[strings.ToLower(metricType)]
}

// Внесение изменений в саму мапу 
func (ms MemStorage) UpdateMetric(metricName string, metricType MetricType, metricValue interface{}) {

	if metric, ok := ms[metricName]; ok {
		if metric.Type != metricType {
			return 
		}

		switch metricType {
		case Gauge:
			ms[metricName] = Metric{Type: metricType, Value: metricValue}
		case Counter:
			if prevValue, ok := metric.Value.(int64); ok {
				if newValue, ok := metricValue.(int64); ok {
					ms[metricName] = Metric{Type: metricType, Value: prevValue + newValue}
				}
			}
		}
	} else {
		ms[metricName] = Metric{Type: metricType, Value: metricValue}
	}
}

// Парсинг запроса на изменение метрик
func (ms *MemStorage) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим путь из URL
	seg := strings.Split(r.URL.Path, "/")
	if len(seg) != 5 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	metricType := seg[2] 
	metricName := seg[3]
	metricValue := seg[4]

	var value interface{}

	switch metricType {
	case "gauge":
		gaugeValue, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		value = gaugeValue
	case "counter":
		counterValue, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		value = counterValue
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	ms.UpdateMetric(metricName, GetMetricTypeFromString(metricType), value)
	w.WriteHeader(http.StatusOK)
}

func main() {
	storage := make(MemStorage)
	
	http.HandleFunc("/update/", storage.MetricsHandler)

	fmt.Println("Server listening on http://localhost:8080")
	err := http.ListenAndServe(`:8080`, nil)
    if err != nil {
        panic(err)
    }
}
