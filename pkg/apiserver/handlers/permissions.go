package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
)

type permissions struct {
	// *models.UserScheme // @TODO init current user from constructor
}

func permHandler() *permissions {
	return &permissions{}
}

func (p *permissions) get(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		m      *models.Permission
		result []models.PermissionScheme
		err    error
	)

	if m, err = models.NewPermission(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing permission model - %s", err)
	}

	if result, err = m.FindAll(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting permissions list - %s", err)
	}

	return result, http.StatusOK, nil
}

func (p *permissions) getByID(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		id     string = chi.URLParam(r, "ID")
		m      *models.Permission
		result *models.PermissionScheme
		err    error
	)

	if id == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("error getting permission id")
	}

	if m, err = models.NewPermission(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing permission model - %s", err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting permission '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}

func (p *permissions) create(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		request *models.PermissionScheme = &models.PermissionScheme{}
		result  *models.PermissionScheme
		m       *models.Permission
		id      string
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error binding input data - %s", err)
	}

	if m, err = models.NewPermission(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing permission model - %s", err)
	}

	if id, err = m.Create(request); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error creating permission - %s", err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting permission object '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}

func (p *permissions) update(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		request *models.PermissionScheme = &models.PermissionScheme{}
		result  *models.PermissionScheme
		m       *models.Permission
		id      string = chi.URLParam(r, "ID")
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error binding input data - %s", err)
	}

	if m, err = models.NewPermission(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing permission model - %s", err)
	}

	if err = m.Update(id, request); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error updating permission '%s' - %s", id, err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting permission '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}

func (p *permissions) remove(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		m   *models.Permission
		id  string = chi.URLParam(r, "ID")
		err error
	)

	if id == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("error getting permission id")
	}

	if m, err = models.NewPermission(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing permission model - %s", err)
	}

	if err = m.Delete(id); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error deleting permission - %s", err)
	}

	return nil, http.StatusNoContent, nil
}
