# Get complete user model @TODO

```javascript
db.users.aggregate([
    {
        $unwind: "$roles"
    },
    {
        $lookup: {
            from: "roles",
            localField: "roles",
            foreignField: "_id",
            as: "acl"
        }
    },
    {
        $unwind: "$acl"
    },
    {
        $unwind: "$acl.apps"
    },
    {
        $lookup: {
            from: "applications",
            localField: "acl.apps",
            foreignField: "_id",
            as: "acl.apps"
        }   
    },
    {
        $unwind: "$acl.apps"
    },
    {
        $group: {
            _id: '$_id',
            acl: {
                $push: '$acl.apps'
            }
        }
    },
    {
        $lookup: {
            from: 'users',
            localField: '_id',
            foreignField: '_id',
            as: 'userDetails'
        }
    },
    {
        $unwind: "$userDetails"
    },
    {
        $addFields: {
            'userDetails.acl': '$acl'
        }
    },
    {
        $replaceRoot: {
            newRoot: '$userDetails'
        }
    }
]).pretty()
```
__Expected result__

```javascript
{
	"_id" : ObjectId("6243308399fd59c176c52fd4"),
	"name" : "admin",
	"email" : "root@dev.null",
	"password" : "$2a$11$lAT02Pq3MiHefYLYM6ZrUO79swRZAHeE0x0/RX13lIRouX72Hzwr2",
	"roles" : [
		ObjectId("6242d43e99fd59c176c52fd4"),
		ObjectId("62a4a90e5d2fe1005407cd8b")
	],
	"acl" : [
		{
			"_id" : ObjectId("6242d43e99fd59c176c52fd3"),
			"name" : "User management application",
			"path" : "/users"
		},
		{
			"_id" : ObjectId("6245984799fd59c176c52fd5"),
			"name" : "Applications management",
			"path" : "/applications"
		},
		{
			"_id" : ObjectId("6246d923ad35f14740a5fa79"),
			"name" : "Permissions management",
			"path" : "/permissions"
		},
		{
			"_id" : ObjectId("62937983ec569fe63ccffebc"),
			"name" : "Roles management",
			"path" : "/roles"
		},
		{
			"_id" : ObjectId("6242d43e99fd59c176c52fd3"),
			"name" : "User management application",
			"path" : "/users"
		}
	]
}
```
