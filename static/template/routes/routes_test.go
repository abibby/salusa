package routes_test

import (
	"testing"

	"abibby.com/salusa/router"
	"abibby.com/salusa/static/template/routes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitRoutes(t *testing.T) {
	r := router.New()
	require.NotNil(t, r)

	routes.InitRoutes(r)

	paths := map[string]bool{}
	for _, route := range r.Routes() {
		paths[route.Path] = true
	}
	for _, want := range []string{"/", "/login", "/user/create", "/docs/*", "/api/user", "/api/user/{id}"} {
		assert.True(t, paths[want], "missing route %q", want)
	}
}
