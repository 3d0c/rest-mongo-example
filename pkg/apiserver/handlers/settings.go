package handlers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
	"github.com/teal-seagull/lyre-be-v4/pkg/helpers"
)

type settings struct{}

func settingsHandler() *settings {
	return &settings{}
}

func (*settings) get(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		userID string = r.URL.Query().Get("user_id")
		appID  string = r.URL.Query().Get("app_id")
		m      *models.Setting
		result []models.SettingScheme
		err    error
	)

	if m, err = models.NewSetting(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing setting model - %s", err)
	}

	if result, err = m.FindAll(userID, appID); err != nil {
		if err == helpers.ErrNotFound {
			return nil, http.StatusNotFound, fmt.Errorf("settings for user_id='%s', app_id='%s' not found", userID, appID)
		}
		return nil, http.StatusInternalServerError, fmt.Errorf("error getting setting list - %s", err)
	}

	return result, http.StatusOK, nil
}

func (*settings) getByID(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		id     string = chi.URLParam(r, "ID")
		m      *models.Setting
		result *models.SettingScheme
		err    error
	)

	if id == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("error getting setting id")
	}

	if m, err = models.NewSetting(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing setting model - %s", err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting setting '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}
func (*settings) create(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		request *models.SettingScheme = &models.SettingScheme{}
		result  *models.SettingScheme
		m       *models.Setting
		id      string
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error binding input data - %s", err)
	}

	if m, err = models.NewSetting(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing setting model - %s", err)
	}

	if id, err = m.Create(request); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error creating setting - %s", err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting setting object '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}

func (*settings) update(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		request *models.SettingScheme = &models.SettingScheme{}
		result  *models.SettingScheme
		m       *models.Setting
		id      string = chi.URLParam(r, "ID")
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		return nil, http.StatusBadRequest, fmt.Errorf("error binding input data - %s", err)
	}

	if m, err = models.NewSetting(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing setting model - %s", err)
	}

	if err = m.Update(id, request); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error updating setting '%s' - %s", id, err)
	}

	if result, err = m.FindByID(id); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting setting '%s' - %s", id, err)
	}

	return result, http.StatusOK, nil
}

// @TODO Cascade delete
func (*settings) remove(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		m   *models.Setting
		id  string = chi.URLParam(r, "ID")
		err error
	)

	if id == "" {
		return nil, http.StatusBadRequest, fmt.Errorf("error getting setting id")
	}

	if m, err = models.NewSetting(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing setting model - %s", err)
	}

	if err = m.Delete(id); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error deleting setting - %s", err)
	}

	return nil, http.StatusNoContent, nil
}
