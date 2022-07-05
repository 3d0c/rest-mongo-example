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

// SettingScheme type
// @TODO add expanded User and ParameterEx meta fields
type SettingScheme struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	UserID     primitive.ObjectID   `bson:"user_id,omitempty" json:"user_id"`
	AppID      primitive.ObjectID   `bson:"app_id,omitempty" json:"app_id"`
	Parameters []primitive.ObjectID `bson:"parameters,omitempty" json:"parameters"`
}

// Bind interface
// @TODO Add validation package
func (s *SettingScheme) Bind(r *http.Request) error {
	if s.UserID.IsZero() {
		return fmt.Errorf("user id is required")
	}
	if s.AppID.IsZero() {
		return fmt.Errorf("application id is required")
	}
	if len(s.Parameters) == 0 {
		return fmt.Errorf("parameters array can't be empty")
	}

	s.ID = primitive.NilObjectID

	return nil
}

// Setting model
type Setting struct {
	*base
}

// NewSetting parameter model constructor
func NewSetting() (*Setting, error) {
	return &Setting{
		base: &base{
			Collection: DB().Collection("settings"),
		},
	}, nil
}

// FindAll finds all settings
func (s *Setting) FindAll(userID, appID interface{}) ([]SettingScheme, error) {
	var (
		result []SettingScheme
		elem   SettingScheme
		cursor *mongo.Cursor
		oid    primitive.ObjectID
		match  = make([]bson.M, 0)
		and    = make([]bson.M, 0)
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if oid, err = helpers.Typ(userID); err == nil && !oid.IsZero() {
		and = append(and, bson.M{"user_id": oid})
	}
	if oid, err = helpers.Typ(appID); err == nil && !oid.IsZero() {
		and = append(and, bson.M{"app_id": oid})
	}

	if len(and) > 0 {
		match = append(match, bson.M{"$match": bson.M{"$and": and}})
	}

	if cursor, err = s.Aggregate(ctx, match); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		if err = cursor.Decode(&elem); err != nil {
			return nil, err
		}

		result = append(result, elem)
	}

	if len(result) == 0 {
		return nil, helpers.ErrNotFound
	}

	return result, nil
}

// FindByID add $match by id
func (s *Setting) FindByID(i interface{}) (*SettingScheme, error) {
	var (
		oid primitive.ObjectID
		err error
	)

	if oid, err = helpers.Typ(i); i != nil || oid.IsZero() {
		return nil, fmt.Errorf("wrong setting id - %s", err)
	}

	return s.find(bson.M{"_id": oid})
}

// find general find function
func (s *Setting) find(match bson.M) (*SettingScheme, error) {
	var (
		result SettingScheme
		cursor *mongo.Cursor
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if cursor, err = s.Find(ctx, match); err != nil {
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

// Create creates new Setting document into `Settings` collection
// returns oid as hex encoded string and error
func (s *Setting) Create(setting *SettingScheme) (string, error) {
	var (
		result *mongo.InsertOneResult
		oid    primitive.ObjectID
		ok     bool
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if result, err = s.InsertOne(ctx, setting); err != nil {
		return "", err
	}

	if oid, ok = result.InsertedID.(primitive.ObjectID); !ok {
		return "", fmt.Errorf("error convering InsertID to ObjectId")
	}

	return oid.Hex(), nil
}
