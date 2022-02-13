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
// It matches user permission with application math and method
func IsPermit(_ http.ResponseWriter, r *http.Request) (interface{}, int, error) {
	var (
		appname string
		// err  error
		user *models.UserScheme
		perm *models.PermissionScheme
	)

	getAppName := func(path string) string {
		root := "/" + config.TheConfig().Server.APIVersion
		part := strings.TrimPrefix(path, root)

		if len(part) < 1 {
			return ""
		}

		return strings.Split(part, "/")[1]
	}

	if appname = getAppName(r.URL.Path); appname == "" {
		return nil, http.StatusForbidden, fmt.Errorf("error parsing '%s' - unexpected path", r.URL.Path)
	}

	if user = r.Context().Value(models.UserSchemeType{}).(*models.UserScheme); user == nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error casting to *models.UserScheme")
	}

	if perm = user.GetPermission(appname); perm == nil {
		return nil, http.StatusForbidden, fmt.Errorf("user doesn't have access to appplication '%s'", appname)
	}

	if !perm.IsAllowed(r.Method) {
		return nil, http.StatusForbidden, fmt.Errorf("user doesn't have such permission '%s' to application '%s'", r.Method, appname)
	}

	return nil, http.StatusOK, nil
}
