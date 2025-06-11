package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/segmentio/kafka-go"
)

func main() {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://g1.globo.com/rss/g1/")
	if err != nil {
		log.Fatal(err)
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "raw-news",
	})

	for _, item := range feed.Items {
		news := map[string]string{
			"title": item.Title,
			"link":  item.Link,
			"date":  item.Published,
		}
		jsonNews, _ := json.Marshal(news)

		err := writer.WriteMessages(context.Background(),
			kafka.Message{Value: jsonNews},
		)
		if err != nil {
			log.Printf("Erro ao enviar: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
}
