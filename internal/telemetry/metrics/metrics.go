package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var meter = otel.Meter("github.com/loczek/go-link-shortener")

var (
	CacheRequestsCounter metric.Int64Counter
	CollisionsCounter    metric.Int64Counter
	HttpRequestsCounter  metric.Int64Counter
	ReportCounter        metric.Int64Counter
)

func init() {
	CacheRequestsCounter = createIntCounter(
		"cache.requests",
		metric.WithDescription("Number of cache requests."),
		metric.WithUnit("{count}"),
	)
	CollisionsCounter = createIntCounter(
		"collisions",
		metric.WithDescription("Number of hash collisions."),
		metric.WithUnit("{count}"),
	)
	HttpRequestsCounter = createIntCounter(
		"http.requests",
		metric.WithDescription("Number of api hits"),
		metric.WithUnit("{count}"),
	)
	ReportCounter = createIntCounter(
		"reports",
		metric.WithDescription("Number of reports made"),
		metric.WithUnit("{count}"),
	)
}

func createIntCounter(name string, options ...metric.Int64CounterOption) metric.Int64Counter {
	counter, err := meter.Int64Counter(
		name, options...,
	)
	if err != nil {
		panic(err)
	}
	return counter
}
