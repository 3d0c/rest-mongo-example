package models

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/teal-seagull/lyre-be-v4/pkg/helpers"
)

// PermissionScheme nolint
type PermissionScheme struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name,omitempty" json:"name"`
	Description string             `bson:"description,omitempty" json:"description"`
	Methods     []string           `bson:"methods,omitempty" json:"methods,omitempty"`
}

// IsAllowed checks is method allowed or not
func (ps *PermissionScheme) IsAllowed(method string) bool {
	for i := range ps.Methods {
		if ps.Methods[i] == method {
			return true
		}
	}

	return false
}

// Bind interface
// @TODO Add validation package
func (ps *PermissionScheme) Bind(r *http.Request) error {
	if ps.Name == "" {
		return fmt.Errorf("permission name is required")
	}

	for _, m := range ps.Methods {
		if !helpers.IsValidMethod(m) {
			return fmt.Errorf("unexpected method '%s'", m)
		}
	}

	ps.ID = primitive.NilObjectID

	return nil
}

// Permission model
type Permission struct {
	*base
}

// NewPermission permission model constructor
func NewPermission() (*Permission, error) {
	return &Permission{
		base: &base{
			Collection: DB().Collection("permissions"),
		},
	}, nil
}

// FindAll finds all permissions
func (p *Permission) FindAll() ([]PermissionScheme, error) {
	var (
		result []PermissionScheme
		elem   PermissionScheme
		cursor *mongo.Cursor
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if cursor, err = p.Find(ctx, all); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		if err = cursor.Decode(&elem); err != nil {
			return nil, err
		}

		result = append(result, elem)
	}

	return result, nil
}

// FindByID add $match by id
func (p *Permission) FindByID(i interface{}) (*PermissionScheme, error) {
	var (
		oid primitive.ObjectID
		err error
	)

	switch v := i.(type) {
	case string:
		if oid, err = primitive.ObjectIDFromHex(i.(string)); err != nil {
			return nil, err
		}

	case primitive.ObjectID:
		oid = i.(primitive.ObjectID)

	default:
		return nil, fmt.Errorf("wrong input type '%s', expecting (string) or (ObjectID)", v)
	}

	return p.find(bson.M{"_id": oid})
}

// find general find function
func (p *Permission) find(match bson.M) (*PermissionScheme, error) {
	var (
		result PermissionScheme
		cursor *mongo.Cursor
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if cursor, err = p.Find(ctx, match); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err = cursor.Decode(&result); err != nil {
			return nil, err
		}
		// returns result only if cursor.Next is true
		return &result, nil
	}
	// else returns nil. to prevent initialized but empty structure
	return nil, fmt.Errorf("nothing found")
}

// Create creates new permission document into `permissions` collection
// returns oid as hex encoded string and error
func (p *Permission) Create(app *PermissionScheme) (string, error) {
	var (
		result *mongo.InsertOneResult
		oid    primitive.ObjectID
		ok     bool
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if result, err = p.InsertOne(ctx, app); err != nil {
		return "", err
	}

	if oid, ok = result.InsertedID.(primitive.ObjectID); !ok {
		return "", fmt.Errorf("error convering InsertID to ObjectId")
	}

	return oid.Hex(), nil
}
