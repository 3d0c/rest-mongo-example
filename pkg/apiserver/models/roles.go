package models

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RoleScheme nolint
type RoleScheme struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name,omitempty" json:"name"`
	Description string             `bson:"description,omitempty" json:"description"`
	Apps        []string           `bson:"apps,omitempty" json:"apps,omitempty"`
}

// Bind interface
// @TODO Add validation package
func (rs *RoleScheme) Bind(r *http.Request) error {
	var (
		am  *Application
		err error
	)

	if rs.Name == "" {
		return fmt.Errorf("role name is required")
	}

	if am, err = NewApplication(); err != nil {
		return fmt.Errorf("error initializing application model - %s", err)
	}

	for _, aid := range rs.Apps {
		if _, err = am.FindByID(aid); err != nil {
			if err == ErrNotFound {
				return fmt.Errorf("application '%s' is not exist", aid)
			}
			return fmt.Errorf("error finding application '%s' - %s", aid, err)
		}
	}

	rs.ID = primitive.NilObjectID

	return nil
}

// Role model
type Role struct {
	*base
}

// NewRole role model constructor
func NewRole() (*Role, error) {
	return &Role{
		base: &base{
			Collection: DB().Collection("roles"),
		},
	}, nil
}

// FindAll finds all roles
func (r *Role) FindAll() ([]RoleScheme, error) {
	var (
		result []RoleScheme
		elem   RoleScheme
		cursor *mongo.Cursor
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if cursor, err = r.Find(ctx, all); err != nil {
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
func (r *Role) FindByID(i interface{}) (*RoleScheme, error) {
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

	return r.find(bson.M{"_id": oid})
}

// find general find function
func (r *Role) find(match bson.M) (*RoleScheme, error) {
	var (
		result RoleScheme
		cursor *mongo.Cursor
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if cursor, err = r.Find(ctx, match); err != nil {
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

// Create creates new role document into `roles` collection
// returns oid as hex encoded string and error
func (r *Role) Create(role *RoleScheme) (string, error) {
	var (
		result *mongo.InsertOneResult
		oid    primitive.ObjectID
		ok     bool
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if result, err = r.InsertOne(ctx, role); err != nil {
		return "", err
	}

	if oid, ok = result.InsertedID.(primitive.ObjectID); !ok {
		return "", fmt.Errorf("error convering InsertID to ObjectId")
	}

	return oid.Hex(), nil
}
