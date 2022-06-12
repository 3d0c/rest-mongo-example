package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
	"github.com/teal-seagull/lyre-be-v4/pkg/config"
)

// IsPermit middleware should be third after IsAuthorized and GetUser
//
// It matches user permission with application path and method
func IsPermit(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		appPath string
		user    *models.UserScheme
	)

	getAppPath := func(path string) string {
		root := "/" + config.TheConfig().Server.APIVersion + "/"
		path = strings.TrimPrefix(path, root)

		parts := strings.Split(path, "/")
		if len(parts) < 1 {
			return ""
		}

		return "/" + parts[0]
	}

	if appPath = getAppPath(r.URL.Path); appPath == "" {
		return nil, http.StatusForbidden, fmt.Errorf("error parsing '%s' - unexpected path", r.URL.Path)
	}

	if user = r.Context().Value(models.UserSchemeType{}).(*models.UserScheme); user == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error casting to *models.UserScheme")
	}

	if !user.IsAllowed(appPath) {
		return nil, http.StatusForbidden, fmt.Errorf("user doesn't have permission to application '%s'", appPath)
	}

	return nil, http.StatusOK, nil
}
