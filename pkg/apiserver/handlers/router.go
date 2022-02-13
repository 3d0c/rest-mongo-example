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

	// Auth routes, available for all
	//
	// Login = create session, former POST /auth/login
	r.Post(
		filepath.Join(root, "/sessions"),
		middlewares.Chain(
			sessionsHandler().create,
		),
	)

	// Logout = delete session, former POST /auth/logout
	r.Delete(
		filepath.Join(root, "/sessions"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			middlewares.IsPermit,
			sessionsHandler().remove,
		),
	)

	return r
}
