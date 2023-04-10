package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID       primitive.ObjectID  `bson:"_id,omitempty"`
	AuthorID primitive.ObjectID  `bson:"author_id"`
	PostID   primitive.ObjectID  `bson:"post_id"`
	ParentID *primitive.ObjectID `bson:"parent_id,omitempty"`
	Content  string              `bson:"content"`
}
