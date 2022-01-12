package handlers

import (
	"net/http"
)

type sessions struct {
}

func sessionsHandler() *sessions {
	return &sessions{}
}

func (s *sessions) create(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	result := struct {
		Items []string `json:"items"`
	}{
		[]string{},
	}

	return result, http.StatusOK, nil
}

func (s *sessions) remove(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	return nil, http.StatusNoContent, nil
}
