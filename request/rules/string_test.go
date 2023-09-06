package rules

import (
	"testing"
)

func TestString(t *testing.T) {
	data := map[string]TestCase{
		"alpha-pass":          {"alpha", &ValidationOptions{Value: "a"}, true},
		"alpha-fail":          {"alpha", &ValidationOptions{Value: "1a"}, false},
		"alpha_dash-pass":     {"alpha_dash", &ValidationOptions{Value: "a-"}, true},
		"alpha_dash-fail":     {"alpha_dash", &ValidationOptions{Value: "1a-"}, false},
		"alpha_num-pass":      {"alpha_num", &ValidationOptions{Value: "a1"}, true},
		"alpha_num-fail":      {"alpha_num", &ValidationOptions{Value: "1a!"}, false},
		"numeric-pass":        {"numeric", &ValidationOptions{Value: "123"}, true},
		"numeric-fail":        {"numeric", &ValidationOptions{Value: "123a"}, false},
		"email-pass":          {"email", &ValidationOptions{Value: "user@example.com"}, true},
		"email-fail":          {"email", &ValidationOptions{Value: "not an email"}, false},
		"ends_with-pass":      {"ends_with", &ValidationOptions{Value: "stringend", Arguments: []string{"end"}}, true},
		"ends_with-fail":      {"ends_with", &ValidationOptions{Value: "stringendnot", Arguments: []string{"end"}}, false},
		"starts_with-pass":    {"starts_with", &ValidationOptions{Value: "startstring", Arguments: []string{"start"}}, true},
		"starts_with-fail":    {"starts_with", &ValidationOptions{Value: "notstartstring", Arguments: []string{"start"}}, false},
		"ip_address-pass":     {"ip_address", &ValidationOptions{Value: "192.168.0.1"}, true},
		"ip_address-fail":     {"ip_address", &ValidationOptions{Value: "192.168.0.0.0"}, false},
		"json-pass":           {"json", &ValidationOptions{Value: `{"foo":"bar"}`}, true},
		"json-fail":           {"json", &ValidationOptions{Value: `{foo:"bar"}`}, false},
		"mac_address-pass":    {"mac_address", &ValidationOptions{Value: "00:00:5e:00:53:af"}, true},
		"mac_address-fail":    {"mac_address", &ValidationOptions{Value: "mac address"}, false},
		"not_regex-pass":      {"not_regex", &ValidationOptions{Value: "match", Arguments: []string{"^match\\d+$"}}, true},
		"not_regex-fail":      {"not_regex", &ValidationOptions{Value: "match123", Arguments: []string{"^match\\d+$"}}, false},
		"regex-pass":          {"regex", &ValidationOptions{Value: "match123", Arguments: []string{"^match\\d+$"}}, true},
		"regex-fail":          {"regex", &ValidationOptions{Value: "match", Arguments: []string{"^match\\d+$"}}, false},
		"timezone-pass":       {"timezone", &ValidationOptions{Value: "America/Toronto"}, true},
		"timezone-fail":       {"timezone", &ValidationOptions{Value: "fake timezone"}, false},
		"url-pass":            {"url", &ValidationOptions{Value: "https://example.com/foo"}, true},
		"url-fail":            {"url", &ValidationOptions{Value: "fake url"}, false},
		"uuid-pass":           {"uuid", &ValidationOptions{Value: "d0ee2dc2-f392-4615-8e08-fff28bbc2234"}, true},
		"uuid-fail":           {"uuid", &ValidationOptions{Value: "fake uuid"}, false},
		"length-pass":         {"length", &ValidationOptions{Value: "short", Arguments: []string{"5"}}, true},
		"length-fail":         {"length", &ValidationOptions{Value: "a long string", Arguments: []string{"5"}}, false},
		"length_between-pass": {"length_between", &ValidationOptions{Value: "short", Arguments: []string{"4", "6"}}, true},
		"length_between-fail": {"length_between", &ValidationOptions{Value: "a long string", Arguments: []string{"5", "10"}}, false},
	}

	runTests(t, data)
}
