# ACL

### Complete scheme of how does ACL works

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
- `Read-Write` Allowed API method are `GET`, `POST`
- `Full Access` Allowed API methods are `GET`, `POST`, `PUT`, `DELETE`

To create these permissions one can use `fixtures/mongo.default`

### Example

Suppose there is a logged in user, which has valid token. On each request we're getting a complete `user` model, which could be represented as a JSON like:

```json
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


# TODO

- Unit tests for models using `mongo-mock`
- Database initialisation script/app
- Dockerfile
- CI/CD pipline
- Integration (e2e) test
