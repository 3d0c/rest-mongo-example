package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
)

type users struct {
	// *models.UserScheme
}

func usersHandler() *users {
	return &users{}
}

func (u *users) get(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		um     *models.User
		result []models.UserScheme
		err    error
	)

	if um, err = models.NewUser(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing user model - %s", err)
	}

	if result, err = um.FindAll(chi.URLParam(r, "role")); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting users list - %s", err)
	}

	return result, http.StatusOK, nil
}

func (u *users) getByID(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		uid    string = chi.URLParam(r, "ID")
		um     *models.User
		result *models.UserScheme
		err    error
	)

	if uid == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("error getting user id")
	}

	if um, err = models.NewUser(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing user model - %s", err)
	}

	if result, err = um.FindByID(uid); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting user '%s' - %s", uid, err)
	}

	return result, http.StatusOK, nil
}

func (u *users) create(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		request *models.UserScheme = &models.UserScheme{}
		current *models.UserScheme
		result  *models.UserScheme
		um      *models.User
		uid     string
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error binding input data - %s", err)
	}

	if current = r.Context().Value(models.UserSchemeType{}).(*models.UserScheme); current == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing current user")
	}

	if um, err = models.NewUser(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing user model - %s", err)
	}

	request.CreatedDate = time.Now().UTC()
	request.CreatedBy = current.Name
	request.UpdatedDate = time.Now().UTC()
	request.UpdatedBy = current.Name

	if uid, err = um.Create(request); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error creating user - %s", err)
	}

	if result, err = um.FindByID(uid); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting user '%s' - %s", uid, err)
	}

	return result, http.StatusOK, nil
}

func (u *users) update(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		request *models.UserScheme = &models.UserScheme{}
		current *models.UserScheme
		result  *models.UserScheme
		um      *models.User
		uid     string = chi.URLParam(r, "ID")
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error binding input data - %s", err)
	}

	if current = r.Context().Value(models.UserSchemeType{}).(*models.UserScheme); current == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing current user")
	}

	if um, err = models.NewUser(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing user model - %s", err)
	}

	request.UpdatedDate = time.Now().UTC()
	request.UpdatedBy = current.Name

	if err = um.Update(uid, request); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error updating user '%s' - %s", uid, err)
	}

	if result, err = um.FindByID(uid); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting user '%s' - %s", uid, err)
	}

	return result, http.StatusOK, nil
}

func (u *users) remove(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		um  *models.User
		uid string = chi.URLParam(r, "ID")
		err error
	)

	if uid == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("error getting user id")
	}

	if um, err = models.NewUser(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing user model - %s", err)
	}

	if err = um.Delete(uid); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error deleting user model - %s", err)
	}

	return nil, http.StatusNoContent, nil
}
