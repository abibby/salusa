package rules

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStringEdgeCases(t *testing.T) {
	data := map[string]TestCase{
		"ends_with-no-args":        {"ends_with", &ValidationOptions{Value: "string"}, false},
		"starts_with-no-args":      {"starts_with", &ValidationOptions{Value: "string"}, false},
		"not_regex-no-args":        {"not_regex", &ValidationOptions{Value: "anything"}, true},
		"not_regex-invalid-regex":  {"not_regex", &ValidationOptions{Value: "match", Arguments: []string{"["}}, true},
		"regex-no-args":            {"regex", &ValidationOptions{Value: "match"}, true},
		"regex-invalid-regex":      {"regex", &ValidationOptions{Value: "match", Arguments: []string{"["}}, true},
		"length-non-int-arg":       {"length", &ValidationOptions{Value: "foo", Arguments: []string{"not int"}}, true},
		"length_between-bad-min":   {"length_between", &ValidationOptions{Value: "foo", Arguments: []string{"not int", "5"}}, true},
		"length_between-bad-max":   {"length_between", &ValidationOptions{Value: "foo", Arguments: []string{"1", "not int"}}, true},
		"length_between-fail":      {"length_between", &ValidationOptions{Value: "a long string", Arguments: []string{"1", "3"}}, false},
		"in-pass":                  {"in", &ValidationOptions{Value: "foo", Arguments: []string{"bar", "foo", "baz"}}, true},
		"in-fail":                  {"in", &ValidationOptions{Value: "qux", Arguments: []string{"bar", "foo"}}, false},
		"not_in-pass":              {"not_in", &ValidationOptions{Value: "qux", Arguments: []string{"bar", "foo"}}, true},
		"not_in-fail":              {"not_in", &ValidationOptions{Value: "foo", Arguments: []string{"bar", "foo"}}, false},
		"string-rule-non-string":   {"alpha", &ValidationOptions{Value: 123}, true},
	}

	runTests(t, data)
}

func TestAddTypeRuleMissingHandlers(t *testing.T) {
	onlyInt := &TypeRule{ArgCount: 0, Int: func(value int64, arguments TypeRuleArguments) bool { return true }}
	onlyString := &TypeRule{ArgCount: 0, String: func(value string, arguments TypeRuleArguments) bool { return true }}

	AddTypeRule("only_int_rule", onlyInt)
	AddTypeRule("only_string_rule", onlyString)

	data := map[string]TestCase{
		"int-missing-rule":             {"only_string_rule", &ValidationOptions{Value: 1}, true},
		"uint-missing-rule":            {"only_string_rule", &ValidationOptions{Value: uint(1)}, true},
		"float-missing-rule":           {"only_string_rule", &ValidationOptions{Value: 1.5}, true},
		"time-missing-rule":            {"only_string_rule", &ValidationOptions{Value: time.Now()}, true},
		"slice-missing-rule":           {"only_string_rule", &ValidationOptions{Value: []int{1}}, true},
		"string-missing-rule-pass":     {"only_string_rule", &ValidationOptions{Value: "foo"}, true},
		"string-missing-rule-fail":     {"only_string_rule", &ValidationOptions{Value: "bar"}, true},
		"int-with-int-rule":            {"only_int_rule", &ValidationOptions{Value: 1}, true},
		"non-numeric-field":            {"only_int_rule", &ValidationOptions{Value: struct{}{}}, true},
	}

	runTests(t, data)
}

func TestAddTypeRuleArgCount(t *testing.T) {
	AddTypeRule("needs_one_arg", &TypeRule{
		ArgCount: 1,
		Int:      func(value int64, arguments TypeRuleArguments) bool { return value == arguments.GetInt(0) },
	})

	rule, ok := GetRule("needs_one_arg")
	assert.True(t, ok)

	assert.True(t, rule(&ValidationOptions{Value: 5, Arguments: []string{"5"}}))
	assert.True(t, rule(&ValidationOptions{Value: 5, Arguments: []string{}}))
	assert.False(t, rule(&ValidationOptions{Value: 5, Arguments: []string{"not int"}}))
}

func TestTypeRuleArguments(t *testing.T) {
	args := TypeRuleArguments{"1", "2", "2.5", "2020-01-01T00:00:00Z", "true", "-3"}

	assert.Equal(t, "", args.GetString(10))
	assert.Equal(t, "1", args.GetString(0))
	assert.Equal(t, int64(1), args.GetInt(0))
	assert.Equal(t, uint64(2), args.GetUint(1))
	assert.Equal(t, 2.5, args.GetFloat(2))
	assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), args.GetTime(3))
	assert.True(t, args.GetBoolean(4))
	assert.True(t, args.GetBoolean(0))
	assert.False(t, args.GetBoolean(1))
	assert.Equal(t, int64(-3), args.GetInt(5))
}

func TestGetBoolean(t *testing.T) {
	args := TypeRuleArguments{"1", "true", "yes", "0", "no"}
	assert.True(t, args.GetBoolean(0))
	assert.True(t, args.GetBoolean(1))
	assert.True(t, args.GetBoolean(2))
	assert.False(t, args.GetBoolean(3))
	assert.False(t, args.GetBoolean(4))
}
