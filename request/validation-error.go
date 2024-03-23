package request

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/gorilla/schema"
)

type ValidationError map[string][]string

var _ error = ValidationError{}

func (e ValidationError) Error() string {
	// return "validation error"
	return fmt.Sprintf("validation error (%d)", len(e))
}

func (e ValidationError) HasErrors() bool {
	return len(e) > 0
}

func (e ValidationError) AddError(key string, message string) ValidationError {
	messages := e[key]
	if messages == nil {
		messages = []string{message}
	} else {
		messages = append(messages, message)
	}
	e[key] = messages
	return e
}
func (e ValidationError) Merge(vErr ValidationError) {
	for k, v := range vErr {
		errs := e[k]
		if errs == nil {
			errs = v
		} else {
			errs = append(errs, v...)
		}
		e[k] = errs
	}
}

func fromSchemaMultiError(err schema.MultiError) ValidationError {
	validationErr := ValidationError{}
	for key, subErr := range err {
		if err, ok := subErr.(schema.ConversionError); ok {
			validationErr[key] = []string{
				fmt.Sprintf("should be of type %s", err.Type.String()),
			}
		} else {
			validationErr[key] = []string{subErr.Error()}
		}
	}
	return validationErr
}

func fromJsonUnmarshalTypeError(err *json.UnmarshalTypeError, requestStruct any) ValidationError {
	validationErr := ValidationError{}
	t := reflect.TypeOf(requestStruct)
	key := err.Field
	f, ok := t.Elem().FieldByName(err.Field)
	if ok {
		jsonKey := f.Tag.Get("json")
		if jsonKey != "" {
			key = jsonKey
		}
	}
	validationErr[key] = []string{
		fmt.Sprintf("should be of type %s", err.Type.String()),
	}
	return validationErr
}
