package middlewares

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/teal-seagull/lyre-be-v4/pkg/helpers"
	"github.com/teal-seagull/lyre-be-v4/pkg/log"
)

// Middlewares type
type Middlewares func(res http.ResponseWriter, request *http.Request) (interface{}, int, error)

func Chain(m ...Middlewares) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			err    error
			status int
			result interface{}
		)

		for _, middleware := range m {
			if result, status, err = middleware(w, r); err != nil {
				break
			}
		}

		// Skip output for getfile middleware
		if w.Header().Get("Content-Disposition") != "" {
			return
		}

		// CORS
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(status)

		if err != nil {
			log.TheLogger().Error("HTTP Request",
				zap.Error(err), zap.String("method", r.Method), zap.String("source", r.RemoteAddr), zap.String("URI", r.RequestURI))
			return
		}

		if result != nil {
			helpers.NewJsonResponder(w).Write(result)
		}
	}
}
