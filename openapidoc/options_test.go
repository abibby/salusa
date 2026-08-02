package openapidoc

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	s := &spec.Swagger{}

	s = Info(spec.InfoProps{Title: "My API", Version: "1.0"})(s)
	assert.Equal(t, "My API", s.Info.Title)
	assert.Equal(t, "1.0", s.Info.Version)

	s = BasePath("/api/v1")(s)
	assert.Equal(t, "/api/v1", s.BasePath)

	s = AddSecurityDefinition("Key", &spec.SecuritySchemeProps{
		Type: "apiKey",
		In:   "header",
		Name: "X-Key",
	})(s)
	assert.NotNil(t, s.SecurityDefinitions)
	assert.Equal(t, "apiKey", s.SecurityDefinitions["Key"].Type)

	s = AddSecurityDefinition("Key2", &spec.SecuritySchemeProps{Type: "oauth2"})(s)
	assert.NotNil(t, s.SecurityDefinitions["Key2"])
	assert.Equal(t, "oauth2", s.SecurityDefinitions["Key2"].Type)

	s = &spec.Swagger{}
	s = AddDefaultSecurityDefinition()(s)
	assert.NotNil(t, s.SecurityDefinitions[DefaultSecurityDefinitionName])
	assert.Equal(t, "apiKey", s.SecurityDefinitions[DefaultSecurityDefinitionName].Type)
}
