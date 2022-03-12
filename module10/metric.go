package module10

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	//HTTPReqDuration metric:http_request_duration_seconds
	HTTPReqDuration *prometheus.HistogramVec
	//HTTPReqTotal metric:http_request_total
	HTTPReqTotal *prometheus.CounterVec
)

func init() {

	// request duration
	HTTPReqDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "The HTTP request latencies in seconds.",
		Buckets: nil,
	}, []string{"method", "path"})

	// request total count
	HTTPReqTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests made.",
	}, []string{"method", "path", "status"})

	prometheus.MustRegister(
		HTTPReqDuration,
		HTTPReqTotal,
	)
}