package handlers

import (
	"fmt"
	"net/http"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
)

type user struct {
	// *models.UserScheme @TODO
}

func userHandler() *user {
	return &user{}
}

func (u *user) get(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		user *models.UserScheme
	)

	if user = r.Context().Value(models.UserSchemeType{}).(*models.UserScheme); user == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error casting to *models.UserScheme")
	}

	return user, http.StatusOK, nil
}
