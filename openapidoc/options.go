package openapidoc

import "github.com/go-openapi/spec"

type SwaggerOption func(*spec.Swagger) *spec.Swagger

const DefaultSecurityDefinitionName = "Bearer"

func Info(info spec.InfoProps) SwaggerOption {
	return func(s *spec.Swagger) *spec.Swagger {
		s.Info = &spec.Info{
			InfoProps: info,
		}
		return s
	}
}

func BasePath(basePath string) SwaggerOption {
	return func(s *spec.Swagger) *spec.Swagger {
		s.BasePath = basePath
		return s
	}
}

func SecurityDefinition(name string, securityScheme *spec.SecuritySchemeProps) SwaggerOption {
	return func(s *spec.Swagger) *spec.Swagger {
		if s.SecurityDefinitions == nil {
			s.SecurityDefinitions = spec.SecurityDefinitions{}
		}
		s.SecurityDefinitions[name] = &spec.SecurityScheme{SecuritySchemeProps: *securityScheme}
		return s
	}
}

func DefaultSecurityDefinition() SwaggerOption {
	return SecurityDefinition(DefaultSecurityDefinitionName, &spec.SecuritySchemeProps{
		Type:        "apiKey",
		In:          "header",
		Name:        "Authorization",
		Description: "The api key must have the `Bearer ` prefix, e.g. `Authorization: Bearer abcd1234`",
	})
}
