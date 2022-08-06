package handlers

import (
	"fmt"
	"net/http"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
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
