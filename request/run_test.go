package request

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun_unmarshals_data_from_query_string(t *testing.T) {
	type Request struct {
		Foo string `query:"foo"`
	}

	httpRequest := httptest.NewRequest("GET", "http://0.0.0.0/?foo=bar", http.NoBody)
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.NoError(t, err)
	assert.Equal(t, "bar", structRequest.Foo)
}

func TestRun_fails_with_invalid_type_from_query_string(t *testing.T) {
	type Request struct {
		Foo int `query:"foo"`
	}

	httpRequest := httptest.NewRequest("GET", "http://0.0.0.0/?foo=bar", http.NoBody)
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.Error(t, err)
	assert.IsType(t, ValidationError{}, err)
	assert.Equal(t, err, ValidationError{
		"foo": []string{"should be of type int"},
	})
}

func TestRun_unmarshals_json_data_from_body(t *testing.T) {
	type Request struct {
		Foo string `json:"foo"`
	}

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", bytes.NewBuffer([]byte(`{ "foo": "bar" }`)))
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.NoError(t, err)
	assert.Equal(t, "bar", structRequest.Foo)
}

func TestRun_fails_with_invalid_type_from_body(t *testing.T) {
	type Request struct {
		Foo int `json:"foo"`
	}

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", bytes.NewBuffer([]byte(`{ "foo": "bar" }`)))
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.Error(t, err)
	assert.IsType(t, ValidationError{}, err)
	assert.Equal(t, err, ValidationError{
		"foo": []string{"should be of type int"},
	})
}

func TestRun_fails_with_message(t *testing.T) {
	type Request struct {
		Foo int `json:"foo" validate:"max:1"`
	}

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", bytes.NewBuffer([]byte(`{ "foo": 10 }`)))
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.Error(t, err)
	assert.IsType(t, ValidationError{}, err)
	assert.Equal(t, err, ValidationError{
		"foo": []string{"The foo must not be greater than 1."},
	})
}

func TestRun_doesnt_fail_with_extra_query_params(t *testing.T) {
	type Request struct {
		Foo int `json:"foo"`
	}

	httpRequest := httptest.NewRequest("GET", "http://0.0.0.0/?bar=foo", http.NoBody)
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.NoError(t, err)
}

func TestRun_works_multiple_times(t *testing.T) {
	type Request struct {
		Foo int `json:"foo"`
	}

	httpRequest := httptest.NewRequest("GET", "http://0.0.0.0/?bar=foo", bytes.NewBuffer([]byte(`{ "foo": 10 }`)))

	structRequest := &Request{}
	err := Run(httpRequest, structRequest)
	assert.NoError(t, err)
	assert.Equal(t, 10, structRequest.Foo)

	structRequest = &Request{}
	err = Run(httpRequest, structRequest)
	assert.NoError(t, err)
	assert.Equal(t, 10, structRequest.Foo)
}

func TestRun_query_string_only_tagged_fields(t *testing.T) {
	type Request struct {
		Foo string `query:"foo"`
		Bar string
	}

	httpRequest := httptest.NewRequest("GET", "http://0.0.0.0/?foo=foo&Bar=bar", http.NoBody)
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Foo: "foo",
	}, structRequest)
}

func TestRun_query_url_body(t *testing.T) {
	type Request struct {
		Foo string `json:"foo"`
	}

	httpRequest := httptest.NewRequest("GET", "http://0.0.0.0/?foo=foo&Bar=bar", bytes.NewBuffer([]byte(`foo=bar`)))
	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.NoError(t, err)
	assert.Equal(t, &Request{
		Foo: "bar",
	}, structRequest)
}
