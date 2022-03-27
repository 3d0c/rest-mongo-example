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

	// Auth route, available for all
	r.Post(
		filepath.Join(root, "/sessions"),
		middlewares.Chain(
			middlewares.IsValidContentType,
			sessionsHandler().create,
		),
	)

	// Auth route, available for logged in user, so it should have proper token
	r.Delete(
		filepath.Join(root, "/sessions"),
		middlewares.Chain(
			middlewares.IsAuthorized,
			middlewares.GetUser,
			sessionsHandler().remove,
		),
	)

	return r
}
