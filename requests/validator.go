package requests

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/abibby/requests/rules"
)

type Validator interface {
	Valid() error
}

func Validate(request *http.Request, keys []string, v any) error {
	s, err := getStruct(reflect.ValueOf(v))
	if err != nil {
		return fmt.Errorf("Validate mast take a struct or pointer to a struct: %w", err)
	}
	t := s.Type()

	vErr := ValidationError{}

	for i := 0; i < s.NumField(); i++ {
		ft := t.Field(i)
		fv := s.Field(i)
		name := getName(ft)
		err := validateField(name, request, keys, ft, fv)
		if err != nil {
			vErr.Merge(err)
		}
	}

	if vErr.HasErrors() {
		return vErr
	}

	return nil
}

func validateField(attribute string, request *http.Request, keys []string, ft reflect.StructField, fv reflect.Value) ValidationError {
	validate, ok := ft.Tag.Lookup("validate")
	if !ok {
		return nil
	}

	vErr := ValidationError{}

	if validator, ok := fv.Interface().(Validator); ok {
		err := validator.Valid()
		if err != nil {
			vErr.AddError(attribute, err.Error())
		}
	}

	rulesStr := strings.Split(validate, "|")
	for _, ruleStr := range rulesStr {
		ruleName, argsStr := split(ruleStr, ":")
		args := filterZeros(strings.Split(argsStr, ","))
		hasKey := includes(keys, attribute)
		if !hasKey {
			if ruleName == "required" {
				vErr.AddError(attribute, getMessage(ruleName, &MessageOptions{
					Attribute: attribute,
					Value:     fv.Interface(),
					Arguments: args,
				}))
			} else {
				return nil
			}
		} else {
			rule, ok := rules.GetRule(ruleName)
			if !ok {
				continue
			}

			valid := rule(&rules.ValidationOptions{
				Value:     fv.Interface(),
				Arguments: args,
				Request:   request,
				Name:      attribute,
			})
			if !valid {
				vErr.AddError(attribute, getMessage(ruleName, &MessageOptions{
					Attribute: attribute,
					Value:     fv.Interface(),
					Arguments: args,
				}))
			}

		}
	}
	return vErr
}

func getName(f reflect.StructField) string {
	name, ok := f.Tag.Lookup("json")
	if ok {
		return name
	}

	name, ok = f.Tag.Lookup("query")
	if ok {
		return name
	}

	name, ok = f.Tag.Lookup("di")
	if ok {
		return name
	}

	return f.Name
}

func getStruct(v reflect.Value) (reflect.Value, error) {
	switch v.Kind() {
	case reflect.Struct:
		return v, nil
	case reflect.Interface, reflect.Pointer:
		return getStruct(v.Elem())
	default:
		return reflect.Value{}, fmt.Errorf("expected struct found %s", v.Kind())
	}
}

func split(s, sep string) (string, string) {
	parts := strings.SplitN(s, sep, 2)
	switch len(parts) {
	case 0:
		return "", ""
	case 1:
		return parts[0], ""
	default:
		return parts[0], parts[1]
	}
}

func includes[T comparable](haystack []T, needle T) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
func filterZeros[T comparable](array []T) []T {
	newArray := []T{}
	var zero T
	for _, v := range array {
		if v != zero {
			newArray = append(newArray, v)
		}
	}
	return newArray
}
