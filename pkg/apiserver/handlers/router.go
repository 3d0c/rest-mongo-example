package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/middlewares"
	"github.com/teal-seagull/lyre-be-v4/pkg/config"
	"github.com/teal-seagull/lyre-be-v4/pkg/log"
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
		root+"/sessions",
		middlewares.Chain(
			middlewares.IsValidContentType,
			sessionsHandler().create,
		),
	)
	// Logout
	r.Delete(
		root+"/sessions",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			sessionsHandler().remove,
		),
	)

	// Users management
	// List all users in the system
	r.Get(
		root+"/users",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			usersHandler().get,
		),
	)
	// Get specific user
	r.Get(
		root+"/users/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			usersHandler().getByID,
		),
	)
	// Create user
	r.Post(
		root+"/users",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			usersHandler().create,
		),
	)
	// Update user
	r.Put(
		root+"/users/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			usersHandler().update,
		),
	)
	// Update user password
	r.Put(
		root+"/users/{ID}/password",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			usersHandler().updatePassword,
		),
	)
	// Remove user
	r.Delete(
		root+"/users/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			usersHandler().remove,
		),
	)

	// User controller. Used to get user from token
	// Get current user
	r.Get(
		root+"/user",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			userHandler().get,
		),
	)

	// Applications management
	// List all applications
	r.Get(
		root+"/applications",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			appsHandler().get,
		),
	)
	// Get specific application
	r.Get(
		root+"/applications/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			appsHandler().getByID,
		),
	)
	// Create application
	r.Post(
		root+"/applications",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			appsHandler().create,
		),
	)
	// Update applications
	r.Put(
		root+"/applications/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			appsHandler().update,
		),
	)
	// Remove application
	r.Delete(
		root+"/applications/{ID}",
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
		root+"/permissions",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			permHandler().get,
		),
	)
	// Get specific application
	r.Get(
		root+"/permissions/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			permHandler().getByID,
		),
	)
	// Create permission
	r.Post(
		root+"/permissions",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			permHandler().create,
		),
	)
	// Update permission
	r.Put(
		root+"/permissions/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			permHandler().update,
		),
	)
	// Remove application
	r.Delete(
		root+"/permissions/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			permHandler().remove,
		),
	)

	// Roles management
	// List all roles
	r.Get(
		root+"/roles",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			rolesHandler().get,
		),
	)
	// Get specific role
	r.Get(
		root+"/roles/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			rolesHandler().getByID,
		),
	)
	// Create role
	r.Post(
		root+"/roles",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			rolesHandler().create,
		),
	)
	// Update role
	r.Put(
		root+"/roles/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			rolesHandler().update,
		),
	)
	// Remove role
	r.Delete(
		root+"/roles/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			rolesHandler().remove,
		),
	)

	// Parameters management
	// List all parameters
	r.Get(
		root+"/parameters",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			parametersHandler().get,
		),
	)
	// Get specific parameter
	r.Get(
		root+"/parameters/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			parametersHandler().getByID,
		),
	)
	// Create parameter
	r.Post(
		root+"/parameters",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			parametersHandler().create,
		),
	)
	// Update parameter
	r.Put(
		root+"/parameters/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			parametersHandler().update,
		),
	)
	// Remove role
	r.Delete(
		root+"/parameters/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			parametersHandler().remove,
		),
	)

	// Settings management
	// List all settings
	// @TODO should be available for admin user only
	// This endpoint accepts following parameters:
	// ?user_id=
	// ?app_id=
	// ?user_id=&app_id
	r.Get(
		root+"/settings",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			settingsHandler().get,
		),
	)
	// Get specific setting
	r.Get(
		root+"/settings/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			settingsHandler().getByID,
		),
	)
	// Create setting
	r.Post(
		root+"/settings",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			settingsHandler().create,
		),
	)
	// Update setting
	r.Put(
		root+"/settings/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			settingsHandler().update,
		),
	)
	// Remove setting
	r.Delete(
		root+"/setting/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			settingsHandler().remove,
		),
	)

	// Document view
	// Get documents list
	r.Get(
		root+"/docview",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			docviewHandler().get,
		),
	)
	// Get document's file
	r.Get(
		root+"/docview/{ID}",
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			docviewHandler().getFile,
		),
	)

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.TheLogger().Debug("registered", zap.String("method", method), zap.String("route", route))
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		log.TheLogger().Debug("logging error", zap.Error(err))
	}

	return r
}
