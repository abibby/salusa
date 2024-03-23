package auth_test

import (
	"testing"

	"github.com/abibby/salusa/auth"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		Name   string
		Claims *auth.Claims
	}{
		{
			Name:   "empty",
			Claims: &auth.Claims{},
		},
		{
			Name: "sub",
			Claims: &auth.Claims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject: "sub",
				},
			},
		},
		{
			Name: "type",
			Claims: &auth.Claims{
				Type: "type",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			token, err := auth.GenerateTokenFrom(tc.Claims)
			assert.NoError(t, err)
			newClaims, err := auth.Parse(token)
			assert.NoError(t, err)
			assert.Equal(t, tc.Claims, newClaims)
		})
	}
}
