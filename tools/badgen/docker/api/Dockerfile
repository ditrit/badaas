# builder image
FROM golang:1.19-alpine AS builder
RUN apk add build-base
# TODO
WORKDIR /app
COPY badaas/ ../badaas
COPY example/ .
RUN CGO_ENABLED=1 go build -a -o badaas .


# final image for end users
FROM alpine:3.16.2
RUN addgroup -S badaas \
    && adduser -S badaas -G badaas
USER badaas
COPY --from=builder /app/badaas .
COPY example/badaas/docker/api/badaas.yml .
EXPOSE 8000
ENTRYPOINT ["./badaas", "--config_path", "badaas.yml"]