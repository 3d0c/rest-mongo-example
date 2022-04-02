/*
    Please, note, that _id fields are hardcoded. This because of using
    standard uniqe object id as a refernce. Following _ids are unique.
    So it's safe to use this provisioning.

    Users
        - admin:default_password

    @TODO add propper indices.
*/
db = db.getSiblingDB('v4');

db.createCollection('permissions');
db.permissions.insertMany([
    {
        "_id" : ObjectId("620524134a84ecd9ac78f61d"),
        "name" : "Read-Only",
        "description" : "Read only access to Application",
        "methods" : [
            "GET"
        ]
    },
    {
        "_id" : ObjectId("620524134a84ecd9ac78f61e"),
        "name" : "Read-Write",
        "description" : "Read and Create access to Application",
        "methods" : [
            "GET",
            "POST"
        ]
    },
    {
        "_id" : ObjectId("620524134a84ecd9ac78f61f"),
        "name" : "Full-Access",
        "description" : "Full access to the Application",
        "methods" : [
            "GET",
            "POST",
            "PUT",
            "DELETE"
        ]
    },
]);

db.createCollection('applications');
db.applications.insertMany([
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
]);

db.createCollection('users');
db.users.insertMany([
    {
        "_id" : ObjectId("6243308399fd59c176c52fd4"),
        "name" : "admin",
        "email" : "root@dev.null",
        "password" : "$2a$11$lAT02Pq3MiHefYLYM6ZrUO79swRZAHeE0x0/RX13lIRouX72Hzwr2",
        "acl" : [
            {
                "application" : {
                    "_id" : ObjectId("6242d43e99fd59c176c52fd3")
                },
                "permissions" : {
                    "_id" : ObjectId("620524134a84ecd9ac78f61f")
                }
            },
            {
                "application" : {
                    "_id" : ObjectId("6245984799fd59c176c52fd5")
                },
                "permissions" : {
                    "_id" : ObjectId("620524134a84ecd9ac78f61f")
                }
            },
            {
                "application" : {
                    "_id" : ObjectId("6246d923ad35f14740a5fa79")
                },
                "permissions" : {
                    "_id" : ObjectId("620524134a84ecd9ac78f61f")
                }
            }
        ]
    }
])
