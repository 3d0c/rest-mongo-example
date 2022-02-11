package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AppDetails nolint
type AppDetails struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name,omitempty"`
	Path string             `bson:"path,omitempty"`
}

// ApplicationScheme model
type ApplicationScheme struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name,omitempty"`
	Path string             `bson:"path,omitempty"`
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
