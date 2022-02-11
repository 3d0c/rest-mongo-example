package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PermDetails nolint
type PermDetails struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
	Methods     []string           `bson:"methods,omitempty"`
}

// PermissionDetails nolint
type PermissionScheme struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
	Methods     []string           `bson:"methods,omitempty"`
}

// Permission model
type Permission struct {
	*base
	*mongo.Collection
}

// NewPermission permission model constructor
func NewPermission() (*Permission, error) {
	return &Permission{
		base:       &base{},
		Collection: DB().Collection("permissions"),
	}, nil
}
