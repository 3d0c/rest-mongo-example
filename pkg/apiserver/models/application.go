package models

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ApplicationScheme model
// @index uniq name
// @index uniq path
// @index (name,path) uniq
type ApplicationScheme struct {
	ID   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name string             `bson:"name,omitempty" json:"name"`
	Path string             `bson:"path,omitempty" json:"path"`
}

// Bind interface
// @TODO Add validation package
// @TODO validate application path, must contain "/"
func (as *ApplicationScheme) Bind(r *http.Request) error {
	if as.Name == "" {
		return fmt.Errorf("application name is required")
	}
	if as.Path == "" {
		return fmt.Errorf("application path is required")
	}

	as.ID = primitive.NilObjectID

	return nil
}

// Application model
type Application struct {
	*base
}

// NewApplication permission model constructor
func NewApplication() (*Application, error) {
	return &Application{
		base: &base{
			Collection: DB().Collection("applications"),
		},
	}, nil
}

// FindAll finds all applications
func (a *Application) FindAll() ([]ApplicationScheme, error) {
	var (
		result []ApplicationScheme
		elem   ApplicationScheme
		cursor *mongo.Cursor
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if cursor, err = a.Find(ctx, all); err != nil {
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
func (a *Application) FindByID(i interface{}) (*ApplicationScheme, error) {
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

	return a.find(bson.M{"_id": oid})
}

// find general find function
func (a *Application) find(match bson.M) (*ApplicationScheme, error) {
	var (
		result ApplicationScheme
		cursor *mongo.Cursor
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if cursor, err = a.Find(ctx, match); err != nil {
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

// Create creates new application document into `applications` collection
// returns oid as hex encoded string and error
// @TODO add validation for IDs of applications and permission check its existence
//       could be done as exception of uniq index
func (a *Application) Create(app *ApplicationScheme) (string, error) {
	var (
		result *mongo.InsertOneResult
		oid    primitive.ObjectID
		ok     bool
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if result, err = a.InsertOne(ctx, app); err != nil {
		return "", err
	}

	if oid, ok = result.InsertedID.(primitive.ObjectID); !ok {
		return "", fmt.Errorf("error convering InsertID to ObjectId")
	}

	return oid.Hex(), nil
}

// Update updates specific application, pretty the same as Create method
func (a *Application) Update(id string, app *ApplicationScheme) error {
	var (
		oid primitive.ObjectID
		err error
	)

	if oid, err = primitive.ObjectIDFromHex(id); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if _, err = a.ReplaceOne(ctx, bson.M{"_id": oid}, app); err != nil {
		return err
	}

	return nil
}
