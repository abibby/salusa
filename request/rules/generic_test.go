package rules

import (
	"testing"
)

func TestGeneric(t *testing.T) {
	data := map[string]TestCase{
		"gt-int-pass":          {"gt", &ValidationOptions{Value: 2, Arguments: []string{"1"}}, true},
		"gt-int-fail":          {"gt", &ValidationOptions{Value: 1, Arguments: []string{"1"}}, false},
		"gt-uint-pass":         {"gt", &ValidationOptions{Value: uint(2), Arguments: []string{"1"}}, true},
		"gt-uint-fail":         {"gt", &ValidationOptions{Value: uint(1), Arguments: []string{"2"}}, false},
		"gt-float-pass":        {"gt", &ValidationOptions{Value: 2.5, Arguments: []string{"1.5"}}, true},
		"gt-float-fail":        {"gt", &ValidationOptions{Value: 1.5, Arguments: []string{"1.5"}}, false},
		"gt-string-pass":       {"gt", &ValidationOptions{Value: "b", Arguments: []string{"a"}}, true},
		"gt-string-fail":       {"gt", &ValidationOptions{Value: "a", Arguments: []string{"b"}}, false},
		"gte-pass":             {"gte", &ValidationOptions{Value: 1, Arguments: []string{"1"}}, true},
		"gte-fail":             {"gte", &ValidationOptions{Value: 0, Arguments: []string{"1"}}, false},
		"lt-pass":              {"lt", &ValidationOptions{Value: 1, Arguments: []string{"2"}}, true},
		"lt-fail":              {"lt", &ValidationOptions{Value: 2, Arguments: []string{"2"}}, false},
		"lte-pass":             {"lte", &ValidationOptions{Value: 2, Arguments: []string{"2"}}, true},
		"lte-fail":             {"lte", &ValidationOptions{Value: 3, Arguments: []string{"2"}}, false},
		"max-uint-pass":        {"max", &ValidationOptions{Value: uint(1), Arguments: []string{"2"}}, true},
		"max-uint-fail":        {"max", &ValidationOptions{Value: uint(3), Arguments: []string{"2"}}, false},
		"max-float-pass":       {"max", &ValidationOptions{Value: 1.5, Arguments: []string{"2"}}, true},
		"max-float-fail":       {"max", &ValidationOptions{Value: 2.5, Arguments: []string{"2"}}, false},
		"max-array-pass":       {"max", &ValidationOptions{Value: []int{1}, Arguments: []string{"2"}}, true},
		"max-array-fail":       {"max", &ValidationOptions{Value: []int{1, 2, 3}, Arguments: []string{"2"}}, false},
		"min-uint-pass":        {"min", &ValidationOptions{Value: uint(3), Arguments: []string{"2"}}, true},
		"min-uint-fail":        {"min", &ValidationOptions{Value: uint(1), Arguments: []string{"2"}}, false},
		"min-float-pass":       {"min", &ValidationOptions{Value: 2.5, Arguments: []string{"2"}}, true},
		"min-float-fail":       {"min", &ValidationOptions{Value: 1.5, Arguments: []string{"2"}}, false},
		"min-array-pass":       {"min", &ValidationOptions{Value: []int{1, 2, 3}, Arguments: []string{"2"}}, true},
		"min-array-fail":       {"min", &ValidationOptions{Value: []int{1}, Arguments: []string{"2"}}, false},
		"multiple_of-uint-pass": {"multiple_of", &ValidationOptions{Value: uint(10), Arguments: []string{"5"}}, true},
		"multiple_of-uint-fail": {"multiple_of", &ValidationOptions{Value: uint(6), Arguments: []string{"5"}}, false},
		"multiple_of-float-pass": {"multiple_of", &ValidationOptions{Value: 10.0, Arguments: []string{"5"}}, true},
		"multiple_of-float-fail": {"multiple_of", &ValidationOptions{Value: 6.5, Arguments: []string{"5"}}, false},
	}

	runTests(t, data)
}
