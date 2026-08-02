package request

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun_decode_slice_from_query(t *testing.T) {
	type Request struct {
		Items []int `query:"items"`
	}

	httpRequest := httptest.NewRequest("GET", "http://0.0.0.0/?items=1&items=2", http.NoBody)
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.NoError(t, err)
	assert.Equal(t, &Request{Items: []int{1, 2}}, structRequest)
}

func TestRun_decode_slice_error(t *testing.T) {
	type Request struct {
		Items *[]int `query:"items"`
	}

	httpRequest := httptest.NewRequest("GET", "http://0.0.0.0/?items=foo", http.NoBody)
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.Error(t, err)
	assert.IsType(t, ValidationError{}, err)
}

func TestRun_invalid_json_body(t *testing.T) {
	type Request struct {
		Foo string `json:"foo"`
	}

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", bytes.NewBuffer([]byte(`not json`)))
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not decode body")
}

func TestRun_invalid_url_encoded_body(t *testing.T) {
	type Request struct {
		Foo string `json:"foo"`
	}

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", bytes.NewBuffer([]byte(`foo=%zz`)))
	httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.Error(t, err)
}

func TestRun_invalid_multipart_body(t *testing.T) {
	type Request struct {
		Foo string `json:"foo"`
	}

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", bytes.NewBuffer([]byte(`not multipart`)))
	httpRequest.Header.Add("Content-Type", "multipart/form-data; boundary=abc")
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.Error(t, err)
}

func TestRun_multipart_form_value(t *testing.T) {
	type Request struct {
		Foo string `json:"foo"`
	}

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	assert.NoError(t, writer.WriteField("foo", "bar"))
	assert.NoError(t, writer.Close())

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", buff)
	httpRequest.Header.Add("Content-Type", writer.FormDataContentType())
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.NoError(t, err)
	assert.Equal(t, &Request{Foo: "bar"}, structRequest)
}

func TestRun_multipart_missing_file(t *testing.T) {
	type Request struct {
		Foo fs.File `json:"foo"`
	}

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	assert.NoError(t, writer.Close())

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", buff)
	httpRequest.Header.Add("Content-Type", writer.FormDataContentType())
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.Error(t, err)
	assert.IsType(t, ValidationError{}, err)
}

func TestRun_body_read_error(t *testing.T) {
	type Request struct {
		Foo string `json:"foo"`
	}

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", nil)
	httpRequest.Body = failingReadCloser{}
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read error")
}

func TestRun_non_struct(t *testing.T) {
	httpRequest := httptest.NewRequest("GET", "http://0.0.0.0/", http.NoBody)
	err := Run(httpRequest, 5)
	assert.Error(t, err)
}

func TestRun_invalid_time_json(t *testing.T) {
	type Request struct {
		T time.Time `json:"t"`
	}

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", bytes.NewBuffer([]byte(`{"t": "not a time"}`)))
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.Error(t, err)
	assert.IsType(t, ValidationError{}, err)
}

func TestRun_json_field_not_in_body(t *testing.T) {
	type Request struct {
		Foo string `json:"foo"`
		Bar string `json:"bar"`
	}

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", bytes.NewBuffer([]byte(`{"foo": "x"}`)))
	structRequest := &Request{}

	err := Run(httpRequest, structRequest)

	assert.NoError(t, err)
	assert.Equal(t, &Request{Foo: "x"}, structRequest)
}

func TestFileInfo(t *testing.T) {
	type Request struct {
		Foo fs.File `json:"foo"`
	}

	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	part, err := writer.CreateFormFile("foo", "foo.txt")
	assert.NoError(t, err)
	_, err = part.Write([]byte("foo content"))
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	httpRequest := httptest.NewRequest("POST", "http://0.0.0.0/", buff)
	httpRequest.Header.Add("Content-Type", writer.FormDataContentType())
	structRequest := &Request{}

	err = Run(httpRequest, structRequest)
	assert.NoError(t, err)

	defer structRequest.Foo.Close()

	info, err := structRequest.Foo.Stat()
	assert.NoError(t, err)
	assert.Equal(t, "foo.txt", info.Name())
	assert.Equal(t, int64(len("foo content")), info.Size())
	assert.False(t, info.IsDir())
	assert.Equal(t, fs.FileMode(0644), info.Mode())
	assert.NotZero(t, info.ModTime())
	assert.Nil(t, info.Sys())

	content, err := io.ReadAll(structRequest.Foo)
	assert.NoError(t, err)
	assert.Equal(t, "foo content", string(content))
}

type failingReadCloser struct{}

func (failingReadCloser) Read([]byte) (int, error) { return 0, errors.New("read error") }
func (failingReadCloser) Close() error             { return nil }
