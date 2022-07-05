package helpers

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Typ returns ObjectID from string or oid
func Typ(i interface{}) (primitive.ObjectID, error) {
	var (
		oid primitive.ObjectID
		err error
	)

	switch v := i.(type) {
	case string:
		if oid, err = primitive.ObjectIDFromHex(i.(string)); err != nil {
			return primitive.NilObjectID, err
		}

	case primitive.ObjectID:
		oid = i.(primitive.ObjectID)

	default:
		return primitive.NilObjectID, fmt.Errorf("wrong input type '%s', expecting (string) or (ObjectID)", v)
	}

	return oid, nil
}
