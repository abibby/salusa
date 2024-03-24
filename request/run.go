package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/schema"
)

func Run(requestHttp *http.Request, requestStruct any) error {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.SetAliasTag("query")

	err := decoder.Decode(requestStruct, requestHttp.URL.Query())
	if multiErr, ok := err.(schema.MultiError); ok {
		return fromSchemaMultiError(multiErr)
	} else if err != nil {
		return fmt.Errorf("could not decode query string: %w", err)
	}

	if requestHttp.Body != http.NoBody {
		defer requestHttp.Body.Close()
		body, err := io.ReadAll(requestHttp.Body)
		if err != nil {
			return err
		}

		requestHttp.Body = io.NopCloser(bytes.NewBuffer(body))

		contentType := requestHttp.Header.Get("Content-Type")
		switch contentType {
		case "application/x-www-form-urlencoded":
			bodyDecoder := schema.NewDecoder()
			bodyDecoder.IgnoreUnknownKeys(true)
			bodyDecoder.SetAliasTag("json")

			v, err := url.ParseQuery(string(body))
			if err != nil {
				return err
			}
			err = bodyDecoder.Decode(requestStruct, v)
			if multiErr, ok := err.(schema.MultiError); ok {
				return fromSchemaMultiError(multiErr)
			} else if err != nil {
				return fmt.Errorf("could not decode body: %w", err)
			}
		default:
			err := json.Unmarshal(body, requestStruct)
			if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
				return fromJsonUnmarshalTypeError(jsonErr, requestStruct)
			} else if err != nil {
				return fmt.Errorf("could not decode body: %w", err)
			}
		}
	}

	err = Validate(requestHttp, requestStruct)
	if err != nil {
		return err
	}
	return nil
}
