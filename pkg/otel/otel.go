package otel

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel/exporters/prometheus"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var serviceName = semconv.ServiceNameKey.String("otel-nostr")

var meter metric.Meter

func serveMetrics() {
	log.Printf("serving metrics at localhost:2223/metrics")
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2223", nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}

func NewTracerProvider() trace.TracerProvider {

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			serviceName,
		)),
	)

	otel.SetTracerProvider(tp)

	newMeter("nostr-ex")
	go serveMetrics()

	return otel.GetTracerProvider()
}

func newMeter(meterName string) {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	m := provider.Meter(meterName)
	meter = m
}

func GetMeter() metric.Meter {
	return meter
}
