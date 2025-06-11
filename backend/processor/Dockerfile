FROM golang:1.24.4
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY processor/ ./processor/
WORKDIR /app/processor
RUN go build -o processor .

CMD ["./processor"]
