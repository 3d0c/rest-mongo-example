package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
)

type roles struct {
	// *models.UserScheme // @TODO init current user from constructor
}

func rolesHandler() *roles {
	return &roles{}
}

func (rl *roles) get(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		m      *models.Role
		result []models.RoleScheme
		err    error
	)

	if m, err = models.NewRole(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing role model - %s", err)
	}

	if result, err = m.FindAll(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting roles list - %s", err)
	}

	return result, http.StatusOK, nil
}

func (rl *roles) getByID(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		id     string = chi.URLParam(r, "ID")
		m      *models.Role
		result *models.RoleScheme
		err    error
	)

	if id == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("error getting role id")
	}

	if m, err = models.NewRole(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing role model - %s", err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting role '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}

func (rl *roles) create(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		request *models.RoleScheme = &models.RoleScheme{}
		result  *models.RoleScheme
		m       *models.Role
		id      string
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error binding input data - %s", err)
	}

	if m, err = models.NewRole(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing role model - %s", err)
	}

	if id, err = m.Create(request); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error creating role - %s", err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting role object '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}

func (rl *roles) update(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		request *models.RoleScheme = &models.RoleScheme{}
		result  *models.RoleScheme
		m       *models.Role
		id      string = chi.URLParam(r, "ID")
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error binding input data - %s", err)
	}

	if m, err = models.NewRole(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing role model - %s", err)
	}

	if err = m.Update(id, request); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error updating role '%s' - %s", id, err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting role '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}

func (rl *roles) remove(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		m   *models.Role
		id  string = chi.URLParam(r, "ID")
		err error
	)

	if id == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("error getting role id")
	}

	if m, err = models.NewRole(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing role model - %s", err)
	}

	if err = m.Delete(id); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error deleting role - %s", err)
	}

	return nil, http.StatusNoContent, nil
}
