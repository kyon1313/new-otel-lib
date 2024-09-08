package otelBuilder

import (
	apw_logging "otel-library/logs"
	apw_metrics "otel-library/metrics"
	apw_tracing "otel-library/tracing"
)

type Otel struct {
	Tracing apw_tracing.OtelTracing
	Metrics apw_metrics.OtelMetric
	Logs    apw_logging.OtelLogging
}

func NewOtel(tracing apw_tracing.OtelTracing, metrics apw_metrics.OtelMetric, logs apw_logging.OtelLogging) *Otel {
	return &Otel{
		Tracing: tracing,
		Metrics: metrics,
		Logs:    logs,
	}
}

func (o *Otel) GetTracing() apw_tracing.OtelTracing {
	return o.Tracing
}

func (o *Otel) GetMetrics() apw_metrics.OtelMetric {
	return o.Metrics
}

func (o *Otel) GetLogs() apw_logging.OtelLogging {
	return o.Logs
}
