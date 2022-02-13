package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/teal-seagull/lyre-be-v4/pkg/helpers"
)

const (
	testUserID = "123"
)

func TestIsAuthorized(t *testing.T) {
	var (
		status     int
		validToken string
		err        error
	)

	if validToken, err = helpers.CreateToken(testUserID); err != nil {
		t.Fatalf("Error creating valid token - %s\n", err)
	}

	testCases := []struct {
		description        string
		req                *http.Request
		expectedStatusCode int
		middleware         func(http.ResponseWriter, *http.Request) (interface{}, int, error)
		assertFn           func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{
			description:        "No authorization header",
			req:                httptest.NewRequest(http.MethodGet, "https:", nil),
			expectedStatusCode: http.StatusUnauthorized,
			middleware:         IsAuthorized,
			assertFn:           assert.NotNil,
		},
		{
			description: "Wrong token",
			req: withHeaders(
				httptest.NewRequest(http.MethodPut, "https:", nil),
				map[string]string{
					authHeaderKey: bearerSchemePrefix + "123",
				},
			),
			expectedStatusCode: http.StatusUnauthorized,
			middleware:         IsAuthorized,
			assertFn:           assert.NotNil,
		},
		{
			description: "Happy Pass",
			req: withHeaders(
				httptest.NewRequest(http.MethodPut, "https:", nil),
				map[string]string{
					authHeaderKey: bearerSchemePrefix + validToken,
				},
			),
			expectedStatusCode: http.StatusOK,
			middleware:         IsAuthorized,
			assertFn:           assert.Nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			_, status, err = testCase.middleware(nil, testCase.req)
			testCase.assertFn(t, err)
			assert.Equal(t, testCase.expectedStatusCode, status)
		})
	}
}

func TestIsAuthorizedContext(t *testing.T) {
	var (
		status     int
		validToken string
		err        error
		r          *http.Request
	)

	if validToken, err = helpers.CreateToken(testUserID); err != nil {
		t.Fatalf("Error creating valid token - %s\n", err)
	}

	r = withHeaders(
		httptest.NewRequest(http.MethodPut, "https://a/api/v1/application/u", nil),
		map[string]string{
			authHeaderKey: bearerSchemePrefix + validToken,
		},
	)

	_, status, err = IsAuthorized(nil, r)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, status)

	v := r.Context().Value(helpers.AuthClaims{})
	assert.NotNil(t, v)

	claims, ok := v.(*helpers.AuthClaims)
	assert.Equal(t, true, ok)

	obtained := claims.GetUserID()
	assert.Nil(t, err)

	assert.Equal(t, testUserID, obtained)
}

func withHeaders(req *http.Request, headers map[string]string) *http.Request {
	for h, v := range headers {
		req.Header.Set(h, v)
	}

	return req
}
