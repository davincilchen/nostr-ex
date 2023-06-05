package repo

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func NewConsumerMetrics(meterName string) *Metrics {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	meter := provider.Meter(meterName)

	success, err := meter.Int64Counter("success_counter", metric.WithDescription("how many process success"))
	if err != nil {
		log.Fatal(err) //TODO:
	}

	fail, err := meter.Int64Counter("fail_counter", metric.WithDescription("how many process success"))
	if err != nil {
		log.Fatal(err) //TODO:
	}

	return &Metrics{
		success: success,
		fail:    fail,
	}

}

type Metrics struct {
	success metric.Int64Counter
	fail    metric.Int64Counter
}

func (t *Metrics) Success(ctx context.Context) {
	t.success.Add(ctx, 1)
}

func (t *Metrics) Fail(ctx context.Context) {
	t.fail.Add(ctx, 1)
}
