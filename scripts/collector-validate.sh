docker run --rm \
    -v "$PWD/collector/config.yaml:/etc/otelcol-contrib/config.yaml:ro" \
    otel/opentelemetry-collector-contrib \
    validate --config=/etc/otelcol-contrib/config.yaml
