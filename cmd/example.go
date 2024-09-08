package main

import (
	"context"
	"net/http"
	"otel-library/errs"
	_logging "otel-library/logs"
	"otel-library/otelBuilder"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	metric_sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

const (
	endpointURl = "http://k8s-eck-eckstack-1b5d1a8c6a-a5d028e65b82cd1b.elb.us-east-2.amazonaws.com:8200"
)

var otelConfig = initOtel()

func initOtel() *otelBuilder.Otel {
	ctx := context.Background()
	l := _logging.NewOtelLogging()

	batchOpts := []trace.BatchSpanProcessorOption{
		trace.WithBatchTimeout(time.Second * 10),
	}

	metricOps := []metric_sdk.PeriodicReaderOption{
		metric_sdk.WithInterval(time.Second * 30),
		metric_sdk.WithTimeout(time.Second * 60),
	}

	tracing, metrics, err := otelBuilder.NewOtelBuilder().
		// WithConsoleExporter().
		WithAuthHeader("WdXC96655L2E8xU29yHx1ZXv").
		WithEndpointURL(endpointURl).
		WithInsecure(true).
		WithServiceName("testing-api").
		WithTraceBatchSpanProcessorOption(batchOpts...).
		WithMetricPeriodicReaderOption(metricOps...).
		Build(ctx, l)

	if err != nil {
		l.Error("Failed to initialize OpenTelemetry", err)
		return nil
	}

	return otelBuilder.NewOtel(tracing, metrics, l)
}

func main() {

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		_, span := otelConfig.Tracing.GetTracer().Start(r.Context(), r.URL.Path)
		defer otelConfig.Tracing.EndSpan(span)

		counter, _ := otelConfig.Metrics.CreateCounter(r.URL.Path,
			metric.WithDescription("Total number calling this api"),
			metric.WithUnit("request"),
		)

		counter.Add(r.Context(), 1)

		// testFunc(ctx)
		w.Write([]byte("Hello World!"))
	})

	http.HandleFunc("/test2", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otelConfig.Tracing.GetTracer().Start(r.Context(), r.URL.Path)
		defer otelConfig.Tracing.EndSpan(span)

		counter, _ := otelConfig.Metrics.CreateCounter(r.URL.Path,
			metric.WithDescription("Total number calling this api2"),
			metric.WithUnit("request"),
		)

		counter.Add(r.Context(), 1)
		testFunc(ctx)
		w.Write([]byte("Hello World!"))
	})

	http.ListenAndServe(":8080", nil)
}

func testFunc(ctx context.Context) {
	ctx, span := otelConfig.Tracing.GetTracer().Start(ctx, "Test Func 1")
	defer otelConfig.Tracing.EndSpan(span)

	otelConfig.Tracing.AddAttributes(span, nil, attribute.String("testFunc", "2"))

	time.Sleep(time.Second * 3)
	testFunc2(ctx)
}

func testFunc2(ctx context.Context) {
	ctx, span := otelConfig.Tracing.GetTracer().Start(ctx, "Test Func 2")
	otelConfig.Tracing.RecordError(ctx, span, errs.CreateError(errs.BAD_REQUEST, errs.BadRequest, "Something Went Wrong"))
	defer otelConfig.Tracing.EndSpan(span)
}
