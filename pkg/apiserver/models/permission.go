package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PermissionScheme nolint
type PermissionScheme struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name,omitempty" json:"name"`
	Description string             `bson:"description,omitempty" json:"description"`
	Methods     []string           `bson:"methods,omitempty" json:"methods,omitempty"`
}

// IsAllowed checks is method allowed or not
func (p *PermissionScheme) IsAllowed(method string) bool {
	for i := range p.Methods {
		if p.Methods[i] == method {
			return true
		}
	}

	return false
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
