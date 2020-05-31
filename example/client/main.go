package main

import (
    "context"
    "log"
    "fmt"
    "time"
    "net/http"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"

    "go.opentelemetry.io/otel/api/global"
    "go.opentelemetry.io/otel/api/kv"

    "go.opentelemetry.io/otel/exporters/trace/jaeger"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    // "go.opentelemetry.io/contrib/plugins/go.mongodb.org/mongo-driver"
  )


func initTracer() func() {
	// Create and install Jaeger export pipeline
	_, flush, err := jaeger.NewExportPipeline(
		jaeger.WithCollectorEndpoint("http://jaeger:14268/api/traces"),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: "trace-demo",
			Tags: []kv.KeyValue{
				kv.String("exporter", "jaeger"),
				kv.Float64("float", 312.23),
			},
		}),
		jaeger.RegisterAsGlobal(),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
	)
	if err != nil {
		log.Fatal(err)
	}

	return func() {
		flush()
	}
}

func main() {
  fn := initTracer()
  defer fn()

  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://database:5000"))
  if err != nil {
    log.Fatal(err)
  }
  _ = client.Database("testing").Collection("animals")
  fmt.Println("):")

  tracer := global.Tracer("component-main")
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    _, span := tracer.Start(ctx, "/")
    defer span.End()
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
  })
  http.ListenAndServe(":8080", nil)
}
