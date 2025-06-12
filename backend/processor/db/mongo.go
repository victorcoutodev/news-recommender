package db

import (
	"context"
	"log"
	"news-recommender/m/processor/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Cliente *mongo.Client

func ConectarMongo(uri string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	log.Println("✅ Conectado ao MongoDB com sucesso!")
	Cliente = client
	return nil
}

// SalvarNoticia insere uma notícia na coleção "noticias" se ainda não existir pelo campo Link
func SalvarNoticia(n model.Noticia) error {
	collection := Cliente.Database("newsdb").Collection("noticias")

	if n.CriadoEm.IsZero() {
		n.CriadoEm = time.Now()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"link": n.Link}
	update := bson.M{"$setOnInsert": n}
	opts := options.Update().SetUpsert(true)

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Printf("Erro ao salvar no MongoDB: %v", err)
		return err
	}

	if result.UpsertedCount > 0 {
		log.Printf("✅ Notícia inserida no MongoDB: %s", n.Titulo)
	} else {
		log.Printf("ℹ️ Notícia já existe, não inserida: %s", n.Titulo)
	}

	return nil
}
