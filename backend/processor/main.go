package main

import (
	"context"
	"encoding/json"
	"log"
	"news-recommender/m/processor/classifier"
	"news-recommender/m/processor/db"
	"news-recommender/m/processor/model"
	"os"
	"time"

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

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	err := db.ConectarMongo(mongoURI)
	if err != nil {
		log.Fatalf("❌ Erro ao conectar ao MongoDB: %v", err)
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

	// 3. Loop de leitura de mensagens do Kafka
	log.Println("🔄 Iniciando o processamento de notícias...")
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

		// 5. Monta o objeto NewsItem para o Redis
		newsItem := redisstore.NewsItem{
			Title:       news["title"],
			Description: news["description"],
			Category:    category,
			Link:        news["link"],
		}

		// 6. Salva no Redis
		err = storage.Save(newsItem)
		if err != nil {
			log.Printf("Erro ao salvar no Redis: %v", err)
		}

		// 7. Monta o objeto Noticia para o Mongo
		noticia := model.Noticia{
			Titulo:    news["title"],
			Link:      news["link"],
			Categoria: category,
			Fonte:     "G1",
			CriadoEm:  time.Now(),
		}

		// 8. Salva no MongoDB
		err = db.SalvarNoticia(noticia)
		if err != nil {
			log.Printf("Erro ao salvar no MongoDB: %v", err)
		}

		// 9. Log de sucesso
		log.Printf("ClassificadoTest: %s → %s", news["title"], category)
	}
}
