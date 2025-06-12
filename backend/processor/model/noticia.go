package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Noticia struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Titulo    string             `bson:"titulo"`
	Link      string             `bson:"link"`
	Categoria string             `bson:"categoria"`
	CriadoEm  time.Time          `bson:"criado_em"`
	Fonte     string             `bson:"fonte"`
}
