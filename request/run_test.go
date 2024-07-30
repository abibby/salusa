package request

import (
	"bytes"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abibby/nulls"
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
		"foo": []string{"should be of type int: value \"bar\""},
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

func TestRun_multipart_file(t *testing.T) {
	type Request struct {
		Foo fs.File `json:"foo"`
	}

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)

	part, err := writer.CreateFormFile("foo", "foo.txt")
	assert.NoError(t, err)
	_, err = part.Write([]byte("foo content"))
	assert.NoError(t, err)

	writer.Close()

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", buff)
	httpRequest.Header.Add("Content-Type", writer.FormDataContentType())

	structRequest := &Request{}

	err = Run(httpRequest, structRequest)
	assert.NoError(t, err)

	if !assert.NotNil(t, structRequest.Foo) {
		return
	}

	fooContent, err := io.ReadAll(structRequest.Foo)
	assert.NoError(t, err)

	defer structRequest.Foo.Close()

	assert.Equal(t, "foo content", string(fooContent))
}

type IntPtr struct {
	IntPtr *int `query:"int_ptr"`
}
type NullsInt struct {
	IntPtr *nulls.Int `query:"int_ptr"`
}
type TimeReq struct {
	Time time.Time `query:"time"`
}
type TimePtrReq struct {
	Time *time.Time `query:"time"`
}

func ptr[T any](v T) *T {
	return &v
}

func TestRun(t *testing.T) {
	type args struct {
		requestHttp   *http.Request
		requestStruct any
	}
	tests := []struct {
		name        string
		args        args
		wantRequest any
		wantErr     bool
	}{
		{
			name: "int pointer",
			args: args{
				httptest.NewRequest("GET", "https://example.com?int_ptr=1", http.NoBody),
				&IntPtr{},
			},
			wantRequest: &IntPtr{IntPtr: ptr(1)},
			wantErr:     false,
		},
		{
			name: "int pointer empty",
			args: args{
				httptest.NewRequest("GET", "https://example.com?int_ptr=", http.NoBody),
				&IntPtr{},
			},
			wantRequest: &IntPtr{IntPtr: nil},
			wantErr:     false,
		},
		{
			name: "int pointer missing",
			args: args{
				httptest.NewRequest("GET", "https://example.com", http.NoBody),
				&IntPtr{},
			},
			wantRequest: &IntPtr{IntPtr: nil},
			wantErr:     false,
		},
		{
			name: "nulls.Int",
			args: args{
				httptest.NewRequest("GET", "https://example.com?int_ptr=1", http.NoBody),
				&NullsInt{},
			},
			wantRequest: &NullsInt{IntPtr: nulls.NewInt(1)},
			wantErr:     false,
		},
		{
			name: "time",
			args: args{
				httptest.NewRequest("GET", "https://example.com?time=2020-01-01T00:00:00Z", http.NoBody),
				&TimeReq{},
			},
			wantRequest: &TimeReq{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)},
			wantErr:     false,
		},
		{
			name: "time pointer",
			args: args{
				httptest.NewRequest("GET", "https://example.com?time=2020-01-01T00:00:00Z", http.NoBody),
				&TimePtrReq{},
			},
			wantRequest: &TimePtrReq{Time: ptr(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))},
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Run(tt.args.requestHttp, tt.args.requestStruct); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.wantRequest, tt.args.requestStruct)
		})
	}
}
