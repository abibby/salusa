package openapidoc

import "github.com/go-openapi/spec"

type SwaggerOption func(*spec.Swagger) *spec.Swagger

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
