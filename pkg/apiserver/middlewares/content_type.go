package middlewares

import (
	"fmt"
	"net/http"
)

const (
	jsonContentType = "application/json"
)

// IsValidContentType checks header for valid content-type record
func IsValidContentType(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		ct string = r.Header.Get("Content-Type")
	)

	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
		fallthrough
	case http.MethodDelete:
		fallthrough
	case http.MethodPut:
		if ct != jsonContentType {
			return nil, http.StatusBadRequest, fmt.Errorf("wrong Content-type '%s'", ct)
		}
	}

	return nil, http.StatusOK, nil
}
