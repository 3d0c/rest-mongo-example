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
		perm    *models.PermissionScheme
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

	if perm = user.GetPermission(appPath); perm == nil {
		return nil, http.StatusForbidden, fmt.Errorf("user doesn't have access to appplication '%s'", appPath)
	}

	if !perm.IsAllowed(r.Method) {
		return nil, http.StatusForbidden, fmt.Errorf("user doesn't have such permission '%s' to application '%s'", r.Method, appPath)
	}

	return nil, http.StatusOK, nil
}
