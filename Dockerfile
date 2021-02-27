FROM golang:1.15 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on

WORKDIR /go/src/github.com/Le0tk0k/blog-server

COPY . .
RUN go mod download
RUN go build .

FROM alpine:3.13.0

ENV DOCKERIZE_VERSION v0.6.1

RUN apk add --no-cache bash ca-certificates curl openssl tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

COPY ./db/migrations ./db/migrations
COPY --from=build /go/src/github.com/Le0tk0k/blog-server/blog-server /

RUN chmod a+x blog-server
CMD ./blog-server