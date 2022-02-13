package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
	"github.com/teal-seagull/lyre-be-v4/pkg/helpers"
)

// GetUser middleware should be second after IsAuthorized
//
// It creates complete user model and passes it into request context.
// So it's available in any middleware and handler
func GetUser(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		err       error
		claims    *helpers.AuthClaims
		userModel *models.User
		user      *models.UserScheme
	)

	if userModel, err = models.NewUser(); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if claims = r.Context().Value(helpers.AuthClaims{}).(*helpers.AuthClaims); claims == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error casting to *helpers.AuthClaim")
	}

	if user, err = userModel.Find(claims.GetUserID()); err != nil {
		return nil, http.StatusUnauthorized, fmt.Errorf("no such user found - %s", err)
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, models.UserSchemeType{}, user)

	*r = *r.WithContext(ctx)

	return nil, http.StatusOK, nil
}
