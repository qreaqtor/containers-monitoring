FROM golang:1.23.3 AS builder
COPY . /app
WORKDIR /app/services
RUN go work sync
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./pinger/cmd/

FROM alpine:latest
COPY --from=builder /app/services/pinger/config /config
COPY --from=builder /bin/app /app

CMD ["/app"]