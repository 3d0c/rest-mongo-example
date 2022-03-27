package handlers

import (
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
	"github.com/teal-seagull/lyre-be-v4/pkg/helpers"
	"github.com/teal-seagull/lyre-be-v4/pkg/log"
)

type sessions struct {
	logger *zap.Logger
}

func sessionsHandler() *sessions {
	return &sessions{
		logger: log.TheLogger().With(zap.String("handler", "sessionsHandler")),
	}
}

func (s *sessions) create(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		// request is used only for binding user name or email and password
		request *models.UserScheme = &models.UserScheme{}
		us      *models.UserScheme
		sm      *models.Session
		ss      *models.SessionScheme
		u       *models.User
		err     error
	)

	if err = render.Bind(r, request); err != nil {
		s.logger.Error("error binding input data", zap.Error(err))
		return nil, http.StatusBadRequest, nil
	}

	if u, err = models.NewUser(); err != nil {
		s.logger.Error("error initializing user model", zap.Error(err))
		return nil, http.StatusInternalServerError, nil
	}

	if us, err = u.FindByName(request); err != nil {
		s.logger.Error("error finding user by", zap.String("name", request.Name), zap.String("email", request.Email), zap.Error(err))
		return nil, http.StatusNotFound, nil
	}

	if !helpers.ComparePasswords(*us.Password, *request.Password) {
		s.logger.Error("error comparing password for", zap.String("user", request.Name), zap.String("email", request.Email))
		return nil, http.StatusUnauthorized, nil
	}

	if sm, err = models.NewSession(); err != nil {
		s.logger.Error("error initializing session model", zap.Error(err))
		return nil, http.StatusInternalServerError, nil
	}

	if ss, err = sm.Create(us.ID); err != nil {
		s.logger.Error("error creating session", zap.Error(err))
		return nil, http.StatusInternalServerError, nil
	}

	return ss, http.StatusOK, nil
}

func (s *sessions) remove(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		sm    *models.Session
		token string
		err   error
	)

	if token = r.Context().Value(helpers.TokenStringType{}).(string); token == "" {
		s.logger.Error("error getting token from context")
		return nil, http.StatusInternalServerError, nil
	}

	if sm, err = models.NewSession(); err != nil {
		s.logger.Error("error initializing session model", zap.Error(err))
		return nil, http.StatusInternalServerError, nil
	}

	if err = sm.Remove(token); err != nil {
		s.logger.Error("error removing session token", zap.Error(err))
		return nil, http.StatusNotFound, nil
	}

	return nil, http.StatusNoContent, nil
}
