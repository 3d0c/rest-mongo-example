package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ApplicationScheme model
type ApplicationScheme struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name,omitempty" json:"name"`
	Path string             `bson:"path,omitempty" json:"path"`
}

// Application model
type Application struct {
	*base
	*mongo.Collection
}

// NewApplication permission model constructor
func NewApplication() (*Application, error) {
	return &Application{
		base:       &base{},
		Collection: DB().Collection("applications"),
	}, nil
}
