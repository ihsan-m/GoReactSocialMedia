package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty"`
	AuthorID  primitive.ObjectID   `bson:"author_id"`
	Caption   string               `bson:"caption"`
	Emojis    []string             `bson:"emojis"`
	Gifs      []string             `bson:"gifs"`
	Photos    []string             `bson:"photos"`
	Location  string               `bson:"location"`
	Platforms []string             `bson:"platforms"`
	Comments  []primitive.ObjectID `bson:"comments"`
}
