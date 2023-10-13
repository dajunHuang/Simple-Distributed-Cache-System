#syntax=docker/dockerfile:1
FROM ubuntu:20.04

WORKDIR /root
COPY . .

RUN apt update && apt install -y protobuf-compiler wget
RUN wget https://go.dev/dl/go1.21.3.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz

ENV GOPATH="/root/go"
ENV PATH="/usr/local/go/bin:/root/go/bin:$PATH"

RUN go env -w GOPROXY=https://goproxy.cn
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

RUN protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative cache/cache.proto
RUN go build server.go client.go
