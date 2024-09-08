package otelBuilder

import (
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/metric"
	oteltrace "go.opentelemetry.io/otel/sdk/trace"
)

// Console Exporter, only for testing
func newConsoleMetricExporter() (metric.Exporter, error) {
	return stdoutmetric.New(stdoutmetric.WithPrettyPrint())
}

func newConsoleTraceExporter() (oteltrace.SpanExporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}
