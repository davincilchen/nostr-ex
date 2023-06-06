package repo

import (
	"context"
	"fmt"
	"log"
	"nostr-ex/pkg/otel"
	"time"

	"go.opentelemetry.io/otel/metric"
	//"go.opentelemetry.io/otel/metric"
)

var shareMetrics *ShareMetrics

func NewMetrics(meterName string) *Metrics {
	// exporter, err := prometheus.New()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	// meter := provider.Meter(meterName)
	meter := otel.GetMeter()
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

type ShareMetrics struct {
	queueSize metric.Int64UpDownCounter
}

func GetShareMetrics() *ShareMetrics {
	if shareMetrics != nil {
		return shareMetrics
	}

	shareMetrics = newShareMetrics()
	return shareMetrics
}

func newShareMetrics() *ShareMetrics {
	meterName := "MQ"
	meter := otel.GetMeter()
	s := ""
	s2 := ""
	s = fmt.Sprintf("%s queue_size", meterName)
	s2 = fmt.Sprintf("%s: how many massage still in queue", meterName)
	queueSize, err := meter.Int64UpDownCounter(s, metric.WithDescription(s2))
	if err != nil {
		log.Fatal(err) //TODO:
	}

	return &ShareMetrics{
		queueSize: queueSize,
	}

}

func (t *ShareMetrics) Enqueue(ctx context.Context) {
	t.queueSize.Add(ctx, 1)
}

func (t *ShareMetrics) Dequeue(ctx context.Context) {
	t.queueSize.Add(ctx, -1)
}
