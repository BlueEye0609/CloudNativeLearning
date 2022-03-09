package metrics

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricsNamespace = "httpserver"
)

// var totalRequests = prometheus.NewCounterVec(
// 	prometheus.CounterOpts{
// 		Name: "http_requests_total",
// 		Help: "Number of get requests.",
// 	},
// 	[]string{"path"},
// )

// var responseStatus = prometheus.NewCounterVec(
// 	prometheus.CounterOpts{
// 		Name: "response_status",
// 		Help: "Status of HTTP response",
// 	},
// 	[]string{"status"},
// )

// var httpDuration = prometheus.NewHistogramVec(
// 	prometheus.HistogramOpts{
// 		namespace: MetricsNamespace,
// 		Name:      "execution_latency_seconds",
// 		Help:      "Time Spend",
// 		Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
// 	},
// 	[]string{"step"},
// )

func Register() {
	// prometheus.Register(totalRequests)
	// prometheus.Register(responseStatus)
	// prometheus.Register(httpDuration)
	err := prometheus.Register(functionLatency)
	if err != nil {
		fmt.Println(err)
	}
}

var (
	functionLatency = CreateExecutionTimerMetric(MetricsNamespace, "Time spent.")
)

// newExecutionTimer provides a timer for updater's RunOnce execution
func NewTimer() *ExecutionTimer {
	return NewExecutionTimer(functionLatency)
}

// NewExecutionTimer provides a timer for admission latency; call ObserveXXX() on it to measure
func NewExecutionTimer(histo *prometheus.HistogramVec) *ExecutionTimer {
	now := time.Now()
	return &ExecutionTimer{
		histo: histo,
		start: now,
		last:  now,
	}
}

// Oberseve Total measures the execution time from the creation of the Execution Timer
func (t *ExecutionTimer) ObserveTotal() {
	(*t.histo).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}

// CreateExecutionTimeMetrics prepares a new histogram labled with execution step
func CreateExecutionTimerMetric(namespace string, help string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"},
	)
}

// ExecutionTimer measures execution time of a computation, split into major steps
type ExecutionTimer struct {
	histo *prometheus.HistogramVec
	start time.Time
	last  time.Time
}
