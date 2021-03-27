package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post definition
type Post struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	From     string             `json:"from" bson:"from"`
	Message  string             `json:"message" bson:"message" validate:"required,max=255"`
	Likes    int                `json:"likes,omitempty" bson:"likes,omitempty"`
	LikedBy  []string           `json:"liked_by,omitempty" bson:"liked_by,omitempty"`
	Comments []Comment          `json:"comments" bson:"comments"`
}

// Comment definition
type Comment struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	From    string             `json:"from" bson:"from"`
	Content string             `json:"content" bson:"content" validate:"required,max=150"`
}
