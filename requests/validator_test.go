package requests

import (
	"testing"

	"github.com/abibby/requests/rules"
	"github.com/stretchr/testify/assert"
)

func Test_Validate_fails_with_non_struct_arguments(t *testing.T) {
	err := Validate(nil, []string{}, 1)

	assert.Error(t, err)
}

func Test_Validate_generates_errors_on_failing_rules(t *testing.T) {
	rules.AddRule("should_fail", func(*rules.ValidationOptions) bool {
		return false
	})

	type Request struct {
		Foo int `validate:"should_fail"`
	}

	err := Validate(nil, []string{"Foo"}, &Request{})

	assert.Equal(t, ValidationError{
		"Foo": []string{"should_fail"},
	}, err)
}

func Test_Validate_ignores_failing_rules_if_no_value_is_passed(t *testing.T) {
	rules.AddRule("should_fail", func(*rules.ValidationOptions) bool {
		return false
	})

	type Request struct {
		Foo int `validate:"should_fail"`
	}

	err := Validate(nil, []string{}, &Request{})

	assert.NoError(t, err)
}

func Test_Validate_generates_no_errors_on_passing_rules(t *testing.T) {
	rules.AddRule("should_pass", func(*rules.ValidationOptions) bool {
		return true
	})

	type Request struct {
		Foo int `validate:"should_pass"`
	}

	err := Validate(nil, []string{"Foo"}, &Request{})

	assert.NoError(t, err)
}

func Test_Validate_multiple_errors(t *testing.T) {
	rules.AddRule("should_fail_1", func(*rules.ValidationOptions) bool {
		return false
	})
	rules.AddRule("should_fail_2", func(*rules.ValidationOptions) bool {
		return false
	})

	type Request struct {
		Foo int `validate:"should_fail_1"`
		Bar int `validate:"should_fail_1|should_fail_2"`
	}

	err := Validate(nil, []string{"Foo", "Bar"}, &Request{})

	assert.Equal(t, ValidationError{
		"Foo": []string{"should_fail_1"},
		"Bar": []string{"should_fail_1", "should_fail_2"},
	}, err)
}

func Test_Validate_uses_args(t *testing.T) {
	rules.AddRule("has_args", func(options *rules.ValidationOptions) bool {
		return options.Value.(string) == options.Arguments[0]
	})

	type Request struct {
		Foo string `validate:"has_args:foo"`
		Bar string `validate:"has_args:bar"`
	}

	err := Validate(nil, []string{"Foo", "Bar"}, &Request{
		Foo: "foo",
		Bar: "foo",
	})

	assert.Equal(t, ValidationError{
		"Bar": []string{"has_args bar"},
	}, err)
}
func Test_Validate_required(t *testing.T) {
	type Request struct {
		Foo int `validate:"required"`
	}

	err := Validate(nil, []string{}, &Request{})

	assert.Equal(t, ValidationError{
		"Foo": []string{"The Foo field is required."},
	}, err)
}
