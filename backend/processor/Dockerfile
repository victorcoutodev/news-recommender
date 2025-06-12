# Stage 1: Build do binário Go
FROM golang:1.24.4 AS builder

WORKDIR /app

# Copia dependências e baixa
COPY go.mod go.sum ./
RUN go mod download

# Copia o código do processor e builda
COPY . . 
WORKDIR /app/processor
RUN go build -o processor main.go

# Stage 2: Imagem final com runtime leve
FROM ubuntu:22.04

WORKDIR /app

# Instala curl e outros pacotes básicos necessários
RUN apt-get update && apt-get install -y curl ca-certificates && rm -rf /var/lib/apt/lists/*

# Instala Ollama (mantém igual)
RUN curl -L https://ollama.ai/install.sh | sh

# Puxa o modelo ONDE? Ideal: já durante build para evitar puxar toda vez no container iniciar
# Então vamos puxar o modelo durante o build, para não travar no startup

# Copia o binário do builder
COPY --from=builder /app/processor/processor .

# Comentado pois isso deve rodar no host antes do docker subir
# RUN ollama pull mistral

# Define variável de ambiente (exemplo, ajuste conforme necessidade)
ENV KAFKA_BROKER=kafka:9092

# Comando para iniciar o app diretamente
CMD ["./processor"]
