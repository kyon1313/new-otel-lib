package middleware

import (
	apw_metrics "otel-library/metrics"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/metric"
)

// Gin middleware to handle metrics collection for each endpoint.
// The middleware increments a counter each time an endpoint is hit, ensuring that metrics are consistently collected across all endpoints.
func MetricsMiddleware(metrics apw_metrics.OtelMetric) gin.HandlerFunc {
	return func(c *gin.Context) {
		counter, _ := metrics.CreateCounter(c.Request.URL.Path,
			metric.WithDescription("Total number calling this api"),
			metric.WithUnit("request"),
		)

		counter.Add(c.Request.Context(), 1)
		c.Next()
	}
}
