FROM golang:1.22 AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o stress-test .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/stress-test .

ENTRYPOINT ["./stress-test"]