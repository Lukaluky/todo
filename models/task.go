package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Task struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	Completed bool               `bson:"completed" json:"completed"`
	UserID    uint               `bson:"user_id" json:"user_id"`
}