package handlers

import (
	"net/http"
)

func nilHandler(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	return nil, http.StatusOK, nil
}
