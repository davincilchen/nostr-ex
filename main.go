package main

import (
	"nostr-ex/pkg/app/server"
	"nostr-ex/pkg/config"
	//_ "github.com/go-sql-driver/mysql"
	// "go.opentelemetry.io/otel/exporters/prometheus"
	// "go.opentelemetry.io/otel/metric/global"
	// controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	// //"go.opentelemetry.io/otel/sdk/metric/export/aggregation"
	// //processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	// //selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	// "go.opentelemetry.io/otel/sdk/resource"
	// semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

const confPath = "./config.json"

//var serviceName = semconv.ServiceNameKey.String("nostr-ex")

// func initMeter() {
// 	c := controller.New(
// 		processor.NewFactory(
// 			selector.NewWithHistogramDistribution(),
// 			aggregation.CumulativeTemporalitySelector(),
// 			processor.WithMemory(true),
// 		),
// 		controller.WithResource(resource.NewWithAttributes(
// 			semconv.SchemaURL,
// 			serviceName,
// 		)),
// 	)
// 	metricExporter, err := prometheus.New(prometheus.Config{}, c)
// 	if err != nil {
// 		log.Fatalf("failed to install metric exporter, %v", err)
// 	}
// 	global.SetMeterProvider(metricExporter.MeterProvider())

// 	http.HandleFunc("/", metricExporter.ServeHTTP)
// 	go func() {
// 		_ = http.ListenAndServe(":2222", nil)
// 	}()
// 	fmt.Println("Prometheus server running on :2222")
// }

func main() {

	//initTracer()
	//initMeter()

	//cfg, err := config.New(confPath)
	cfg, _ := config.New(confPath)

	svr := server.New(cfg)
	svr.Serve()

}
