FROM golang:1.24.4
WORKDIR /app

COPY go.mod go.sum ./ 
RUN go mod download

COPY collector/ ./collector/
WORKDIR /app/collector
RUN go build -o collector .

CMD ["./collector"]
