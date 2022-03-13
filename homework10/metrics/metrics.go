package metrics

import {
	"io"
	"os"
	"strings"
	"time"
	
	"github.com/prometheus/client_golang/prometheus/promhttp"
}

func Register() {
  err := prometheus.Register(functionLatency)
  if err != nil {
    fmt.Println(err)
  }
}

const (
  MetricsNamespace = "default"
)

func NewTimer() *ExecutionTimer {
  return NewExecutionTimer(functionLatency)
}

var (
  functionLatency = CreateExecutionTimeMetric(MetricsNamespace, "Time spent")
)

funct NewExecutionTimer(histo *prometheus.HistogramVec) *ExecutionTimer {
  now := time.Now()
  return &ExecutionTimer {
    histo: histo,
    start: now,
    last:  now,
  }
}

func (t *ExecutionTimer) ObseerveTotal() {
  (*t.histo).WithLabelValues("total").Observe(time.Now().Sub(t.start).Seconds())
}

func CreateExecutionTimeMetric(namespace string, help string) *prometheus.HistogramVec {
  return prometheus.NewHistogramVec(
    prometheus.HistogramOpts{
      Namespace: namespace,
      Name:      "execution_latency_seconds",
      Help:      help,
      Buckets:   prometheus.ExponentialBuckets(0.001, 2, 15),
    }, []string{"step"}
  )
}
