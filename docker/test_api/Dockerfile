# DEVELOPMENT ONLY, DO NOT USE FOR PRODUCTION
FROM golang:1.19-alpine
RUN mkdir /badaas
WORKDIR /badaas
COPY go.mod .
COPY go.sum .
RUN go mod download
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOFLAGS=-buildvcs=false