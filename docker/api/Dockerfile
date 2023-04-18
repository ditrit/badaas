# builder image
FROM golang:1.19-alpine as builder
WORKDIR /app
RUN apk add build-base
COPY commands/ ./commands
COPY configuration/ ./configuration
COPY controllers/ ./controllers
COPY examples/ ./examples
COPY httperrors/ ./httperrors
COPY logger/ ./logger
COPY persistence/ ./persistence
COPY resources/ ./resources
COPY router/ ./router
COPY services/ ./services
COPY utils/ ./utils
COPY validators/ ./validators
COPY badaas.go .
COPY go.mod .
COPY go.sum .
RUN CGO_ENABLED=1 go build -a -o badaas .
RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot


# final image for end users
FROM alpine:3.16.2
COPY --from=builder /etc/passwd /etc/passwd
USER badaas
COPY --from=builder /app/badaas .
COPY docker/api/badaas.yml .
EXPOSE 8000
ENTRYPOINT ["./badaas", "--config_path", "badaas.yml"]