package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"github.com/gorilla/mux"
	"sort"
	"flag"
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

func (ms MemStorage) GetValueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	metricType := GetMetricTypeFromString(vars["metricType"])
	metricName := vars["metricName"]
	if metric, ok := ms[metricName]; ok {
		if metric.Type != metricType {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		value := fmt.Sprintf("%v", metric.Value)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(value))
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (ms MemStorage) GetAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<html><body><h1>Metrics</h1>")
	var metricNames []string
	for metricName := range ms {
		metricNames = append(metricNames, metricName)
	}

	sort.Strings(metricNames)

	for _, metricName := range metricNames {
		metric := ms[metricName]
		fmt.Fprintf(w, "<p>%s: %v</p>", metricName, metric.Value)
	}
	fmt.Fprintf(w, "</body></html>")
}


func main() {
	var serverAddress string
	flag.StringVar(&serverAddress, "a", "localhost:8080", "The value for the -a flag")
	flag.Parse()

	storage := make(MemStorage)
	router := mux.NewRouter()
	router.HandleFunc("/value/{metricType}/{metricName}", storage.GetValueHandler).Methods("GET")
	router.HandleFunc("/update/{metricType}/{metricName}/{metriccValue}", storage.MetricsHandler).Methods("POST")
	router.HandleFunc("/", storage.GetAllMetricsHandler).Methods("GET")

	fmt.Printf("Server listening on http://%s", serverAddress)
	listenPort := fmt.Sprintf(":%s", strings.Split(serverAddress, ":")[1])
	err := http.ListenAndServe(listenPort, router)
	if err != nil {
		panic(err)
	}
}
