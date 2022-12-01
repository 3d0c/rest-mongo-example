package models

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/teal-seagull/lyre-be-v4/pkg/config"
	"github.com/teal-seagull/lyre-be-v4/pkg/helpers"
	"github.com/teal-seagull/lyre-be-v4/pkg/sap"
)

// UserSchemeType name for context
type UserSchemeType struct{}

// UserScheme params
type UserScheme struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	Name        string               `bson:"username" json:"user_name"`
	Email       string               `bson:"email" json:"email"`
	Password    *string              `bson:"password" json:"password,omitempty"`
	ACL         []ApplicationScheme  `bson:"acl" json:"acl,omitempty"`
	Roles       []primitive.ObjectID `bson:"roles" json:"roles"`
	Settings    []primitive.ObjectID `bson:"settings" json:"settings"`
	Avatar      string               `bson:"avatar" json:"avatar"`
	FirstName   string               `bson:"first_name" json:"first_name"`
	LastName    string               `bson:"last_name" json:"last_name"`
	PlantID     string               `bson:"plant_id" json:"plant_id"`
	Printer     string               `bson:"printer" json:"printer"`
	LastLogin   time.Time            `bson:"last_login" json:"last_login,omitempty"`
	CreatedBy   string               `bson:"created_by" json:"created_by"`
	CreatedDate time.Time            `bson:"created_date" json:"created_date,omitempty"`
	UpdatedBy   string               `bson:"updated_by" json:"updated_by"`
	UpdatedDate time.Time            `bson:"updated_date" json:"updated_date,omitempty"`
	Status      int                  `bson:"status" json:"status"`
	SapID       string               `bson:"sap_id" json:"sap_id"`
	SapPwd      string               `bson:"sap_pwd" json:"sap_pwd"`
}

// Bind interface
// TODO Add validation package
func (u *UserScheme) Bind(r *http.Request) error {
	var (
		fileName string
		err      error
	)

	if config.TheConfig().SAP.ValidateUser {
		if err = sap.ValidateUser(u.SapID, u.SapPwd); err != nil {
			return ErrSapUserNotFound
		}
	}

	u.ID = primitive.NilObjectID

	if u.Password == nil {
		return fmt.Errorf("password is required")
	}

	if u.Avatar != "" {
		if fileName, err = helpers.ParseAndSaveImage(u.Avatar); err != nil {
			return err
		}

		u.Avatar = fileName
	}

	return nil
}

// IsAllowed checks is it allowed application for user
func (u *UserScheme) IsAllowed(path string) bool {
	for i := range u.ACL {
		if u.ACL[i].Path == path {
			return true
		}
	}

	return false
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
		match = bson.M{"username": us.Name}
	}

	if us, err = u.find(match); err != nil {
		return nil, err
	}

	return us, nil
}

// FindAll finds all
func (u *User) FindAll(roleID string) ([]UserScheme, error) {
	var (
		result []UserScheme
		elem   UserScheme
		cursor *mongo.Cursor
		match  bson.M = all
		oid    primitive.ObjectID
		err    error
	)

	if roleID != "" {
		if oid, err = primitive.ObjectIDFromHex(roleID); err != nil {
			return nil, err
		}
		match = bson.M{"roles": oid}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if cursor, err = u.Aggregate(ctx, completeUserModel(match)); err != nil {
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
// as a result, because of "omitempty", there is no "password field" in response
func (u *UserScheme) MarshalJSON() ([]byte, error) {
	type tmp UserScheme

	u.Password = nil

	return json.Marshal(&struct {
		*tmp
	}{
		tmp: (*tmp)(u),
	})
}
