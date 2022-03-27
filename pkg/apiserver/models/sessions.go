package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/teal-seagull/lyre-be-v4/pkg/helpers"
)

// SessionSchemeType name for context
type SessionSchemeType struct{}

// SessionScheme params
// @index token uinq
// @index user_id
// @@index (token, user_id) uniq
type SessionScheme struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	UserID primitive.ObjectID `bson:"user_id,omitempty" json:"-"`
	Token  string             `bson:"token" json:"token"`
}

// Session model
type Session struct {
	*base
	*mongo.Collection
}

// NewSession Session model constructor
func NewSession() (*Session, error) {
	return &Session{
		base:       &base{},
		Collection: DB().Collection("sessions"),
	}, nil
}

// Exists checks token exists
// @index token
func (s *Session) Exists(token string) (bool, error) {
	var (
		count int64
		err   error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if count, err = s.CountDocuments(ctx, bson.M{"token": token}); err != nil {
		return false, err
	}

	return count > 0, nil
}

// Create session
func (s *Session) Create(uid primitive.ObjectID) (*SessionScheme, error) {
	var (
		token string
		err   error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if token, err = helpers.CreateToken(uid.Hex()); err != nil {
		return nil, err
	}

	if _, err = s.InsertOne(ctx, bson.M{"user_id": uid, "token": token}); err != nil {
		return nil, err
	}

	return &SessionScheme{Token: token}, nil
}

// Remove session by match
func (s *Session) Remove(i interface{}) error {
	var (
		match bson.M
		err   error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	switch v := i.(type) {
	case string:
		match = bson.M{"token": i.(string)}

	case primitive.ObjectID:
		match = bson.M{"_id": i.(primitive.ObjectID)}

	default:
		return fmt.Errorf("wrong input type '%s', expecting (string) or (ObjectID)", v)
	}

	_, err = s.DeleteOne(ctx, match)

	return err
}
