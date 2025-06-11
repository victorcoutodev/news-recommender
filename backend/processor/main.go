package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func classifyWithAI(text string) string {
	// Simulação - Substitua pela chamada real à API do DeepSeek
	return "tecnologia"
}

func main() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "raw-news",
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Erro: %v", err)
			continue
		}

		var news map[string]string
		json.Unmarshal(msg.Value, &news)

		// Classificação com IA (exemplo simplificado)
		category := classifyWithAI(news["title"])
		news["category"] = category

		// Salvar no MongoDB (próximo passo)
		log.Printf("Notícia classificada: %s -> %s", news["title"], category)
	}
}
