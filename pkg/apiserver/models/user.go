package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserSchemeType name for context
type UserSchemeType struct{}

// UserScheme params
type UserScheme struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
	ACL   []ACLScheme        `bson:"acl"`
}

// GetPermission returns permission by application name (mean path)
func (u *UserScheme) GetPermission(name string) *PermissionScheme {
	for _, acl := range u.ACL {
		if acl.Application.Path == name {
			return acl.Permissions
		}
	}

	return nil
}

// User model
type User struct {
	*base
	*mongo.Collection
}

// NewUser user model constructor
func NewUser() (*User, error) {
	return &User{
		base:       &base{},
		Collection: DB().Collection("users"),
	}, nil
}

// Find single user by id
// TODO ADD USERID AS AN ARG AND PASS IT INTO completeUserModel as $match!!!
func (u *User) Find(userID string) (*UserScheme, error) {
	var (
		result UserScheme
		cursor *mongo.Cursor
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if cursor, err = u.Aggregate(ctx, completeUserModel()); err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if cursor.Next(ctx) {
		if err = cursor.Decode(&result); err != nil {
			return nil, err
		}
	}

	return &result, nil
}
