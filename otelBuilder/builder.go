package otelBuilder

import (
	"context"
	"fmt"
	apw_logging "otel-library/logs"
	apw_metrics "otel-library/metrics"
	apw_tracing "otel-library/tracing"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const (
	tracerName = "default-tracer"
	meterName  = "default-meter"
)

type Header map[string]string

type OtelBuilder struct {
	serviceName        string
	traceOpts          []trace.BatchSpanProcessorOption
	metricOpts         []metric.PeriodicReaderOption
	traceExporterOpts  []otlptracehttp.Option
	metricExporterOpts []otlpmetrichttp.Option
	useConsoleExporter bool
}

func NewOtelBuilder() *OtelBuilder {
	return &OtelBuilder{}
}

// WithInsecure configures whether to use an insecure (non-TLS) connection to the OTLP endpoint.
func (b *OtelBuilder) WithInsecure(insecure bool) *OtelBuilder {
	if insecure {
		b.traceExporterOpts = append(b.traceExporterOpts, otlptracehttp.WithInsecure())
		b.metricExporterOpts = append(b.metricExporterOpts, otlpmetrichttp.WithInsecure())
	}
	return b
}

// WithEndpoint sets the OTLP endpoint URL for both tracing and metrics.
// WithEndpoint sets the OTLP endpoint URL for both tracing and metrics.
func (b *OtelBuilder) WithEndpointURL(otlpEndpoint string) *OtelBuilder {
	if otlpEndpoint != "" {
		b.traceExporterOpts = append(b.traceExporterOpts, otlptracehttp.WithEndpointURL(otlpEndpoint))
		b.metricExporterOpts = append(b.metricExporterOpts, otlpmetrichttp.WithEndpointURL(otlpEndpoint))
	}
	return b
}

// WithHeaders adds HTTP headers to be included in requests to the OTLP endpoint.
func (b *OtelBuilder) WithHeaders(headers Header) *OtelBuilder {
	headerMap := make(map[string]string)
	for key, value := range headers {
		headerMap[key] = value
	}
	b.traceExporterOpts = append(b.traceExporterOpts, otlptracehttp.WithHeaders(headerMap))
	b.metricExporterOpts = append(b.metricExporterOpts, otlpmetrichttp.WithHeaders(headerMap))
	return b
}

// WithAuthHeader adds an authorization header with a Bearer token for authentication.
func (b *OtelBuilder) WithAuthHeader(token string) *OtelBuilder {
	return b.WithHeaders(Header{
		"Authorization": "Bearer " + token,
	})
}

// WithServiceName sets the name of the service that will be reported in tracing and metrics data.
func (b *OtelBuilder) WithServiceName(serviceName string) *OtelBuilder {
	if serviceName != "" {
		b.serviceName = serviceName
	}
	return b
}

// WithTraceBatchSpanProcessorOption allows configuring options for the BatchSpanProcessor.
func (b *OtelBuilder) WithTraceBatchSpanProcessorOption(opts ...trace.BatchSpanProcessorOption) *OtelBuilder {
	b.traceOpts = opts
	return b
}

// WithMetricPeriodicReaderOption allows configuring options for the PeriodicReader.
func (b *OtelBuilder) WithMetricPeriodicReaderOption(opts ...metric.PeriodicReaderOption) *OtelBuilder {
	b.metricOpts = opts
	return b
}

// WithConsoleExporter enables the console exporter for debugging purposes.
func (b *OtelBuilder) WithConsoleExporter() *OtelBuilder {
	b.useConsoleExporter = true
	return b
}

// Build creates and returns OtelTracing and OtelMetrics instances using the configured options.
func (b *OtelBuilder) Build(ctx context.Context, l apw_logging.OtelLogging) (apw_tracing.OtelTracing, apw_metrics.OtelMetric, error) {
	var traceExporter trace.SpanExporter
	var metricExporter metric.Exporter
	var err error

	if b.useConsoleExporter {
		traceExporter, err = newConsoleTraceExporter()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create console trace exporter: %w", err)
		}
		metricExporter, err = newConsoleMetricExporter()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create console metric exporter: %w", err)
		}
	} else {
		traceExporter, err = otlptracehttp.New(ctx, b.traceExporterOpts...)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create OTLP trace exporter with options %v: %w", b.traceExporterOpts, err)
		}
		metricExporter, err = otlpmetrichttp.New(ctx, b.metricExporterOpts...)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create OTLP metric exporter with options %v: %w", b.metricExporterOpts, err)
		}
	}

	resourceOpts := resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(b.serviceName))
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter, b.traceOpts...),
		trace.WithResource(resourceOpts),
	)

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter, b.metricOpts...)),
		metric.WithResource(resourceOpts),
	)

	return apw_tracing.NewTracing(tracerProvider.Tracer(b.serviceName), l), apw_metrics.NewMetric(meterProvider.Meter(b.serviceName)), nil
}
