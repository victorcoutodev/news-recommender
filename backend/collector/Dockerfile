# Stage 1: build
FROM golang:1.24.4 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . . 
WORKDIR /app/collector
RUN go build -o collector main.go

# Stage 2: imagem final
FROM ubuntu:22.04

WORKDIR /app
RUN apt-get update && apt-get install -y curl ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/collector/collector .

CMD ["./collector"]
