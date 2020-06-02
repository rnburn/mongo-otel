FROM ubuntu:18.04

ADD client /client

RUN apt-get update \
  && apt-get install --no-install-recommends --no-install-suggests -y \
       golang \
       ca-certificates \
       git \
  && go get go.opentelemetry.io/otel \
  && go get github.com/apache/thrift/lib/go/thrift \
  && go get google.golang.org/api/support/bundler \
  && go get go.mongodb.org/mongo-driver/mongo \
  && go get go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver \
  && cd /client \
  && go build
