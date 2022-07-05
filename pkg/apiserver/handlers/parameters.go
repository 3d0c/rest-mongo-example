package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
)

type parameters struct{}

func parametersHandler() *parameters {
	return &parameters{}
}

func (p *parameters) get(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		m      *models.Parameter
		result []models.ParameterScheme
		err    error
	)

	if m, err = models.NewParameter(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing parameter model - %s", err)
	}

	if result, err = m.FindAll(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting parameters list - %s", err)
	}

	return result, http.StatusOK, nil
}

func (p *parameters) getByID(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		id     string = chi.URLParam(r, "ID")
		m      *models.Parameter
		result *models.ParameterScheme
		err    error
	)

	if id == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("error getting parameter id")
	}

	if m, err = models.NewParameter(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing parameter model - %s", err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting parameter '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}

func (p *parameters) create(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		request *models.ParameterScheme = &models.ParameterScheme{}
		result  *models.ParameterScheme
		m       *models.Parameter
		id      string
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error binding input data - %s", err)
	}

	if m, err = models.NewParameter(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing parameter model - %s", err)
	}

	if id, err = m.Create(request); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error creating parameter - %s", err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting parameter object '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}

func (p *parameters) update(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		request *models.ParameterScheme = &models.ParameterScheme{}
		result  *models.ParameterScheme
		m       *models.Parameter
		id      string = chi.URLParam(r, "ID")
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error binding input data - %s", err)
	}

	if m, err = models.NewParameter(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing parameter model - %s", err)
	}

	if err = m.Update(id, request); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error updating parameter '%s' - %s", id, err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting parameter '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}

// @TODO Cascade delete
func (p *parameters) remove(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		m   *models.Parameter
		id  string = chi.URLParam(r, "ID")
		err error
	)

	if id == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("error getting parameter id")
	}

	if m, err = models.NewParameter(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing parameter model - %s", err)
	}

	if err = m.Delete(id); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error deleting parameter - %s", err)
	}

	return nil, http.StatusNoContent, nil
}
