package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var promHttpCalls = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "sg_http_calls",
		Help: "A counter for the total number of http calls",
	},
	[]string{"handler", "code"},
)

// duration is partitioned by the HTTP method and handler. It uses custom
// buckets based on the expected request duration.
var promHttpLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "sg_http_latency",
		Help:    "A histogram of latencies for requests.",
		Buckets: []float64{.25, .5, 1, 2.5, 5, 10},
	},
	[]string{"handler", "code"},
)

// responseSize has no labels, making it a zero-dimensional ObserverVec.
var promResponseSize = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "sg_response_size_bytes",
		Help:    "A histogram of response sizes for requests.",
		Buckets: []float64{200, 500, 900, 1500},
	},
	[]string{},
)

var allMetrics []prometheus.Collector = []prometheus.Collector{
	promHttpCalls, promHttpLatency, promResponseSize,
}
