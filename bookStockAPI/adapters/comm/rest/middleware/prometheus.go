package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

//RequestCounterVec counts total request per endpoint and method
var (
	RequestCounterVec = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "http",
			Subsystem: "requests",
			Name:      "number_of_requests",
			Help:      "Total number of requests handled by the API",
		},
		[]string{"endpoint", "method"},
	)
)

//RequestDurationGauge calculates avg time per endpoint and method
var (
	RequestDurationGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "http",
			Subsystem: "requests",
			Name:      "request_duration",
			Help:      "Avg time of requests handled by the API",
		},
		[]string{"endpoint", "method"},
	)
)

//MetricsMiddleware writes increments request count and calculates request duration and writes to Prometheus metrics
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := strings.TrimRight(r.URL.Path, "/")
		ignoreRoutes := make([]string, 0)
		ignoreRoutes = append(ignoreRoutes, "/metrics", "/swagger")
		if !contains(ignoreRoutes, url) {
			startTime := time.Now()
			log.Println(r.RequestURI)
			RequestCounterVec.WithLabelValues(url, r.Method).Inc()
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
			duration := time.Now().Sub(startTime).Seconds() * 1000 //Duration in milliseconds
			RequestDurationGauge.WithLabelValues(url, r.Method).Set(duration)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
