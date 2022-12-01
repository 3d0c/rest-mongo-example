package models

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	all = bson.M{}

	// ErrNotFound error type
	ErrNotFound = errors.New("nothing found")
	// ErrSapUserNotFound SAP user validation error
	ErrSapUserNotFound = errors.New("SAP user not found")
)

type base struct {
	*mongo.Collection
}

// Delete is a common method for initialized collection
// @TODO clean all associated sessions in transaction
func (b *base) Delete(id string) error {
	var (
		oid primitive.ObjectID
		err error
	)

	if oid, err = primitive.ObjectIDFromHex(id); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if _, err = b.DeleteOne(ctx, bson.M{"_id": oid}); err != nil {
		return err
	}

	return nil
}

// Update generic update method
func (b *base) Update(id string, i interface{}) error {
	var (
		oid primitive.ObjectID
		err error
	)

	if oid, err = primitive.ObjectIDFromHex(id); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if _, err = b.ReplaceOne(ctx, bson.M{"_id": oid}, i); err != nil {
		return err
	}

	return nil
}
