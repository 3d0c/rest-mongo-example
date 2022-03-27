# ACL

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
- ACLs

ACLs-field is a meta field. Which means it stores only referenced IDs in database.

ACL is an array of ACL structures, which has following fields

- Application : ObjectID
- Permission  : ObjectID

This is a references to corresponding collection entities.

#### `applications`

Is a definition of any particular application, which API we want to expose.  
Mandatory fields:

- name
- path

#### `permissions`

Permission is what particular user allowed to do with particular application. There are three default permissions:

- `Read-Only` Allowed API method is `GET`
- `Read-Write` Allowed API methods are `GET`, `POST`
- `Full Access` Allowed API methods are `GET`, `POST`, `PUT`, `DELETE`

To create these permissions one can use `fixtures/mongo.default`

### Example

Suppose there is a logged in user, which has valid token. On each request we're getting a complete `user` model, which could be represented as a JSON like:

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
				]
			}
		}
	]
}
```

At the middleware level, before we accept request in the controller, we can validate is it possible at all. Here we have two important things: Application Path `app_details.path` and permissions list for this application `perm_details.methods`.

As an example application got a request like (assuming user is authorised and we've got a complete user model)

__First one:__

```
GET localhost:8443/api/v1/sample/users
                          ^^^^^^
                      app_details.path
```

As a result it's allowed to do this method.

__Second one:__

```
POST localhost:8443/api/v1/another/documents
                          ^^^^^^
                      app_details.path
```

As a result user will get `403 Forbidden` because there is no such method in allowed list.

# Routing chain

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
			middlewares.IsPermit,
			logHandler().get,
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

# API Specification

### 1. Create a session (login)

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

### 2. Remove the session (logout)

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



# TODO

- Database initialisation script/app
- Dockerfile
- CI/CD pipline
- Integration (e2e) test
