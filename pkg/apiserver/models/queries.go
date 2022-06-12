package models

import (
	"go.mongodb.org/mongo-driver/bson"
)

func completeUserModel(match bson.M) []bson.M {
	return []bson.M{
		{
			"$unwind": "$roles",
		},
		{
			"$lookup": bson.M{
				"from":         "roles",
				"localField":   "roles",
				"foreignField": "_id",
				"as":           "acl",
			},
		},
		{
			"$unwind": "$acl",
		},
		{
			"$unwind": "$acl.apps",
		},
		{
			"$lookup": bson.M{
				"from":         "applications",
				"localField":   "acl.apps",
				"foreignField": "_id",
				"as":           "acl.apps",
			},
		},
		{
			"$unwind": "$acl.apps",
		},
		{
			"$group": bson.M{
				"_id": "$_id",
				"acl": bson.M{
					"$push": "$acl.apps",
				},
			},
		},
		{
			"$lookup": bson.M{
				"from":         "users",
				"localField":   "_id",
				"foreignField": "_id",
				"as":           "userDetails",
			},
		},
		{
			"$unwind": "$userDetails",
		},
		{
			"$addFields": bson.M{
				"userDetails.acl": "$acl",
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": "$userDetails",
			},
		},

		{
			"$match": match,
		},
	}
}
