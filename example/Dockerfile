FROM ubuntu:18.04

ADD client /client

RUN apt-get update \
  && apt-get install --no-install-recommends --no-install-suggests -y \
       golang \
       ca-certificates \
       git \
  && go get go.mongodb.org/mongo-driver/mongo \
  && cd /client \
  && go build
