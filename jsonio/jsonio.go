package jsonio

import (
	"encoding/json"
	"io"
)

type JsonReader io.PipeReader

var _ io.Reader = (*JsonReader)(nil)

func NewReader(v any) *JsonReader {
	r, w := io.Pipe()
	go func() {
		err := json.NewEncoder(w).Encode(v)
		w.CloseWithError(err)
	}()
	return (*JsonReader)(r)
}

// Read implements io.Reader.
func (j *JsonReader) Read(p []byte) (n int, err error) {
	return (*io.PipeReader)(j).Read(p)
}

type JsonWriter io.PipeWriter

var _ io.Writer = (*JsonWriter)(nil)

func NewWriter(v any) *JsonWriter {
	r, w := io.Pipe()
	go func() {
		err := json.NewDecoder(r).Decode(v)
		w.CloseWithError(err)
	}()
	return (*JsonWriter)(w)
}

// Read implements io.Writer.
func (j *JsonWriter) Write(p []byte) (n int, err error) {
	return (*io.PipeWriter)(j).Write(p)
}
