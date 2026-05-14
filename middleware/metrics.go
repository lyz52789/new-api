package middleware

import (
	"time"

	appmetrics "github.com/QuantumNous/new-api/metrics"
	"github.com/gin-gonic/gin"
)

func PrometheusMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		tag := getMetricsRouteTag(c)
		appmetrics.IncHTTPInFlight(tag)
		start := time.Now()

		defer func() {
			tag = getMetricsRouteTag(c)
			appmetrics.DecHTTPInFlight(tag)
			appmetrics.ObserveHTTPRequest(c.Request.Method, c.FullPath(), c.Writer.Status(), tag, time.Since(start))
		}()

		c.Next()
	}
}

func getMetricsRouteTag(c *gin.Context) string {
	tag, _ := c.Get(RouteTagKey)
	if tagString, ok := tag.(string); ok && tagString != "" {
		return tagString
	}
	return "web"
}
