package main

import (
    "context"
    "log"
    "fmt"
    "time"
    "net/http"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
  )

func main() {
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://database:5000"))
  if err != nil {
    log.Fatal(err)
  }
  _ = client.Database("testing").Collection("animals")
  fmt.Println("):")

  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
  })
  http.ListenAndServe(":8080", nil)
}
