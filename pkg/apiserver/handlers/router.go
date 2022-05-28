package handlers

import (
	"path/filepath"

	"github.com/go-chi/chi"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/middlewares"
	"github.com/teal-seagull/lyre-be-v4/pkg/config"
)

// SetupRouter sets up endpoints
func SetupRouter(cfg config.Server) *chi.Mux {
	var (
		root string = "/" + cfg.APIVersion
	)

	r := chi.NewRouter()

	// Preflight OPTIONS, needed for CORS
	r.Options("/*", middlewares.Chain(nilHandler))

	// Authentication
	// Login. Available for all
	r.Post(
		filepath.Join(root, "/sessions"),
		middlewares.Chain(
			middlewares.IsValidContentType,
			sessionsHandler().create,
		),
	)
	// Logout
	r.Delete(
		filepath.Join(root, "/sessions"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			sessionsHandler().remove,
		),
	)

	// Users management
	// List all users in the system
	r.Get(
		filepath.Join(root, "/users"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			usersHandler().get,
		),
	)
	// Get specific user
	r.Get(
		filepath.Join(root, "/users/{ID}"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			usersHandler().getByID,
		),
	)
	// Create user
	r.Post(
		filepath.Join(root, "/users"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			usersHandler().create,
		),
	)
	// Update user
	r.Put(
		filepath.Join(root, "/users/{ID}"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			usersHandler().update,
		),
	)
	// Remove user
	r.Delete(
		filepath.Join(root, "/users/{ID}"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			usersHandler().remove,
		),
	)

	// Applications management
	// List all applications
	r.Get(
		filepath.Join(root, "/applications"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			appsHandler().get,
		),
	)
	// Get specific application
	r.Get(
		filepath.Join(root, "/applications/{ID}"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			appsHandler().getByID,
		),
	)
	// Create application
	r.Post(
		filepath.Join(root, "/applications"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			appsHandler().create,
		),
	)
	// Update applications
	r.Put(
		filepath.Join(root, "/applications/{ID}"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			appsHandler().update,
		),
	)
	// Remove application
	r.Delete(
		filepath.Join(root, "/applications/{ID}"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			appsHandler().remove,
		),
	)

	// Permissions management
	// List all permissions
	r.Get(
		filepath.Join(root, "/permissions"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			permHandler().get,
		),
	)
	// Get specific application
	r.Get(
		filepath.Join(root, "/permissions/{ID}"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			permHandler().getByID,
		),
	)
	// Create permission
	r.Post(
		filepath.Join(root, "/permissions"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			permHandler().create,
		),
	)
	// Update permission
	r.Put(
		filepath.Join(root, "/permissions/{ID}"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			permHandler().update,
		),
	)
	// Remove application
	r.Delete(
		filepath.Join(root, "/permissions/{ID}"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			permHandler().remove,
		),
	)

	return r
}
