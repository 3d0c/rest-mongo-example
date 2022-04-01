package models

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/teal-seagull/lyre-be-v4/pkg/helpers"
)

// UserSchemeType name for context
type UserSchemeType struct{}

// UserScheme params
type UserScheme struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"user_name"`
	Email    string             `bson:"email" json:"email"`
	Password *string            `bson:"password" json:"password,omitempty"`
	ACL      []ACLScheme        `bson:"acl" json:"acl"`
}

// Bind interface
// TODO Add validation package
func (u *UserScheme) Bind(r *http.Request) error {
	if u.Password == nil {
		return fmt.Errorf("password is required")
	}

	u.ID = primitive.NilObjectID

	return nil
}

// GetPermission returns permission by application path
func (u *UserScheme) GetPermission(path string) *PermissionScheme {
	for _, acl := range u.ACL {
		if acl.Application.Path == path {
			return acl.Permissions
		}
	}

	return nil
}

// User model
type User struct {
	*base
}

// NewUser user model constructor
func NewUser() (*User, error) {
	return &User{
		base: &base{
			Collection: DB().Collection("users"),
		},
	}, nil
}

// FindByID add $match by id
func (u *User) FindByID(i interface{}) (*UserScheme, error) {
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

	return u.find(bson.M{"_id": oid})
}

// FindByName find by email or name
func (u *User) FindByName(us *UserScheme) (*UserScheme, error) {
	var (
		match bson.M
		err   error
	)

	if us.Email != "" {
		match = bson.M{"email": us.Email}
	}
	if us.Name != "" {
		match = bson.M{"name": us.Name}
	}

	if us, err = u.find(match); err != nil {
		return nil, err
	}

	return us, nil
}

// FindAll finds all
func (u *User) FindAll() ([]UserScheme, error) {
	var (
		result []UserScheme
		elem   UserScheme
		cursor *mongo.Cursor
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if cursor, err = u.Aggregate(ctx, completeUserModel(all)); err != nil {
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

// find general find user function, expects $match returns complete user model
// @index uniq name
// @index uniq email
func (u *User) find(match bson.M) (*UserScheme, error) {
	var (
		result UserScheme
		cursor *mongo.Cursor
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if cursor, err = u.Aggregate(ctx, completeUserModel(match)); err != nil {
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

// Create creates new user document into `users` collection
// returns oid as hex encoded string and error
// @TODO add validation for IDs of applications and permission check its existence
func (u *User) Create(user *UserScheme) (string, error) {
	var (
		result *mongo.InsertOneResult
		oid    primitive.ObjectID
		ok     bool
		err    error
	)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Please note, that Password at this step is plain, hash it first
	if *user.Password, err = helpers.HashPassword(*user.Password); err != nil {
		return "", err
	}

	if result, err = u.InsertOne(ctx, user); err != nil {
		return "", err
	}

	if oid, ok = result.InsertedID.(primitive.ObjectID); !ok {
		return "", fmt.Errorf("error convering InsertID to ObjectId")
	}

	return oid.Hex(), nil
}

// MarshalJSON cleans password field on response
// as a result, because of "omitempty", there is no "password field" in respone
func (u *UserScheme) MarshalJSON() ([]byte, error) {
	type tmp UserScheme

	u.Password = nil

	return json.Marshal(&struct {
		*tmp
	}{
		tmp: (*tmp)(u),
	})
}
