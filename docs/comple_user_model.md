# Get complete user model

```javascript
db.users.aggregate([
	{
		$unwind: "$acl"
	},		
	{
		$lookup: {
			from: "applications",
			localField: "acl.application_id",
			foreignField: "_id",
			as: "acl.app_details"
		}
	},
	{
		$lookup: {
			from: "permissions",
			localField: "acl.permission_id",
			foreignField: "_id",
			as: "acl.perm_details"
		}
	},
	{
        $unwind: {
            path: "$acl.perm_details"
		}
	},
	{
        $unwind: {
            path: "$acl.app_details"
		}
	},
    {
        $group: {
            _id: '$_id',
            acl: {
                $push: '$acl'
            }
        }
    },	
    {
        $lookup: {
            from: 'users',
            localField: '_id',
            foreignField: '_id',
            as: 'aclDetails'
        }
    },
    {
        $unwind: {
            path: '$aclDetails'
        }
    },
    {
        $addFields: {
            'aclDetails.acl': '$acl'
        }
    },
    {
        $replaceRoot: {
            newRoot: '$aclDetails'
        }
    }    
]).pretty()
```
__Expected result__

```javascript
{
	"_id" : ObjectId("620527ed4a84ecd9ac78f623"),
	"name" : "admin",
	"email" : "root@dev.null",
	"acl" : [
		{
			"application" : ObjectId("620524994a84ecd9ac78f620"),
			"permission" : ObjectId("620524134a84ecd9ac78f61f"),
			"app_details" : {
				"_id" : ObjectId("620524994a84ecd9ac78f620"),
				"name" : "Sample One",
				"path" : "/sample"
			},
			"perm_details" : {
				"_id" : ObjectId("620524134a84ecd9ac78f61f"),
				"name" : "Full Access",
				"description" : "Full access to Application",
				"methods" : [
					"GET",
					"POST",
					"PUT",
					"DELETE"
				]
			}
		},
		{
			"application" : ObjectId("620527c04a84ecd9ac78f622"),
			"permission" : ObjectId("620524134a84ecd9ac78f61f"),
			"app_details" : {
				"_id" : ObjectId("620527c04a84ecd9ac78f622"),
				"name" : "Another one",
				"path" : "/another"
			},
			"perm_details" : {
				"_id" : ObjectId("620524134a84ecd9ac78f61f"),
				"name" : "Full Access",
				"description" : "Full access to Application",
				"methods" : [
					"GET",
					"POST",
					"PUT",
					"DELETE"
				]
			}
		}
	]
}
```
