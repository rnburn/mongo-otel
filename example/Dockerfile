FROM ubuntu:18.04

RUN apt-get update \
  && apt-get install --no-install-recommends --no-install-suggests -y \
       golang \
       ca-certificates \
       git \
  && go get go.mongodb.org/mongo-driver/mongo
