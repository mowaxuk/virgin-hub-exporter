# 🛠 Build stage: compile the Go binary using Go 1.23
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o exporter

# 🏃 Runtime stage: minimal Alpine container with just the binary
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/exporter .
EXPOSE 9877
ENTRYPOINT ["./exporter"]
