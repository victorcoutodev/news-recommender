package main

import (
	"context"
	"encoding/json"
	"log"
	"news-recommender/m/processor/classifier"
	"os"

	redisstore "news-recommender/m/processor/redis"

	"github.com/segmentio/kafka-go"
)

func main() {
	// 0. Configuração do Kafka Broker
	// Obtém o endereço do broker a partir da variável de ambiente
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9093" // fallback para desenvolvimento local
	}

	// 1. Inicializa o classificador
	classifier := classifier.NewOllamaClassifier("mistral")

	// 2. Configura o Kafka Consumer
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "raw-news",
	})
	defer reader.Close()

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379"
	}
	storage := redisstore.NewRedisStorage(addr)

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Erro no Kafka: %v", err)
			continue
		}

		var news map[string]string
		if err := json.Unmarshal(msg.Value, &news); err != nil {
			log.Printf("Erro ao decodificar JSON: %v", err)
			continue
		}

		// 4. Classificação com Ollama
		category, err := classifier.Classify(news["title"] + " " + news["description"])
		if err != nil {
			log.Printf("Erro na classificação: %v", err)
			category = "erro"
		}

		news["category"] = category

		newsItem := redisstore.NewsItem{
			Title:       news["title"],
			Description: news["description"],
			Category:    category,
			Link:        news["link"],
		}

		err = storage.Save(newsItem)
		if err != nil {
			log.Printf("Erro ao salvar no Redis: %v", err)
		}

		log.Printf("ClassificadoTest: %s → %s", news["title"], category)
	}
}
