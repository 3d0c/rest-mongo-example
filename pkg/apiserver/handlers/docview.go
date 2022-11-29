package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
	"github.com/teal-seagull/lyre-be-v4/pkg/config"
)

type docview struct {
	// *models.UserScheme // @TODO init current user from constructor
}

func docviewHandler() *docview {
	return &docview{}
}

func (*docview) get(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		item   = r.URL.Query().Get("item")
		typeid = r.URL.Query().Get("type")
		m      *models.Document
		result []models.DocumentScheme
		err    error
	)

	if item == "" || typeid == "" {
		return nil, http.StatusInternalServerError, fmt.Errorf("wrong request, item and typeid are mandatory")
	}

	if m, err = models.NewDocument(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing DocView model - %s", err)
	}

	if result, err = m.Find(item, typeid); err != nil {
		return nil, http.StatusNotFound, fmt.Errorf("error getting documents list - %s", err)
	}

	return result, http.StatusOK, nil
}

func (*docview) getFile(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		id   = chi.URLParam(r, "ID")
		path = config.TheConfig().Docview.Path
		name string
		m    *models.Document
		err  error
	)

	if id == "" {
		return nil, http.StatusInternalServerError, fmt.Errorf("wrong request, item and typeid are mandatory")
	}

	if m, err = models.NewDocument(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error initializing DocView model - %s", err)
	}

	if name, err = m.Download(id, path); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error downloading document content - %s", err)
	}

	// Manually adding CORS here, because this request is ignored by chain view
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.Header().Set(
		"Content-Disposition", "attachment; filename="+strconv.Quote(name))
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeFile(w, r, path+name)

	return nil, http.StatusOK, nil
}
