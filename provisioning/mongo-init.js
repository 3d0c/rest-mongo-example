/*
    Please, note, that _id fields are hardcoded. This because of using
    standard uniqe object id as a refernce. Following _ids are unique.
    So it's safe to use this provisioning.

    Users
        - admin:default_password

    @TODO add propper indices.
*/
db = db.getSiblingDB('v4');

db.applications.drop();
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
    },
    {
        "_id" : ObjectId("62c1e55b5e329ab2f7f8b0ef"),
        "name" : "Parameters management",
        "path" : "/parameters"
    },
    {
        "_id" : ObjectId("62c1e55d5e329ab2f7f8b0f0"),
        "name" : "Settings management",
        "path" : "/settings"
    },
    {
        "_id" : ObjectId("62eebee180e2b96494b6bf95"),
        "name" : "Document Viewer",
        "path" : "/docview"
    }    
]);

db.roles.drop();
db.createCollection('roles');
db.roles.insertMany([
    {
        "_id" : ObjectId("6242d43e99fd59c176c52fd4"),
        "name" : "System management",
        "description": "System management role",
        "apps": [
            ObjectId("6242d43e99fd59c176c52fd3"),
            ObjectId("6245984799fd59c176c52fd5"),
            ObjectId("62937983ec569fe63ccffebc"),
            ObjectId("62c1e55b5e329ab2f7f8b0ef"),
            ObjectId("62c1e55d5e329ab2f7f8b0f0"),
            ObjectId("62eebee180e2b96494b6bf95")
        ]
    }
]);
db.roles.createIndex({ "name": 1 }, { unique: true });

db.settings.drop()
db.createCollection('settings');
db.settings.insertMany([
    {
        "_id": ObjectId("62c1c1c95e329ab2f7f8b0ee"),
        "user_id": ObjectId("6243308399fd59c176c52fd4"),
        "app_id": ObjectId("6242d43e99fd59c176c52fd3"),
        "parameters": [
            ObjectId("62b827d15e329ab2f7f8b0e9"),
            ObjectId("62b828485e329ab2f7f8b0ea")
        ]
    }
])
db.users.createIndex({ "user_id": 1, "app_id": 1 }, { unique: true });
db.users.createIndex({ "user_id": 1 });
db.users.createIndex({ "app_id": 1 });

db.users.drop()
db.createCollection('users');
db.users.insertMany([
    {
        "_id" : ObjectId("6243308399fd59c176c52fd4"),
        "username" : "admin",
        "email" : "root@dev.null",
        "password" : "$2a$11$lAT02Pq3MiHefYLYM6ZrUO79swRZAHeE0x0/RX13lIRouX72Hzwr2",
        "roles": [
            ObjectId("6242d43e99fd59c176c52fd4") 
        ],
        "settings": [
            ObjectId("62c1c1c95e329ab2f7f8b0ee")
        ]
    }
])
db.users.createIndex({ "username": 1 }, { unique: true });
db.users.createIndex({ "email": 1 }, { unique: true });
