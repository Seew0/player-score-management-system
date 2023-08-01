package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Player struct {
	ID      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name    string             `bson:"name,omitempty" json:"name"`
	Country string             `bson:"country,omitempty" json:"country"`
	Score   int                `bson:"score,omitempty" json:"score"`
}
