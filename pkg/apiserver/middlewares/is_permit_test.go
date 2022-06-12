package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/teal-seagull/lyre-be-v4/pkg/apiserver/models"
	"github.com/teal-seagull/lyre-be-v4/pkg/config"
)

func TestIsPermit(t *testing.T) {
	config.TheConfig().Server.APIVersion = "api/v1"

	testCases := []struct {
		description    string
		req            *http.Request
		expectedStatus int
		assertFn       func(assert.TestingT, interface{}, ...interface{}) bool
		testUser       *models.UserScheme
	}{
		{
			description:    "Wrong application",
			req:            httptest.NewRequest(http.MethodPut, "https://a/api/v1/application/u", nil),
			expectedStatus: 403,
			assertFn:       assert.NotNil,
			testUser: &models.UserScheme{
				ACL: []models.ApplicationScheme{
					{
						Name: "Example One",
						Path: "/wrongapplication",
					},
				},
			},
		},
		{
			description:    "Happy path",
			req:            httptest.NewRequest(http.MethodGet, "https://a/api/v1/application/u", nil),
			expectedStatus: 200,
			assertFn:       assert.Nil,
			testUser: &models.UserScheme{
				ACL: []models.ApplicationScheme{
					{
						Name: "Example One",
						Path: "/application",
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			ctx := testCase.req.Context()
			ctx = context.WithValue(ctx, models.UserSchemeType{}, testCase.testUser)

			*testCase.req = *testCase.req.WithContext(ctx)

			_, status, err := IsPermit(nil, testCase.req)
			testCase.assertFn(t, err)
			assert.Equal(t, testCase.expectedStatus, status)
		})
	}
}
