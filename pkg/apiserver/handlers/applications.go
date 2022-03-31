package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
)

type applications struct {
	// *models.UserScheme // @TODO init current user from constructor
}

func appsHandler() *applications {
	return &applications{}
}

func (a *applications) get(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		am     *models.Application
		result []models.ApplicationScheme
		err    error
	)

	if am, err = models.NewApplication(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing application model - %s", err)
	}

	if result, err = am.FindAll(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting applications list - %s", err)
	}

	return result, http.StatusOK, nil
}

func (a *applications) getByID(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		aid    string = chi.URLParam(r, "ID")
		am     *models.Application
		result *models.ApplicationScheme
		err    error
	)

	if aid == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("error getting application id")
	}

	if am, err = models.NewApplication(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing application model - %s", err)
	}

	if result, err = am.FindByID(aid); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting application '%s' - %s", aid, err)
	}

	return result, http.StatusOK, nil
}

func (a *applications) create(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		request *models.ApplicationScheme = &models.ApplicationScheme{}
		result  *models.ApplicationScheme
		m       *models.Application
		id      string
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error binding input data - %s", err)
	}

	if m, err = models.NewApplication(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing application model - %s", err)
	}

	if id, err = m.Create(request); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error creating application - %s", err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting application '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}

func (a *applications) update(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		request *models.ApplicationScheme = &models.ApplicationScheme{}
		result  *models.ApplicationScheme
		m       *models.Application
		id      string = chi.URLParam(r, "ID")
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error binding input data - %s", err)
	}

	if m, err = models.NewApplication(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing application model - %s", err)
	}

	if err = m.Update(id, request); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error updating application '%s' - %s", id, err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting application '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}

func (a *applications) remove(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		m   *models.Application
		id  string = chi.URLParam(r, "ID")
		err error
	)

	if id == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("error getting application id")
	}

	if m, err = models.NewApplication(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing application model - %s", err)
	}

	if err = m.Delete(id); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error deleting application - %s", err)
	}

	return nil, http.StatusNoContent, nil
}
