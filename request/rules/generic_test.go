package rules

import (
	"testing"
)

func TestGeneric(t *testing.T) {
	data := map[string]TestCase{
		"gt-int-pass":     {"gt", &ValidationOptions{Value: 2, Arguments: []string{"1"}}, true},
		"gt-int-fail":     {"gt", &ValidationOptions{Value: 1, Arguments: []string{"1"}}, false},
		"gt-uint-pass":    {"gt", &ValidationOptions{Value: uint(2), Arguments: []string{"1"}}, true},
		"gt-uint-fail":    {"gt", &ValidationOptions{Value: uint(1), Arguments: []string{"1"}}, false},
		"gt-float-pass":   {"gt", &ValidationOptions{Value: 2.0, Arguments: []string{"1"}}, true},
		"gt-float-fail":   {"gt", &ValidationOptions{Value: 1.0, Arguments: []string{"1"}}, false},
		"gt-string-pass":  {"gt", &ValidationOptions{Value: "b", Arguments: []string{"a"}}, true},
		"gt-string-fail":  {"gt", &ValidationOptions{Value: "a", Arguments: []string{"a"}}, false},
		"gte-int-pass":    {"gte", &ValidationOptions{Value: 1, Arguments: []string{"1"}}, true},
		"gte-int-fail":    {"gte", &ValidationOptions{Value: 0, Arguments: []string{"1"}}, false},
		"gte-uint-pass":   {"gte", &ValidationOptions{Value: uint(1), Arguments: []string{"1"}}, true},
		"gte-uint-fail":   {"gte", &ValidationOptions{Value: uint(0), Arguments: []string{"1"}}, false},
		"gte-float-pass":  {"gte", &ValidationOptions{Value: 1.0, Arguments: []string{"1"}}, true},
		"gte-float-fail":  {"gte", &ValidationOptions{Value: 0.0, Arguments: []string{"1"}}, false},
		"gte-string-pass": {"gte", &ValidationOptions{Value: "a", Arguments: []string{"a"}}, true},
		"gte-string-fail": {"gte", &ValidationOptions{Value: "a", Arguments: []string{"b"}}, false},
		"lt-int-pass":     {"lt", &ValidationOptions{Value: 0, Arguments: []string{"1"}}, true},
		"lt-int-fail":     {"lt", &ValidationOptions{Value: 1, Arguments: []string{"1"}}, false},
		"lt-uint-pass":    {"lt", &ValidationOptions{Value: uint(0), Arguments: []string{"1"}}, true},
		"lt-uint-fail":    {"lt", &ValidationOptions{Value: uint(1), Arguments: []string{"1"}}, false},
		"lt-float-pass":   {"lt", &ValidationOptions{Value: 0.0, Arguments: []string{"1"}}, true},
		"lt-float-fail":   {"lt", &ValidationOptions{Value: 1.0, Arguments: []string{"1"}}, false},
		"lt-string-pass":  {"lt", &ValidationOptions{Value: "a", Arguments: []string{"b"}}, true},
		"lt-string-fail":  {"lt", &ValidationOptions{Value: "a", Arguments: []string{"a"}}, false},
		"lte-int-pass":    {"lte", &ValidationOptions{Value: 1, Arguments: []string{"1"}}, true},
		"lte-int-fail":    {"lte", &ValidationOptions{Value: 2, Arguments: []string{"1"}}, false},
		"lte-uint-pass":   {"lte", &ValidationOptions{Value: uint(1), Arguments: []string{"1"}}, true},
		"lte-uint-fail":   {"lte", &ValidationOptions{Value: uint(2), Arguments: []string{"1"}}, false},
		"lte-float-pass":  {"lte", &ValidationOptions{Value: 1.0, Arguments: []string{"1"}}, true},
		"lte-float-fail":  {"lte", &ValidationOptions{Value: 2.0, Arguments: []string{"1"}}, false},
		"lte-string-pass": {"lte", &ValidationOptions{Value: "a", Arguments: []string{"a"}}, true},
		"lte-string-fail": {"lte", &ValidationOptions{Value: "b", Arguments: []string{"a"}}, false},
		"max-int-pass":    {"max", &ValidationOptions{Value: 1, Arguments: []string{"1"}}, true},
		"max-int-fail":    {"max", &ValidationOptions{Value: 2, Arguments: []string{"1"}}, false},
		"max-uint-pass":   {"max", &ValidationOptions{Value: uint(1), Arguments: []string{"1"}}, true},
		"max-uint-fail":   {"max", &ValidationOptions{Value: uint(2), Arguments: []string{"1"}}, false},
		"max-float-pass":  {"max", &ValidationOptions{Value: 1.0, Arguments: []string{"1"}}, true},
		"max-float-fail":  {"max", &ValidationOptions{Value: 2.0, Arguments: []string{"1"}}, false},
		"max-string-pass": {"max", &ValidationOptions{Value: "ab", Arguments: []string{"2"}}, true},
		"max-string-fail": {"max", &ValidationOptions{Value: "abc", Arguments: []string{"2"}}, false},
		"max-array-pass":  {"max", &ValidationOptions{Value: []int{1, 2}, Arguments: []string{"2"}}, true},
		"max-array-fail":  {"max", &ValidationOptions{Value: []int{1, 2, 3}, Arguments: []string{"2"}}, false},
		"min-int-pass":    {"min", &ValidationOptions{Value: 1, Arguments: []string{"1"}}, true},
		"min-int-fail":    {"min", &ValidationOptions{Value: 0, Arguments: []string{"1"}}, false},
		"min-uint-pass":   {"min", &ValidationOptions{Value: uint(1), Arguments: []string{"1"}}, true},
		"min-uint-fail":   {"min", &ValidationOptions{Value: uint(0), Arguments: []string{"1"}}, false},
		"min-float-pass":  {"min", &ValidationOptions{Value: 1.0, Arguments: []string{"1"}}, true},
		"min-float-fail":  {"min", &ValidationOptions{Value: 0.0, Arguments: []string{"1"}}, false},
		"min-string-pass": {"min", &ValidationOptions{Value: "ab", Arguments: []string{"2"}}, true},
		"min-string-fail": {"min", &ValidationOptions{Value: "a", Arguments: []string{"2"}}, false},
		"min-array-pass":  {"min", &ValidationOptions{Value: []int{1, 2}, Arguments: []string{"2"}}, true},
		"min-array-fail":  {"min", &ValidationOptions{Value: []int{1}, Arguments: []string{"2"}}, false},
	}

	runTests(t, data)
}
