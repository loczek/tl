# Telemetry

## Metrics, Traces and Logs

- Metrics are collected with otlp and exported with the [otlp prometheus exporter](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/prometheus) and served with the [prometheus client library](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/prometheus) for prometheus to be pulled
- Logs are collected with otlp and exported with [otlp http exporter](go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp) and pushed to loki
- Traces are collected with otlp and exported with [otlp http exporter](http://go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp) and pushed to tempo

## Transport Method

- metrics and logs are transported via http with protobuf payloads

# Notes

## Push vs Pull

- push doesn't require service discovery
- pull has health checks automatically
- pull is the standard


## Links

- [prometheus](https://prometheus.io/docs/prometheus/latest/installation)
- [otlp prometheus exporter example](https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/examples/prometheus/main.go)
- [otlp prometheus exporter go package](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/prometheus)
- [otlp go](https://opentelemetry.io/docs/languages/go/instrumentation)
- [otlp exporters](https://opentelemetry.io/docs/languages/go/exporters/#prometheus-experimental)
- [redis go](https://redis.uptrace.dev/guide/go-redis.html)
- [7tv](https://github.com/SevenTV/API)
- [fiber prometheus (not used)](https://github.com/ansrivas/fiberprometheus)
- [fiber prometheus grafana dashboard (not used)](https://grafana.com/grafana/dashboards/14331-fiber-framework-processes/)
- [loki logs](https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/exporter/lokiexporter/README.md)
- [prometheus example?](https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/examples/prometheus/main.go)
