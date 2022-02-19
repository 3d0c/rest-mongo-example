package models

import (
	"go.mongodb.org/mongo-driver/bson"
)

func completeUserModel() []bson.M {
	return []bson.M{
		{
			"$unwind": "$acl",
		},
		{
			"$lookup": bson.M{
				"from":         "applications",
				"localField":   "acl.application_id",
				"foreignField": "_id",
				"as":           "acl.application",
			},
		},
		{
			"$lookup": bson.M{
				"from":         "permissions",
				"localField":   "acl.permission_id",
				"foreignField": "_id",
				"as":           "acl.permissions",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$acl.permissions",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$acl.application",
			},
		},
		{
			"$group": bson.M{
				"_id": "$_id",
				"acl": bson.M{
					"$push": "$acl",
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "_id",
				"foreignField": "_id",
				"as":           "aclDetails",
			},
		},
		{
			"$unwind": bson.M{
				"path": "$aclDetails",
			},
		},
		{
			"$addFields": bson.M{
				"aclDetails.acl": "$acl",
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$aclDetails",
			},
		},
	}
}