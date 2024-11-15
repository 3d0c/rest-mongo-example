Lyre-Be-V4 Documentation
========================

## Contents
- [Build and Run](#build-and-run)
    - [Development environment](#development-environment)
- [API Specification](#api-specification)
    - [Authentication](#authentication)
        - [Login](#login)
        - [Logout](#logout)
        - [List sessions](#list-sessions)
        - [Logout specific user](#logout-specific-user)
    - [Manage Users](#manage-users)
        - [List users](#list-users)
        - [Show specific user](#show-specific-user)
        - [Create user](#create-user)
        - [Update user](#update-user)
        - [Delete user](#delete-user)
    - [Manage User](#manage-user)
        - [Get user](#get-user)
    - [Manage Applications](#manage-applications)
        - [List applications](#list-applications)
        - [Show specific application](#show-specific-application)
        - [Create application](#create-application)
        - [Update application](#update-application)
        - [Delete application](#delete-application)
    - [Manage permissions](#manage-permissions)
        - [List permissions](#list-permission)
        - [Show specific permission](#show-specific-permission)
        - [Create permission](#create-permission)
        - [Update permission](#update-permission)
        - [Delete permission](#delete-permission)
    - [Manage roles](#manage-roles)
        - [List roles](#list-roles)
        - [Show specific role](#show-specific-role)
        - [Create role](#create-role)
        - [Update role](#update-role)
        - [Delete role](#delete-role)
    - [Manage parameter](#manage-parameters)
        - [List parameters](#list-parameters)
        - [Show specific parameter](#show-specific-parameter)
        - [Create parameter](#create-parameter)
        - [Update parameter](#update-parameter)
        - [Delete parameter](#delete-parameter)
    - [Manage settings](#manage-settings)
        - [List settings](#list-settings)
        - [Show specific setting](#show-specific-setting)
        - [Create setting](#create-setting)
        - [Update setting](#update-setting)
        - [Delete setting](#delete-setting)        
- [Internals](#internals)
    - [ACL](#acl)
    - [Routing chain](#routing-chain)
- [Development](#development)
    - [Of using github](#of-using-github)
    - [Code style](#code-style)

# Build and Run

## Development environment

Prerequisites for this step is to have Docker installed on the system (Docker Desktop for OSX and Windows).

Run a single command inside project root directory:

```sh
docker-compose --verbose up -d
```

Check that it works:

```sh
curl -v -XPOST \
-H "Content-Type: application/json" \
-d '{"user_name":"admin","password":"default_password"}' \
localhost:8443/v1/sessions
```
This request should return `200 OK` and new token.

# API Specification

## Authentication

### Login

Getting session token, which must be provided with all others API calls as
header `Authorization: Bearer TOKEN`

Request:

```applescript
# Endpoint
POST /v1/sessions

# Expected content type
Content-Type: "application/json"

# Payload
{
    "user_name": (string)
    "password":  (string)
}
# or use email instead of user_name
{
    "email":    (string)
    "password": (string)
}
```

Response:

```applescript
Content-Type: "application/json"

# Expected status codes
200 OK
400 Bad request
403 Forbidden
404 Not found
503 Internal server error

# Body
{
    "token": (string)
}
```

Examples:

```applescript
# Happy pass
curl -v -XPOST -H "Content-Type: application/json" \
-d '{"user_name": "admin", "password": "default_password"}' \
localhost:8443/v1/sessions

# Response
< HTTP/1.1 200 OK
< Date: Sat, 26 Mar 2022 11:48:34 GMT
< Content-Length: 216
< Content-Type: text/plain; charset=utf-8
<
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJPYmplY3RJRChcIjYyMDY2YWVhNmQ0NzNmYmUwYWJmNjVmZFwiKSIsImV4cCI6MTY0ODI5ODkxNCwiaXNzIjoibHlyZS1iZS12NCJ9.O9j5_kcrseTN02ZrCKrEtow7tfPByW8RDfOn0MXP0vM"
}

```

```applescript
# User not found
curl -v -XPOST -H "Content-Type: application/json" \
-d '{"user_name": "nosuchuser", "password": "default_password"}' \
localhost:8443/v1/sessions

# Response
< HTTP/1.1 404 Not Found
< Date: Sat, 26 Mar 2022 11:53:46 GMT
< Content-Length: 0
```

```applescript
# Malformed request
curl -v -XPOST -H "Content-Type: application/json" -d \
'{"user_name": "nosuchuser"}' \
localhost:8443/v1/sessions

#Response
< HTTP/1.1 400 Bad Request
< Date: Sat, 26 Mar 2022 11:56:16 GMT
< Content-Length: 0
```

### Logout

Removes session token from `sessions` collection. So it's no more possible to use it. Please note, that this method doesn't check of existence of passed token in database, if token is valid it tries to remove it anyway. It's done for performance reasons, not to do extra database queries.

Request:

```applescript
# Endpoint
DELETE /v1/sessions

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
204 No content
400 Bad request
403 Forbidden
404 Not found
503 Internal server error
```

Please note, the `DELETE` method returns empty body. Only the status code.

Examples:

```applescript
# Request
curl -v -XDELETE -H \
"Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjA2NmFlYTZkNDczZmJlMGFiZjY1ZmQiLCJleHAiOjE2NDgzNzg0MjYsImlzcyI6Imx5cmUtYmUtdjQifQ.FVV0ZSTOCxJXJmh0hdxHdd61saoSPK9MANovhiEtvjQ" \
localhost:8443/v1/sessions

# Response
< HTTP/1.1 204 No Content
< Date: Sun, 27 Mar 2022 09:54:12 GMT
```

## Manage Users

### List Users

List all users. Returns an array of complete user models.

Request:

```applescript
# Endpoint
GET /v1/users

# Available parameters
?role=roleID

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Array of complete user models
```

Examples:

```applescript
# Request
curl -v -XGET \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjA2NmFlYTZkNDczZmJlMGFiZjY1ZmQiLCJleHAiOjE2NDg1NTAzMzIsImlzcyI6Imx5cmUtYmUtdjQifQ.uq9F7SEZ6ze1jVtcEvYNfJa-W7YLyF8TGEgxljx0BJk" \
localhost:8443/v1/users

# Response
< HTTP/1.1 200 OK
[
    {
        "id" : ObjectId("620527ed4a84ecd9ac78f623"),
        "name" : "admin",
        "email" : "root@dev.null",
        "avatar": "path/uuid.ext"
        "acl": [
            {
                "id": "6242d43e99fd59c176c52fd3",
                "name": "User management application",
                "path": "/users"
            },
            {
                "id": "6245984799fd59c176c52fd5",
                "name": "Applications management",
                "path": "/applications"
            },
            {
                "id": "6246d923ad35f14740a5fa79",
                "name": "Permissions management",
                "path": "/permissions"
            },
            {
                "id": "62937983ec569fe63ccffebc",
                "name": "Roles management",
                "path": "/roles"
            },
            {
                "id": "6242d43e99fd59c176c52fd3",
                "name": "User management application",
                "path": "/users"
            }
        ],
        {
                "application" : ObjectId("620527c04a84ecd9ac78f622"),
                "permission" : ObjectId("620524134a84ecd9ac78f61f"),
                "app_details" : {
                    "id" : ObjectId("620527c04a84ecd9ac78f622"),
                    "name" : "Permissions management API",
                    "path" : "/permissions"
                },
                "perm_details" : {
                    "id" : ObjectId("620524134a84ecd9ac78f61f"),
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
        "roles": [
            "6242d43e99fd59c176c52fd4"
        ],
    }
]
```

To get a list of users with specific role use:

```
curl -v -XGET \
-H "Authorization: Bearer \
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTU4MTg4ODYsImlzcyI6Imx5cmUtYmUtdjQifQ.2cLfysvWqtKftvshxyGtTPf2l_z-SmGDKHtahXfDBYc" \
"localhost:8443/v1/users?role=6242d43e99fd59c176c52fd4"
```

### Show specific user

Get user by ID. Returns single user object.

Request:

```applescript
# Endpoint
GET /v1/users/{id}

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Single user model
```

Examples:

```applescript
# Request
curl -v -XGET -H \
"Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjA2NmFlYTZkNDczZmJlMGFiZjY1ZmQiLCJleHAiOjE2NDg1NTQwNjAsImlzcyI6Imx5cmUtYmUtdjQifQ.WuE5zD7o9NCk7_M74OTTVrPND_vY8d78ZYLBkyt-IlY" \
localhost:8443/v1/users/623ec112c8c51a6a37ae839d

# Response
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Tue, 29 Mar 2022 10:56:09 GMT
< Content-Length: 1194
<
{
    "id": "623ec112c8c51a6a37ae839d",
    "user_name": "user1",
    "email": "user1@dev.null",
    "avatar": "path/uuid.ext"
    "acl": [
            {
                "id": "6242d43e99fd59c176c52fd3",
                "name": "User management application",
                "path": "/users"
            },
            {
                "id": "6245984799fd59c176c52fd5",
                "name": "Applications management",
                "path": "/applications"
            },
            {
                "id": "6246d923ad35f14740a5fa79",
                "name": "Permissions management",
                "path": "/permissions"
            },
            {
                "id": "62937983ec569fe63ccffebc",
                "name": "Roles management",
                "path": "/roles"
            },
            {
                "id": "6242d43e99fd59c176c52fd3",
                "name": "User management application",
                "path": "/users"
            }
        ],
        {
            "application": {
                "id": "620527c04a84ecd9ac78f622",
                "name": "Another One",
                "path": "/another"
            },
            "permissions": {
                "id": "620524134a84ecd9ac78f61f",
                "name": "Full Access",
                "description": "Full access to Application",
                "methods": [
                    "GET",
                    "POST",
                    "PUT",
                    "DELETE"
                ]
            }
        }
    ],
    "roles": [
        "6242d43e99fd59c176c52fd4"
    ],
}
```
### Create user

Create new user

Request:

```applescript
# Endpoint
POST /v1/users

# Expected content type
Content-Type: "application/json"

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
# User model in format as following:
{
    "user_name": "user1",
    "email": "user2@dev.null",
    "password": "plain_text_password",
    "avatar": "Image data URI scheme with base64 encoded data output"
    "roles": [
        "6242d43e99fd59c176c52fd4"
    ],    
}
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Complete single user model. See GET /users/{ID} for response example
```

Examples:

```applescript
# Request
curl -v -XPOST \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTExNjI1OTcsImlzcyI6Imx5cmUtYmUtdjQifQ.rKnp0ooe48ies83d5WhZmTCke_0Pi7p5EESbKovfXzY" \
-H 'Content-Type: application/json' \
-d '{
    "user_name": "user4",
    "email": "user4@dev.null",
    "password": "default",
    // "avatar": not provided for example because of length
    "roles": [
        "6242d43e99fd59c176c52fd4"
    ]    
}' localhost:8443/v1/users

# Response
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Tue, 29 Mar 2022 16:26:55 GMT
< Content-Length: 654
<
{
    "id": "6243334ff0326c4cf6986459",
    "user_name": "user4",
    "email": "user4@dev.null",
    "avatar": "path/uuid.ext"
    "acl": [
        {
            "id": "6242d43e99fd59c176c52fd3",
            "name": "User management application",
            "path": "/users"
        },
        {
            "id": "6245984799fd59c176c52fd5",
            "name": "Applications management",
            "path": "/applications"
        },
        {
            "id": "6246d923ad35f14740a5fa79",
            "name": "Permissions management",
            "path": "/permissions"
        },
        {
            "id": "62937983ec569fe63ccffebc",
            "name": "Roles management",
            "path": "/roles"
        },
        {
            "id": "6242d43e99fd59c176c52fd3",
            "name": "User management application",
            "path": "/users"
        }
    ],
    "roles": [
        "6242d43e99fd59c176c52fd4"
    ],    
}
```

### Update user

Please note, that because of MongoDB specific, this request actually replaces the whole document,
the only field could be omitted is the Password. Result is a single user object.

Request:

```applescript
# Endpoint
PUT /v1/users/{ID}

# Expected content type
Content-Type: "application/json"

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
# User model in format as following:
{
    "user_name": "user5555",
    "email": "user6666@dev.null",
    "password": "plain_text_password", # This one is optional
    "avatar": "Image data URI scheme with base64 encoded data output"
    "roles": [
        "6242d43e99fd59c176c52fd4"
    ],    
}
```
Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Bare user object. Password is optional
```

Examples:

```applescript
# Request
curl -v -XPUT \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTExNzczMzYsImlzcyI6Imx5cmUtYmUtdjQifQ.IE_e0z51K8STYfulVWCJpWky8nGOA3qVi416YQr1fhs" \
-H 'Content-Type: application/json' \
-d '{
    "user_name": "user6",
    "email": "user6@dev.null",
    "password": "default",
    "avatar": "path/uuid.ext"
    "roles": [
        "6242d43e99fd59c176c52fd4"
    ],    
}' localhost:8443/v1/users/62436b5ab97ea7529242bad6

# Response
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Tue, 29 Mar 2022 21:42:24 GMT
< Content-Length: 670
<
{
    "id": "62436b5ab97ea7529242bad6",
    "user_name": "user6",
    "email": "user6@dev.null",
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
    "roles": [
        "6242d43e99fd59c176c52fd4"
    ],    
}
```

### Remove user
Removes user from `users` collection.

Request:

```applescript
# Endpoint
DELETE /v1/users/{ID}

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
204 No content
400 Bad request
403 Forbidden
503 Internal server error
```

Please note, the `DELETE` method returns empty body. Only the status code.

Examples:

```applescript
# Request
curl -v -XDELETE \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTExNzczMzYsImlzcyI6Imx5cmUtYmUtdjQifQ.IE_e0z51K8STYfulVWCJpWky8nGOA3qVi416YQr1fhs" \
localhost:8443/v1/users/62436b5ab97ea7529242bad6

# Response
< HTTP/1.1 204 No Content
< Content-Type: application/json
< Date: Tue, 29 Mar 2022 22:11:51 GMT
```

### Update User Password
Updates the user password from `users` collection.

Request:

```applescript
# Endpoint
PUT /v1/users/{ID}/password

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
{
    "old_password": "default_password2", 
    "new_password": "default_password"
}
```

## Manage User

This group of endpoints are using for "self-management". Because user cannot get it's own profile without knowing it's identifier.

### Get user
Get user using ID from token. Returns complete user model or error.

Request:

```applescript
# Endpoint
GET /v1/user

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

```applescript
# Expected status codes
200 OK
403 Forbidden
503 Internal server error

# Body
Single user model
```

Examples:

```applescript
# Request
curl -v -XDELETE \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTExNzczMzYsImlzcyI6Imx5cmUtYmUtdjQifQ.IE_e0z51K8STYfulVWCJpWky8nGOA3qVi416YQr1fhs" \
localhost:8443/v1/user
```

## Manage Applications

### List Applications

List all applications. Returns an array of applications models.

Request:

```applescript
# Endpoint
GET /v1/applications

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
200 OK
403 Forbidden
503 Internal server error

# Body
Array of application models
```

Examples:

```applescript
# Request
curl -v -XGET \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjA2NmFlYTZkNDczZmJlMGFiZjY1ZmQiLCJleHAiOjE2NDg1NTAzMzIsImlzcyI6Imx5cmUtYmUtdjQifQ.uq9F7SEZ6ze1jVtcEvYNfJa-W7YLyF8TGEgxljx0BJk" \
localhost:8443/v1/applications

# Response
< Content-Type: application/json
< Date: Thu, 31 Mar 2022 12:11:58 GMT
< Content-Length: 483
<
[
    {
        "id": "620524994a84ecd9ac78f620",
        "name": "Sample One",
        "path": "/sample"
    },
    {
        "id": "620527c04a84ecd9ac78f622",
        "name": "Another One",
        "path": "/another"
    },
    {
        "id": "6242d43e99fd59c176c52fd3",
        "name": "User management application",
        "path": "/users"
    },
    {
        "id": "6245984799fd59c176c52fd5",
        "name": "Applications management",
        "path": "/applications"
    }
]
```

### Show specific application

Get application by ID. Returns single application object.

Request:

```applescript
# Endpoint
GET /v1/applications/{id}

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Single application model
```

Examples:

```applescript
# Request
curl -v -XGET \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTExNzczMzYsImlzcyI6Imx5cmUtYmUtdjQifQ.IE_e0z51K8STYfulVWCJpWky8nGOA3qVi416YQr1fhs" \
localhost:8443/v1/applications/6245984799fd59c176c52fd5

# Response
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Thu, 31 Mar 2022 13:14:36 GMT
< Content-Length: 108
<
{
    "id": "6245984799fd59c176c52fd5",
    "name": "Applications management",
    "path": "/applications"
}
```

### Create application

Create new application. Path should correspond routing path. It's going to be matched on request. See provisioning of how default applications (like /users) are created and [Routing Chain](#routing-chain) for details.  

Please note. Path MUST start with `/` symbol.

Request:

```applescript
# Endpoint
POST /v1/applications

# Expected content type
Content-Type: "application/json"

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
# Application model in format as following:
{
    "name": "Plant Logger",
    "path": "/plant_logger",
}
```

Examples:

```applescript
# Request
curl -v -XPOST \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTExNzczMzYsImlzcyI6Imx5cmUtYmUtdjQifQ.IE_e0z51K8STYfulVWCJpWky8nGOA3qVi416YQr1fhs" \
-H "Content-Type: application/json" \
-d '{"name": "Plant Logger", "path": "/pant_logger"}' \
localhost:8443/v1/applications

# Response
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Thu, 31 Mar 2022 20:03:24 GMT
< Content-Length: 95
<
{
    "id": "6246090c5f6ad1232cc8fb7a",
    "name": "Plant Logger",
    "path": "/pant_logger"
}
```

### Update application

Please note, that because of MongoDB specific, this request actually replaces the whole document.  
Also, if you change an application path, do not forget to update the [Routing Chain](#routing-chain). Result is a single application object.  

This is the Update - so, all references to this object will be preserved. 

Request:

```applescript
# Endpoint
PUT /v1/applications/{ID}

# Expected content type
Content-Type: "application/json"

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
# Application object in format as following:
{
    "name": "Another Application",
    "path": "/another_one"
}
```
Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Application object
```

Examples:

```applescript
# Request
curl -v -XPUT \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTExNzczMzYsImlzcyI6Imx5cmUtYmUtdjQifQ.IE_e0z51K8STYfulVWCJpWky8nGOA3qVi416YQr1fhs" \
-H "Content-Type: application/json" -d '{"name":"Plant Logger2","path":"pant_logger"}'  \
localhost:8443/v1/applications/6246090c5f6ad1232cc8fb7a

# Response
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Thu, 31 Mar 2022 21:31:01 GMT
< Content-Length: 96
<
{
    "id": "6246090c5f6ad1232cc8fb7a",
    "name": "Plant Logger2",
    "path": "pant_logger"
}
```

### Remove application
Removes application from `applications` collection.  

@TODO and @NOTE:  
_after deletion, references to the removed application still exist in `users` collection, but because it can't be resolved, this application is no longer available. Todo - cascading remove references from `users` collections inside transaction_.

Request:

```applescript
# Endpoint
DELETE /v1/applications/{ID}

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
204 No content
400 Bad request
403 Forbidden
503 Internal server error
```

Please note, the `DELETE` method returns empty body. Only the status code.

Examples:

```applescript
# Request
curl -v -XDELETE \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTExNzczMzYsImlzcyI6Imx5cmUtYmUtdjQifQ.IE_e0z51K8STYfulVWCJpWky8nGOA3qVi416YQr1fhs" \
-H "Content-Type: application/json" \
localhost:8443/v1/applications/6246090c5f6ad1232cc8fb7a

# Response
< HTTP/1.1 204 No Content
< Content-Type: application/json
< Date: Thu, 31 Mar 2022 21:53:43 GMT
```

## Manage Permissions

### List Permissions

List all permissions. Returns an array of permissions objects.

Request:

```applescript
# Endpoint
GET /v1/permissions

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
200 OK
403 Forbidden
503 Internal server error

# Body
Array of permissions objects
```

Examples:

```applescript
# Request
curl -v -XGET \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTExNzczMzYsImlzcyI6Imx5cmUtYmUtdjQifQ.IE_e0z51K8STYfulVWCJpWky8nGOA3qVi416YQr1fhs" \
localhost:8443/v1/permissions

# Response
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 01 Apr 2022 11:45:39 GMT
< Content-Length: 660
<
[
    {
        "id": "620524134a84ecd9ac78f61d",
        "name": "Read-Only",
        "description": "Read only access to Application",
        "methods": [
            "GET"
        ]
    },
    {
        "id": "620524134a84ecd9ac78f61e",
        "name": "Read-Write",
        "description": "Read and Create access to Application",
        "methods": [
            "GET",
            "POST"
        ]
    },
    {
        "id": "620524134a84ecd9ac78f61f",
        "name": "Full Access",
        "description": "Full access to Application",
        "methods": [
            "GET",
            "POST",
            "PUT",
            "DELETE"
        ]
    }
]

```

### Show specific permission

Get permission by ID. Returns single permission object.

Request:

```applescript
# Endpoint
GET /v1/permissions/{id}

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Single permission object
```

Examples:

```applescript
# Request
curl -v -XGET \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTExNzczMzYsImlzcyI6Imx5cmUtYmUtdjQifQ.IE_e0z51K8STYfulVWCJpWky8nGOA3qVi416YQr1fhs" \
localhost:8443/v1/permissions/620524134a84ecd9ac78f61f

# Response
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 01 Apr 2022 11:48:31 GMT
< Content-Length: 203
<
{
    "id": "620524134a84ecd9ac78f61f",
    "name": "Full Access",
    "description": "Full access to Application",
    "methods": [
        "GET",
        "POST",
        "PUT",
        "DELETE"
    ]
}
```

### Create permission

Create new permission.

Request:

```applescript
# Endpoint
POST /v1/permissions

# Expected content type
Content-Type: "application/json"

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
# Permission object in format as following:
# Allowed methods: `GET`, `POST`, `PUT`, `DELETE`
{
    "name": "Plant Logger",
    "methods": [],
}
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Single permission object
```

Examples:

```applescript
# Request
curl -v -XPOST \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTExNzczMzYsImlzcyI6Imx5cmUtYmUtdjQifQ.IE_e0z51K8STYfulVWCJpWky8nGOA3qVi416YQr1fhs" \
-H "Content-Type: application/json" \
-d '{"name":"test1", "methods":["GET"]}' \
localhost:8443/v1/permissions

< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 01 Apr 2022 12:41:42 GMT
< Content-Length: 122
<
{
    "id": "6246f30687ead2746d1340a2",
    "name": "test1",
    "description": "",
    "methods": [
        "GET"
    ]
}
```

### Update permission

Please note, that because of MongoDB specific, this request actually replaces the whole document.  

This is the Update - so, all references to this object will be preserved. 

Request:

```applescript
# Endpoint
PUT /v1/permissions/{ID}

# Expected content type
Content-Type: "application/json"

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
# Permission object in format as following:
# Allowed methods: `GET`, `POST`, `PUT`, `DELETE`
{
    "name": "Plant Logger",
    "methods": [],
}
```
Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Permission object
```

Examples:

```applescript
# Request
curl -v -XPUT \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTExNzczMzYsImlzcyI6Imx5cmUtYmUtdjQifQ.IE_e0z51K8STYfulVWCJpWky8nGOA3qVi416YQr1fhs" \
-H "Content-Type: application/json" \
-d '{"name":"test1", "methods":["GET","POST"]}' \
localhost:8443/v1/permissions/6246f30687ead2746d1340a2

# Response
< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Fri, 01 Apr 2022 12:54:39 GMT
< Content-Length: 138
<
{
    "id": "6246f30687ead2746d1340a2",
    "name": "test1",
    "description": "",
    "methods": [
        "GET",
        "POST"
    ]
}
```

### Remove permission
Removes permission document from `permissions` collection.  

@TODO and @NOTE:  
_after deletion, references to the removed permission still exist in `users` collection, but because it can't be resolved, this permission is no longer available. Todo - cascading remove references from `users` collections inside transaction_.

Request:

```applescript
# Endpoint
DELETE /v1/permissions/{ID}

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
204 No content
400 Bad request
403 Forbidden
503 Internal server error
```

Please note, the `DELETE` method returns empty body. Only the status code.

Examples:

```applescript
# Request
curl -v -XDELETE \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTExNzczMzYsImlzcyI6Imx5cmUtYmUtdjQifQ.IE_e0z51K8STYfulVWCJpWky8nGOA3qVi416YQr1fhs" \
localhost:8443/v1/permissions/6246f30687ead2746d1340a2

# Response
< HTTP/1.1 204 No Content
< Content-Type: application/json
< Date: Fri, 01 Apr 2022 13:02:58 GMT
```

## Manage Roles

### List Roles

List all roles. Returns an array of role objects.

Request:

```applescript
# Endpoint
GET /v1/roles

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
200 OK
403 Forbidden
503 Internal server error

# Body
Array of roles objects
```

Examples:

```applescript
# Request
curl -v -XGET \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTY0Mjk3ODMsImlzcyI6Imx5cmUtYmUtdjQifQ.JQnMKs0hy-GMoEf1Vh021grOefPJmSk649bBkBXN5-Y" \
localhost:8443/v1/roles

# Response
< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token
< Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
< Content-Type: application/json
< Date: Sun, 29 May 2022 15:17:50 GMT
< Content-Length: 330
<
[
    {
        "id": "6242d43e99fd59c176c52fd4",
        "name": "System management",
        "description": "System management role",
        "apps": [
            "6242d43e99fd59c176c52fd3",
            "6245984799fd59c176c52fd5",
            "6246d923ad35f14740a5fa79",
            "62937983ec569fe63ccffebc"
        ]
    }
]
```

### Show specific role

Get role by ID. Returns single role object.

Request:

```applescript
# Endpoint
GET /v1/role/{id}

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Single role object
```

Examples:

```applescript
# Request
curl -v -XGET \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTY0Mjk3ODMsImlzcyI6Imx5cmUtYmUtdjQifQ.JQnMKs0hy-GMoEf1Vh021grOefPJmSk649bBkBXN5-Y" \
localhost:8443/v1/roles/6242d43e99fd59c176c52fd4

# Response
< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token
< Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
< Content-Type: application/json
< Date: Sun, 29 May 2022 15:24:05 GMT
< Content-Length: 282
<
{
    "id": "6242d43e99fd59c176c52fd4",
    "name": "System management",
    "description": "System management role",
    "apps": [
        "6242d43e99fd59c176c52fd3",
        "6245984799fd59c176c52fd5",
        "6246d923ad35f14740a5fa79",
        "62937983ec569fe63ccffebc"
    ]
}
```

### Create role

Create new role.

Request:

```applescript
# Endpoint
POST /v1/roles

# Expected content type
Content-Type: "application/json"

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
# Role object
{
    "name": "Alison's role",
    "description": "Optional description",
    "apps": [Objectid1, Objectid2],
}
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Single role object
```

Examples:

```applescript
# Request
curl -v -XPOST \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTY0Mjk3ODMsImlzcyI6Imx5cmUtYmUtdjQifQ.JQnMKs0hy-GMoEf1Vh021grOefPJmSk649bBkBXN5-Y" \
-H "Content-Type: application/json" \
-d '{"name":"New role#2", "apps":["6246d923ad35f14740a5fa79"]}' \
localhost:8443/v1/roles

< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token
< Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
< Content-Type: application/json
< Date: Sun, 29 May 2022 15:33:23 GMT
< Content-Length: 145
<
{
    "id": "62939243ef10cbbc03b16c41",
    "name": "New role#2",
    "description": "",
    "apps": [
        "6246d923ad35f14740a5fa79"
    ]
}
```

### Update role

Please note, that because of MongoDB specific, this request actually replaces the whole document.  

This is the Update - so, all references to this object will be preserved. 

Request:

```applescript
# Endpoint
PUT /v1/roles/{ID}

# Expected content type
Content-Type: "application/json"

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
# Role object in format as following:
{
    "name": "My new role",
    "apps": [],
}
```
Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Permission object
```

Examples:

```applescript
# Request
curl -v -XPUT \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTY0Mjk3ODMsImlzcyI6Imx5cmUtYmUtdjQifQ.JQnMKs0hy-GMoEf1Vh021grOefPJmSk649bBkBXN5-Y" \
-H "Content-Type: application/json" \
-d '{"name":"test1", "apps":["62937983ec569fe63ccffebc"]}' \
localhost:8443/v1/roles/62939243ef10cbbc03b16c41

# Response
< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token
< Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
< Content-Type: application/json
< Date: Sun, 29 May 2022 15:38:26 GMT
< Content-Length: 140
<
{
    "id": "62939243ef10cbbc03b16c41",
    "name": "test1",
    "description": "",
    "apps": [
        "62937983ec569fe63ccffebc"
    ]
}
```

### Remove role
Removes role document from `roles` collection.  

@TODO and @NOTE:  
_after deletion, references to the removed roles still exist in `users` collection, but because it can't be resolved, this role is no longer available. Todo - cascading remove references from `users` collections inside transaction_.

Request:

```applescript
# Endpoint
DELETE /v1/roles/{ID}

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
204 No content
400 Bad request
403 Forbidden
503 Internal server error
```

Please note, the `DELETE` method returns empty body. Only the status code.

Examples:

```applescript
# Request
curl -v -XDELETE \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTY0Mjk3ODMsImlzcyI6Imx5cmUtYmUtdjQifQ.JQnMKs0hy-GMoEf1Vh021grOefPJmSk649bBkBXN5-Y" \
localhost:8443/v1/roles/62939243ef10cbbc03b16c41

# Response
< HTTP/1.1 204 No Content
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token
< Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
< Content-Type: application/json
< Date: Sun, 29 May 2022 15:42:18 GMT
```

## Manage Parameters

### List Parameters

List all parameters. Returns an array of parameter objects.

Request:

```applescript
# Endpoint
GET /v1/parameters

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
200 OK
403 Forbidden
503 Internal server error

# Body
Array of parameter objects
```

Examples:

```applescript
# Request
curl -v -XGET \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTY0Mjk3ODMsImlzcyI6Imx5cmUtYmUtdjQifQ.JQnMKs0hy-GMoEf1Vh021grOefPJmSk649bBkBXN5-Y" \
localhost:8443/v1/parameters

# Response
< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token
< Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
< Content-Type: application/json
< Date: Sun, 29 May 2022 15:17:50 GMT
< Content-Length: 330
<
[
    {
        "id": "62b827d15e329ab2f7f8b0e9",
        "name": "Sample Parameter 1",
        "description": "Sample Parameter 1",
        "comp_type": "select",
        "data_type": "string",
        "options": [
            "Option 1",
            "Option 2",
            "Option 3"
        ]
    },
    {
        "id": "62b828485e329ab2f7f8b0ea",
        "name": "Sample Parameter 2",
        "description": "Sample Parameter 2",
        "comp_type": "select",
        "data_type": "string",
        "options": [
            "Option 1",
            "Option 2",
            "Option 3"
        ]
    }
]
```

### Show specific parameter

Get parameter by ID. Returns single parameter object.

Request:

```applescript
# Endpoint
GET /v1/parameters/{id}

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Single role object
```

Examples:

```applescript
# Request
curl -v -XGET \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTk0NjYxMjQsImlzcyI6Imx5cmUtYmUtdjQifQ.Lif2940XSSRleQ2AA0qLLNYFt9ZGxyVeipOPrx7WAIc" \
http://localhost:8443/v1/parameters/62b827d15e329ab2f7f8b0e9

# Response
< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token
< Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
< Content-Type: application/json
< Date: Sun, 29 May 2022 15:24:05 GMT
< Content-Length: 282
<
{
    "id": "62b827d15e329ab2f7f8b0e9",
    "name": "Sample Parameter 1",
    "description": "Sample Parameter 1",
    "comp_type": "select",
    "data_type": "string",
    "options": [
        "01",
        "02",
        "03"
    ]
}
```

### Create parameter

Create new parameter.

Request:

```applescript
# Endpoint
POST /v1/parameters

# Expected content type
Content-Type: "application/json"

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
# Parameter object
{
    "name": "New Parameter",
    "description": "Description of new one",
    "comp_type": "select",
    "data_type": "string",
    "options": [
        "01",
        "02",
        "03"
    ]
}
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Single parameter object
```

Examples:

```applescript
# Request
curl -v -XPOST \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTk0NjYxMjQsImlzcyI6Imx5cmUtYmUtdjQifQ.Lif2940XSSRleQ2AA0qLLNYFt9ZGxyVeipOPrx7WAIc" \
-H "Content-Type: application/json" \
-d '{"name":"New Parameter", "comp_type": "select", "data_type": "string", "options":["001"]}' \
localhost:8443/v1/parameters

< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization
< Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
< Content-Type: application/json
< Date: Mon, 04 Jul 2022 20:22:30 GMT
< Content-Length: 184
< 
{
    "id": "62c34c069fa2a6650a510e6f",
    "name": "New Parameter",
    "description": "",
    "comp_type": "select",
    "data_type": "string",
    "options": [
        "001"
    ]
}
```

### Update parameter

Please note, that because of MongoDB specific, this request actually replaces the whole document.  

This is the Update - so, all references to this object will be preserved. 

Request:

```applescript
# Endpoint
PUT /v1/parameters/{ID}

# Expected content type
Content-Type: "application/json"

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
# Parameter object
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Parameter object
```

Examples:

```applescript
# Request
curl -v -XPUT \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiI2MjQzMzA4Mzk5ZmQ1OWMxNzZjNTJmZDQiLCJleHAiOjE2NTk0NjYxMjQsImlzcyI6Imx5cmUtYmUtdjQifQ.Lif2940XSSRleQ2AA0qLLNYFt9ZGxyVeipOPrx7WAIc" \
-H "Content-Type: application/json" \
-d '{"name":"New Parameter Changed", "comp_type": "select", "data_type": "string", "options":["002"]}' \
localhost:8443/v1/parameters/62c34c069fa2a6650a510e6f

# Response
< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Headers: Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token
< Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE
< Content-Type: application/json
< Date: Sun, 29 May 2022 15:38:26 GMT
< Content-Length: 140
<
{
    "id": "62c34c069fa2a6650a510e6f",
    "name": "New Parameter Changed",
    "description": "",
    "comp_type": "select",
    "data_type": "string",
    "options": [
        "002"
    ]
}
```

### Remove parameter
Removes parameter document from `parameters` collection.  

Request:

```applescript
# Endpoint
DELETE /v1/parameters/{ID}

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
204 No content
400 Bad request
403 Forbidden
503 Internal server error
```

Please note, the `DELETE` method returns empty body. Only the status code.


## Manage Settings
Settings is a relation between user, application and parameters for this application. 
To assign settings to user, create parameters, create setting document and add setting id to user's settings array.

Structure of settings collections:

```javascript
{
        "_id" : ObjectId("62c1c1c95e329ab2f7f8b0ee"),
        "user_id" : ObjectId("6243308399fd59c176c52fd4"),
        "app_id" : ObjectId("6242d43e99fd59c176c52fd3"),
        "parameters" : [
                ObjectId("62b827d15e329ab2f7f8b0e9"),
                ObjectId("62b828485e329ab2f7f8b0ea")
        ]
}
```

### List Settings

List all settings. Returns an array of settings objects.

Request:

```applescript
# Endpoint
GET /v1/settings

# Parameters
user_id
app_id
user_id, app_id

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Array of setting objects
```

### Show specific parameter

Get setting by ID. Returns single setting object.

Request:

```applescript
# Endpoint
GET /v1/settings/{id}

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Single setting object
```

### Create setting

Create new setting object.

Request:

```applescript
# Endpoint
POST /v1/setting

# Expected content type
Content-Type: "application/json"

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
# Setting object
{
    "user_id" : "6243308399fd59c176c52fd4",
    "app_id" : "6242d43e99fd59c176c52fd3",
    "parameters" : [
            "62b827d15e329ab2f7f8b0e9",
            "62b828485e329ab2f7f8b0ea"
    ]
}
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Single parameter object
```

### Update setting

Please note, that because of MongoDB specific, this request actually replaces the whole document.  

This is the Update - so, all references to this object will be preserved. 

Request:

```applescript
# Endpoint
PUT /v1/settings/{ID}

# Expected content type
Content-Type: "application/json"

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
# Setting object
```

Response:

```applescript
# Expected status codes
200 OK
400 Bad request
403 Forbidden
503 Internal server error

# Body
Setting object
```

### Remove setting
Removes setting document from `settings` collection.  

Request:

```applescript
# Endpoint
DELETE /v1/settings/{ID}

# Expected authentication header
Authorization: Bearer TOKEN

# Payload
No payload required for this request
```

Response:

```applescript
# Expected status codes
204 No content
400 Bad request
403 Forbidden
503 Internal server error
```

Please note, the `DELETE` method returns empty body. Only the status code.

# Internals
## ACL
### Complete scheme of how does ACL work

Core objects of ACL:

- users
- applications
- permissions

_Please note, that these objects are the MongoDB collection names._

#### `users`

Current implementation of `user` object has three mandatory fields (in terms of `Golang` structure)

- Name
- Email
- Password

*ACL*  
Please note, that `acl` is a meta field. It means, that it isn't stored inside DataBase but expanded on the fly based on `roles` ID's.


Bare `users` element looks like

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
}
```

#### `applications`

Is a definition of any particular application, which API we want to expose.  
Mandatory fields:

- name
- path

### Example

Suppose there is a logged in user, which has valid token. On each request we're getting a complete `user` model, which could be represented as a JSON like:

```javascript
{
    "id": "6243308399fd59c176c52fd4",
    "user_name": "",
    "email": "root@dev.null",
    "acl": [
        {
            "id": "6242d43e99fd59c176c52fd3",
            "name": "User management application",
            "path": "/users"
        },
        {
            "id": "6245984799fd59c176c52fd5",
            "name": "Applications management",
            "path": "/applications"
        },
        {
            "id": "6246d923ad35f14740a5fa79",
            "name": "Permissions management",
            "path": "/permissions"
        },
        {
            "id": "62937983ec569fe63ccffebc",
            "name": "Roles management",
            "path": "/roles"
        },
        {
            "id": "6242d43e99fd59c176c52fd3",
            "name": "User management application",
            "path": "/users"
        }
    ],
    "roles": [
        "6242d43e99fd59c176c52fd4",
        "62a4a90e5d2fe1005407cd8b"
    ],
    "avatar": "",
    "first_name": "",
    "last_name": "",
    "plant_id": "",
    "printer": "",
    "last_login": "0001-01-01T00:00:00Z",
    "created_by": "",
    "created_date": "0001-01-01T00:00:00Z",
    "updated_by": "",
    "updated_date": "0001-01-01T00:00:00Z",
    "status": 0,
    "sap_id": "",
    "sap_pwd": ""
}

```

At the middleware level, before we accept request in the controller, we can validate is it possible at all. Validation is based on `acl[x].path`.

As an example application got a request like (assuming user is authorised and we've got a complete user model)

__First one:__

```
GET localhost:8443/api/v1/sample/users
                          ^^^^^^
                         acl.path
```

As a result it's allowed to do this method.

__Second one:__

```
POST localhost:8443/api/v1/another/documents
                          ^^^^^^
                         acl.path
```

As a result user will get `403 Forbidden` because there is no such method in allowed list.

## Routing chain

Regarding ACL implementation adding new routes should conform following rule:

- `IsAuthorized` Checks whether request has a valid token

- `GetUser` Try to get complete user information

- `IsPermit` Matches permissions with request


Example for some application:

```go
    r.Get(
        filepath.Join(root, "/myapplication"),
        middlewares.Chain(
            middlewares.IsAuthorized,
            middlewares.GetUser,
            appHandler().get,
        ),
    )
```

Example for main application, removing current session (logout):

```go
    r.Delete(
        filepath.Join(root, "/sessions"),
        middlewares.Chain(
            middlewares.IsAuthorized,
            middlewares.GetUser,
            // Take a look, that there is no IsPermit middleware
            sessionsHandler().remove,
        ),
    )
```

## Development

- DO NOT push directly to master branch
- `make lint` before commit
- Create separated branch and do pull Request
    - `git pull`
    - `git checkout -b feature/new-one`
    - do something
    - `git add . && git commit -am 'My new feature'`
    - `git push origin feature/new-one`
    - Do a code review than merge


# TODO

- CI/CD pipeline
- Integration (e2e) test

TODO @3d0c

- rename *m to m, *id to id(string) [unify] // @refactoring
- codegen?
- human readable errors // @refactoring [protocol]
