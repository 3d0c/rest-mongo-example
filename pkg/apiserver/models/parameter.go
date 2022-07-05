package models

import (
	"context"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ParameterScheme type
type ParameterScheme struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description"`
	CompType    string             `bson:"comp_type" json:"comp_type"`
	DataType    string             `bson:"data_type" json:"data_type"`
	Options     []string           `bson:"options" json:"options"`
}

// Bind interface
// @TODO Add validation package
func (ps *ParameterScheme) Bind(r *http.Request) error {
	if ps.Name == "" {
		return fmt.Errorf("parameter name is required")
	}

	ps.ID = primitive.NilObjectID

	return nil
}

// Parameter model
type Parameter struct {
	*base
}

// NewParameter parameter model constructor
func NewParameter() (*Parameter, error) {
	return &Parameter{
		base: &base{
			Collection: DB().Collection("parameters"),
		},
	}, nil
}

// FindAll finds all parameters
func (p *Parameter) FindAll() ([]ParameterScheme, error) {
	var (
		result []ParameterScheme
		elem   ParameterScheme
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
func (p *Parameter) FindByID(i interface{}) (*ParameterScheme, error) {
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
func (p *Parameter) find(match bson.M) (*ParameterScheme, error) {
	var (
		result ParameterScheme
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

// Create creates new parameter document into `parameters` collection
// returns oid as hex encoded string and error
func (p *Parameter) Create(parameter *ParameterScheme) (string, error) {
	var (
		result *mongo.InsertOneResult
		oid    primitive.ObjectID
		ok     bool
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if result, err = p.InsertOne(ctx, parameter); err != nil {
		return "", err
	}

	if oid, ok = result.InsertedID.(primitive.ObjectID); !ok {
		return "", fmt.Errorf("error convering InsertID to ObjectId")
	}

	return oid.Hex(), nil
}
