package repo

import (
	"context"
	"log"
	"time"

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

	fail, err := meter.Int64Counter("fail_counter", metric.WithDescription("how many process fail"))
	if err != nil {
		log.Fatal(err) //TODO:
	}

	duration, err := meter.Int64Histogram("duration_in_milliseconds", metric.WithDescription("duration of process"))
	if err != nil {
		log.Fatal(err) //TODO:
	}

	return &Metrics{
		success:  success,
		fail:     fail,
		duration: duration,
	}

}

type Metrics struct {
	success  metric.Int64Counter
	fail     metric.Int64Counter
	duration metric.Int64Histogram
}

func (t *Metrics) Success(ctx context.Context) {
	t.success.Add(ctx, 1)
}

func (t *Metrics) Fail(ctx context.Context) {
	t.fail.Add(ctx, 1)
}

func (t *Metrics) Duration(ctx context.Context, tt time.Time) {
	t.duration.Record(ctx, time.Since(tt).Milliseconds())
}
