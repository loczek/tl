package metrics

import (
	"context"
	"log/slog"
	"time"

	"github.com/loczek/go-link-shortener/internal/config"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/log/global"
	skdlog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func NewResource() (*sdkresource.Resource, error) {
	return sdkresource.Merge(sdkresource.Default(),
		sdkresource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("go-link-shortener"),
			semconv.ServiceVersion("0.1.0"),
		))
}

func NewMeterProviderPrometheus(res *sdkresource.Resource) (*sdkmetric.MeterProvider, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(exporter),
	)

	otel.SetMeterProvider(meterProvider)

	return meterProvider, nil
}

func NewMeterProviderStdout(res *sdkresource.Resource) (*sdkmetric.MeterProvider, error) {
	exporter, err := stdoutmetric.New()
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter,
			sdkmetric.WithInterval(time.Second*5))),
	)

	otel.SetMeterProvider(meterProvider)

	return meterProvider, nil
}

func NewMeterProviderHttp(res *sdkresource.Resource) (*sdkmetric.MeterProvider, error) {
	exporter, err := otlpmetrichttp.New(context.Background())
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter,
			sdkmetric.WithInterval(time.Second*15))),
	)

	otel.SetMeterProvider(meterProvider)

	return meterProvider, nil
}

func NewLoggerProvider(res *sdkresource.Resource) (*skdlog.LoggerProvider, error) {
	exporter, err := otlploghttp.New(context.Background())
	if err != nil {
		return nil, err
	}

	loggerProvider := skdlog.NewLoggerProvider(
		skdlog.WithProcessor(
			skdlog.NewBatchProcessor(exporter, skdlog.WithExportInterval(time.Second*15)),
		),
		skdlog.WithResource(res),
	)

	if config.IsProd() || config.LOG_TO_STDOUT {
		global.SetLoggerProvider(loggerProvider)

		logger := otelslog.NewLogger("go-link-shortener", otelslog.WithLoggerProvider(loggerProvider))

		slog.SetDefault(logger)
	}

	return loggerProvider, nil
}
