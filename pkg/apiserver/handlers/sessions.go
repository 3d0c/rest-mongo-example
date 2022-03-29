package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
	"github.com/teal-seagull/lyre-be-v4/pkg/helpers"
)

type sessions struct{}

func sessionsHandler() *sessions {
	return &sessions{}
}

func (s *sessions) create(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		// request is used only for binding user name or email and password
		request *models.UserScheme = &models.UserScheme{}
		us      *models.UserScheme
		sm      *models.Session
		ss      *models.SessionScheme
		u       *models.User
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error binding input data - %s", err)
	}

	if u, err = models.NewUser(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing user model - %s", err)
	}

	if us, err = u.FindByName(request); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error finding user by name '%s', email '%s' - %s", request.Name, request.Email, err)
	}

	if !helpers.CompareHashWithPasswords(*us.Password, *request.Password) {
		return nil, http.StatusUnauthorized, fmt.Errorf("error comparing password for user '%s', email '%s'", request.Name, request.Email)
	}

	if sm, err = models.NewSession(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing session model - %s", err)
	}

	if ss, err = sm.Create(us.ID); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error creating session - %s", err)
	}

	return ss, http.StatusOK, nil
}

func (s *sessions) remove(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		sm    *models.Session
		token string
		err   error
	)

	if token = r.Context().Value(helpers.TokenStringType{}).(string); token == "" {
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting token from context")
	}

	if sm, err = models.NewSession(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing session model - %s", err)
	}

	if err = sm.Remove(token); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error removing session token - %s", err)
	}

	return nil, http.StatusNoContent, nil
}
