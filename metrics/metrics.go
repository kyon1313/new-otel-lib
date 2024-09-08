package apw_metrics

import (
	"go.opentelemetry.io/otel/metric"
)

// OtelMetric defines methods for metrics operations.
type OtelMetric interface {
	CreateCounter(name string, opt ...metric.Int64CounterOption) (metric.Int64Counter, error)
	CreateUpDownCounter(name string, opt ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error)
	CreateHistogram(name string, opt ...metric.Int64HistogramOption) (metric.Int64Histogram, error)
	CreateGauge(name string, opt ...metric.Int64GaugeOption) (metric.Int64Gauge, error)
}

// metricImpl is a concrete implementation of OtelMetric.
type metricImpl struct {
	Meter metric.Meter
}

// NewMetric creates a new metric instance.
func NewMetric(meter metric.Meter) OtelMetric {
	return &metricImpl{Meter: meter}
}

func (m *metricImpl) CreateCounter(name string, opt ...metric.Int64CounterOption) (metric.Int64Counter, error) {
	return m.Meter.Int64Counter(name, opt...)
}

func (m *metricImpl) CreateUpDownCounter(name string, opt ...metric.Int64UpDownCounterOption) (metric.Int64UpDownCounter, error) {
	return m.Meter.Int64UpDownCounter(name, opt...)
}

func (m *metricImpl) CreateHistogram(name string, opt ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
	return m.Meter.Int64Histogram(name, opt...)
}

func (m *metricImpl) CreateGauge(name string, opt ...metric.Int64GaugeOption) (metric.Int64Gauge, error) {
	return m.Meter.Int64Gauge(name, opt...)
}
