FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN cd ./cmd/server && GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -o app

FROM scratch
WORKDIR /app
COPY --from=builder /app/cmd/server/app .
COPY --from=builder /app/.env .
EXPOSE 8080
ENTRYPOINT [ "./app" ]