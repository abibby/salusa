package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/schema"
)

func Run(requestHttp *http.Request, requestStruct any) error {
	var decoder = schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.SetAliasTag("query")

	err := decoder.Decode(requestStruct, requestHttp.URL.Query())
	if multiErr, ok := err.(schema.MultiError); ok {
		return fromSchemaMultiError(multiErr)
	} else if err != nil {
		return fmt.Errorf("could decode query string: %w", err)
	}

	if requestHttp.Body != http.NoBody {
		defer requestHttp.Body.Close()

		bodyBuff := bytes.Buffer{}
		body := io.TeeReader(requestHttp.Body, &bodyBuff)

		err := json.NewDecoder(body).Decode(requestStruct)
		if jsonErr, ok := err.(*json.UnmarshalTypeError); ok {
			return fromJsonUnmarshalTypeError(jsonErr, requestStruct)
		} else if err != nil {
			return fmt.Errorf("could decode body: %w", err)
		}

		m := map[string]json.RawMessage{}
		err = json.Unmarshal(bodyBuff.Bytes(), &m)
		if err != nil {
			return fmt.Errorf("could decode body: %w", err)
		}
	}

	err = Validate(requestHttp, requestStruct)
	if err != nil {
		return err
	}
	return nil
}
