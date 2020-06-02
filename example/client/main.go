package main

import (
    "context"
    "log"
    "fmt"
    "time"
    "net/http"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"

    "go.opentelemetry.io/otel/api/global"
    "go.opentelemetry.io/otel/api/kv"

    "go.opentelemetry.io/otel/exporters/trace/jaeger"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    otelmongo "go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver"
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

func insertAnimal(name string) {
  tracer := global.Tracer("component-main")

  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  opts := options.Client()
  opts.Monitor = otelmongo.NewMonitor("mongo", otelmongo.WithTracer(tracer))
  opts.ApplyURI("mongodb://tiger:meow@database:27017")
  client, err := mongo.Connect(ctx, opts)
  if err != nil {
    log.Fatal(err)
  }
  err = client.Ping(ctx, readpref.Primary())
  if err != nil {
    fmt.Println("failed to connect")
    log.Fatal(err)
  }

  collection := client.Database("zoo").Collection("animals")
  _, err = collection.InsertOne(ctx, bson.D{{Key: "animal", Value: name}})
  if err != nil {
    fmt.Println("failed to insert")
    log.Fatal(err)
  }
}

func main() {
  fn := initTracer()
  defer fn()

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    name := r.URL.Path[1:]
    insertAnimal(name)
    fmt.Fprintf(w, "Inserting %s", name)
  })
  http.ListenAndServe(":8080", nil)
}
