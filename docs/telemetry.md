# Telemetry

## Metrics, Traces and Logs

- Metrics are collected with otlp and exported with the [otlp prometheus exporter](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/prometheus) and served with the [prometheus client library](https://pkg.go.dev/go.opentelemetry.io/otel/exporters/prometheus) for prometheus to be pulled
- Logs are collected with otlp and exported with [otlp http exporter](go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp) and pushed to loki
- Traces are not collected

## Transport Method

- metrics and logs are transported via http with protobuf payloads
