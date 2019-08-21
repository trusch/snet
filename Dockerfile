FROM golang:1.12 as builder

ENV GOPATH=/go
ENV GO111MODULE=on
COPY go.mod go.sum /go/src/github.com/trusch/snet/
WORKDIR /go/src/github.com/trusch/snet/
RUN go mod download
COPY cmd pkg .
RUN go install -v github.com/trusch/snet/...
