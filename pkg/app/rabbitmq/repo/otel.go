package repo

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func NewMetrics(meterName string) *Metrics {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	meter := provider.Meter(meterName)

	s := ""
	s2 := ""
	s = fmt.Sprintf("%s success_counter", meterName)
	s2 = fmt.Sprintf("%s: how many process success", meterName)
	success, err := meter.Int64Counter(s, metric.WithDescription(s2))
	if err != nil {
		log.Fatal(err) //TODO:
	}
	s = fmt.Sprintf("%s fail_counter", meterName)
	s2 = fmt.Sprintf("%s: how many process fail", meterName)
	fail, err := meter.Int64Counter(s, metric.WithDescription(s2))
	if err != nil {
		log.Fatal(err) //TODO:
	}

	s = fmt.Sprintf("%s duration_in_milliseconds", meterName)
	s2 = fmt.Sprintf("%s: duration of process", meterName)
	duration, err := meter.Int64Histogram(s, metric.WithDescription(s2))
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
