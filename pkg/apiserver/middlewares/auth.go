package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/teal-seagull/lyre-be-v4/pkg/helpers"
)

const (
	authHeaderKey      = "Authorization"
	bearerSchemePrefix = "Bearer "
)

// IsAuthorized is a middleware which checks whether request contains
// an authorization token, checks it's validity and passes AuthClaim to
// the request context.
func IsAuthorized(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		authHeader string
		err        error
		claims     *helpers.AuthClaims
	)

	if authHeader = r.Header.Get(authHeaderKey); len(authHeader) < 8 {
		return nil, http.StatusUnauthorized, fmt.Errorf("wrong authHeader size")
	}

	tokenString := authHeader[len(bearerSchemePrefix):]

	if claims, err = helpers.VerifyToken(tokenString); err != nil {
		return nil, http.StatusUnauthorized, err
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, helpers.AuthClaims{}, claims)
	ctx = context.WithValue(ctx, helpers.TokenStringType{}, tokenString)

	*r = *r.WithContext(ctx)

	return nil, http.StatusOK, nil
}
