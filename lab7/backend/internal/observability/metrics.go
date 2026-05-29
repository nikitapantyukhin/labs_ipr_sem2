package observability

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sport_platform",
			Subsystem: "backend",
			Name:      "http_requests_total",
			Help:      "Total number of backend HTTP requests grouped by method, route template and status code.",
		},
		[]string{"method", "route", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "sport_platform",
			Subsystem: "backend",
			Name:      "http_request_duration_seconds",
			Help:      "Backend HTTP request duration in seconds.",
			Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
		[]string{"method", "route", "status"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "sport_platform",
			Subsystem: "backend",
			Name:      "http_requests_in_flight",
			Help:      "Current number of backend HTTP requests being processed.",
		},
	)

	businessEventsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sport_platform",
			Subsystem: "backend",
			Name:      "business_events_total",
			Help:      "Business events produced by the sport platform backend.",
		},
		[]string{"event"},
	)
)

func MetricsHandler() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

func PrometheusMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/metrics" {
			ctx.Next()
			return
		}

		startedAt := time.Now()
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		ctx.Next()

		route := ctx.FullPath()
		if route == "" {
			route = "unmatched"
		}
		status := strconv.Itoa(ctx.Writer.Status())

		httpRequestsTotal.WithLabelValues(ctx.Request.Method, route, status).Inc()
		httpRequestDuration.WithLabelValues(ctx.Request.Method, route, status).Observe(time.Since(startedAt).Seconds())
	}
}

func RecordBusinessEvent(event string) {
	businessEventsTotal.WithLabelValues(event).Inc()
}

func init() {
	businessEventsTotal.WithLabelValues("user_registered").Add(0)
	businessEventsTotal.WithLabelValues("club_created").Add(0)
	businessEventsTotal.WithLabelValues("join_request_created").Add(0)
	businessEventsTotal.WithLabelValues("workout_created").Add(0)
}
