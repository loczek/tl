package metrics

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func NewResource() (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("go-link-shortener"),
			semconv.ServiceVersion("0.1.0"),
		))
}

func NewMeterProviderPrometheus(res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	exp, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exp),
	)

	return meterProvider, nil
}

func NewMeterProviderStdout(res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	exp, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp,
			sdkmetric.WithInterval(time.Second*5))),
	)

	return meterProvider, nil
}

func NewMeterProviderHttp(res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	exp, err := otlpmetrichttp.New(context.Background())
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp,
			sdkmetric.WithInterval(time.Second*5))),
	)

	return meterProvider, nil
}
