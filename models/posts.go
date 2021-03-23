package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Post definition
type Post struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	From    string             `json:"from" bson:"from"`
	Message string             `json:"message" bson:"message" validate:"required,max=255"`
}
