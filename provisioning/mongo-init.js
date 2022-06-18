/*
    Please, note, that _id fields are hardcoded. This because of using
    standard uniqe object id as a refernce. Following _ids are unique.
    So it's safe to use this provisioning.

    Users
        - admin:default_password

    @TODO add propper indices.
*/
db = db.getSiblingDB('v4');

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
        "_id" : ObjectId("62937983ec569fe63ccffebc"),
        "name" : "Roles management",
        "path" : "/roles"
    }
]);

db.createCollection('roles');
db.roles.insertMany([
    {
        "_id" : ObjectId("6242d43e99fd59c176c52fd4"),
        "name" : "System management",
        "description": "System management role",
        "apps": [
            ObjectId("6242d43e99fd59c176c52fd3"),
            ObjectId("6245984799fd59c176c52fd5"),
            ObjectId("62937983ec569fe63ccffebc")
        ]
    }
]);
db.roles.createIndex({ "name": 1 }, { unique: true });

db.createCollection('users');
db.users.insertMany([
    {
        "_id" : ObjectId("6243308399fd59c176c52fd4"),
        "username" : "admin",
        "email" : "root@dev.null",
        "password" : "$2a$11$lAT02Pq3MiHefYLYM6ZrUO79swRZAHeE0x0/RX13lIRouX72Hzwr2",
        "roles": [
            ObjectId("6242d43e99fd59c176c52fd4") 
        ]
    }
])
db.users.createIndex({ "username": 1 }, { unique: true });
db.users.createIndex({ "email": 1 }, { unique: true });
