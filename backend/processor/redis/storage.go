package redisstore

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type NewsItem struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Link        string `json:"link"`
}

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr string) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Falha ao conectar ao Redis: %v", err)
	}

	return &RedisStorage{client: client}
}

func (r *RedisStorage) Save(news NewsItem) error {

	hash := sha256.Sum256([]byte(news.Link))
	key := fmt.Sprintf("news:%x", hash)

	jsonData, err := json.Marshal(news)
	if err != nil {
		return err
	}

	err = r.client.Set(ctx, key, jsonData, 24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("erro ao salvar no Redis: %w", err)
	}

	log.Printf("✅ Notícia salva no Redis: [%s]", news.Category)
	return nil
}
