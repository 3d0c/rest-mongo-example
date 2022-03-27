package helpers

import (
	"fmt"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

func TestTokenUserID(t *testing.T) {
	const (
		strUserID = "123"
		intUserID = 123
	)

	testCases := []struct {
		description string
		input       string
		expected    string
		assertFn    func(assert.TestingT, interface{}, ...interface{}) bool
	}{
		{
			description: "Happy Pass",
			input:       strUserID,
			expected:    strUserID,
			assertFn:    assert.Nil,
		},
		{
			description: "Empty userID",
			input:       "",
			expected:    "",
			assertFn:    assert.NotNil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			tokenString, err := CreateToken(testCase.input)
			assert.Nil(t, err)

			authClaim, err := VerifyToken(tokenString)
			assert.Nil(t, err)

			obtained := authClaim.GetUserID()

			assert.Equal(t, testCase.expected, obtained)
		})
	}
}

func TestTokenString(t *testing.T) {
	testCases := []struct {
		description string
		tokenFn     func() string
	}{
		{
			description: "Empty token",
			tokenFn:     func() string { return "" },
		},
		{
			description: "Wrong token payload",
			tokenFn:     func() string { return "xArJEcipJnFmog4sqn2" },
		},
		{
			description: "Invalid signing method",
			tokenFn: func() string {
				token, _ := jwt.NewWithClaims(jwt.SigningMethodNone, &AuthClaims{
					"userID",
					jwt.StandardClaims{
						Issuer: "lyre-be-v4",
					},
				}).SignedString([]byte{})
				return token
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			_, err := VerifyToken(testCase.tokenFn())
			assert.NotNil(t, err)
		})
	}
}

func TestHashPassword(t *testing.T) {
	const (
		password = "default_password"
	)

	hash, err := HashPassword(password)
	assert.Nil(t, err)

	// do not remove. for testing purposes
	fmt.Printf("%s\n", hash)

	result := ComparePasswords(hash, password)
	assert.True(t, result)
}
