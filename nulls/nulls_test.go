// nulls_test.go
package nulls

import (
	"encoding/json"
	"reflect"
	"testing"
	"time" // Import time for testing with time.Time

	"github.com/stretchr/testify/assert"
)

// --- Test Cases ---

func TestNew(t *testing.T) {
	t.Run("Int", func(t *testing.T) {
		val := 42
		n := New(val)
		if !n.Valid {
			t.Errorf("New(%d): expected Valid=true, got false", val)
		}
		if n.V != val {
			t.Errorf("New(%d): expected V=%d, got %d", val, val, n.V)
		}
	})

	t.Run("String", func(t *testing.T) {
		val := "hello"
		n := New(val)
		if !n.Valid {
			t.Errorf("New(%q): expected Valid=true, got false", val)
		}
		if n.V != val {
			t.Errorf("New(%q): expected V=%q, got %q", val, val, n.V)
		}
	})

	t.Run("Struct", func(t *testing.T) {
		type testStruct struct {
			Field string
		}
		val := testStruct{Field: "world"}
		n := New(val)
		if !n.Valid {
			t.Errorf("New(%v): expected Valid=true, got false", val)
		}
		if n.V != val {
			t.Errorf("New(%v): expected V=%v, got %v", val, val, n.V)
		}
	})
}

func TestNullJSON(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name         string
		value        any // The Null[T] value itself
		expectedJSON string
		expectErr    bool // For marshaling errors (rare here)
	}{
		{"Valid Int", New(123), `123`, false},
		{"Invalid Int", Null[int]{}, `null`, false},
		{"Valid String", New("hello world"), `"hello world"`, false},
		{"Invalid String", Null[string]{}, `null`, false},
		{"Valid Bool True", New(true), `true`, false},
		{"Valid Bool False", New(false), `false`, false},
		{"Invalid Bool", Null[bool]{}, `null`, false},
		{"Valid Float", New(3.14), `3.14`, false},
		{"Invalid Float", Null[float64]{}, `null`, false},
		{"Valid Struct", New(testStruct{Name: "Bob", Age: 30}), `{"name":"Bob","age":30}`, false},
		{"Invalid Struct", Null[testStruct]{}, `null`, false},
		{"Valid Time", New(time.Date(2023, 10, 26, 12, 0, 0, 0, time.UTC)), `"2023-10-26T12:00:00Z"`, false},
		{"Invalid Time", Null[time.Time]{}, `null`, false},
		{"Valid Slice", New([]int{1, 2, 3}), `[1,2,3]`, false},
		{"Invalid Slice", Null[[]int]{}, `null`, false},
		{"Valid Pointer", New(func() *int { i := 5; return &i }()), `5`, false}, // Marshals the pointed-to value
		{"Invalid Pointer", Null[*int]{}, `null`, false},                        // Invalid Null pointer
		{"Valid Nil Pointer", New[*int](nil), `null`, false},                    // Valid Null containing a nil pointer
	}

	for _, tt := range tests {
		t.Run("Marshal/"+tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.value)

			if tt.expectErr {
				assert.Error(t, err, "json.Marshal")
			} else {
				assert.NoError(t, err, "json.Marshal")
				if string(b) != tt.expectedJSON {
					t.Errorf("json.Marshal: expected %s, got %s", tt.expectedJSON, string(b))
				}
			}
		})
	}

	unmarshalTests := []struct {
		name        string
		jsonInput   string
		targetPtr   any // Pointer to the Null[T] to unmarshal into
		expectedVal any // The expected Null[T] value after unmarshal
		expectErr   bool
	}{
		{"Valid Int", `456`, new(Null[int]), New(456), false},
		{"Null Int", `null`, new(Null[int]), Null[int]{}, false},
		{"Invalid Int", `"abc"`, new(Null[int]), Null[int]{}, true}, // Expect json unmarshal error
		{"Valid String", `"world"`, new(Null[string]), New("world"), false},
		{"Null String", `null`, new(Null[string]), Null[string]{}, false},
		{"Empty String", `""`, new(Null[string]), New(""), false},
		{"Invalid String", `123`, new(Null[string]), Null[string]{}, true}, // Expect json unmarshal error
		{"Valid Bool True", `true`, new(Null[bool]), New(true), false},
		{"Valid Bool False", `false`, new(Null[bool]), New(false), false},
		{"Null Bool", `null`, new(Null[bool]), Null[bool]{}, false},
		{"Invalid Bool", `"true"`, new(Null[bool]), Null[bool]{}, true}, // Expect json unmarshal error
		{"Valid Float", `9.81`, new(Null[float64]), New(9.81), false},
		{"Null Float", `null`, new(Null[float64]), Null[float64]{}, false},
		{"Valid Struct", `{"name":"Alice","age":25}`, new(Null[testStruct]), New(testStruct{Name: "Alice", Age: 25}), false},
		{"Null Struct", `null`, new(Null[testStruct]), Null[testStruct]{}, false},
		{"Invalid Struct Field", `{"name":"Alice","age":"twenty-five"}`, new(Null[testStruct]), nil, true}, // Expect json unmarshal error
		{"Valid Time", `"2024-01-15T10:30:00Z"`, new(Null[time.Time]), New(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)), false},
		{"Null Time", `null`, new(Null[time.Time]), Null[time.Time]{}, false},
		{"Invalid Time", `"not-a-time"`, new(Null[time.Time]), nil, true}, // Expect json unmarshal error
		{"Valid Slice", `[4, 5, 6]`, new(Null[[]int]), New([]int{4, 5, 6}), false},
		{"Null Slice", `null`, new(Null[[]int]), Null[[]int]{}, false},
		{"Valid Pointer", `10`, new(Null[*int]), New(func() *int { i := 10; return &i }()), false}, // Unmarshals into the pointed-to value
		{"Null Pointer", `null`, new(Null[*int]), Null[*int]{}, false},                             // Unmarshals null into the Null type itself
		{"Valid Null Pointer", `null`, new(Null[*int]), Null[*int]{}, false},                       // Can't unmarshal a nil *pointer* from json `null`, sets Null invalid
	}

	for _, tt := range unmarshalTests {
		t.Run("Unmarshal/"+tt.name, func(t *testing.T) {
			err := json.Unmarshal([]byte(tt.jsonInput), tt.targetPtr)

			if tt.expectErr {
				assert.Error(t, err, "json.Unmarshal")
			} else {
				assert.NoError(t, err, "json.Unmarshal")
				// Use reflection to compare the value pointed to by targetPtr with expectedVal
				targetVal := reflect.ValueOf(tt.targetPtr).Elem().Interface()
				if !reflect.DeepEqual(targetVal, tt.expectedVal) {
					// Special case for time comparison as DeepEqual might fail due to monotonic clock
					if ntTarget, okTarget := targetVal.(Null[time.Time]); okTarget {
						if ntExpected, okExpected := tt.expectedVal.(Null[time.Time]); okExpected {
							if ntTarget.Valid != ntExpected.Valid || (ntTarget.Valid && !ntTarget.V.Equal(ntExpected.V)) {
								t.Errorf("json.Unmarshal: expected %#v, got %#v", tt.expectedVal, targetVal)
							}
							return // Skip default check if time comparison passed/failed
						}
					}
					// Special case for pointer comparison (we compare pointed-to values if both are non-nil)
					if ntTarget, okTarget := targetVal.(Null[*int]); okTarget {
						if ntExpected, okExpected := tt.expectedVal.(Null[*int]); okExpected {
							if ntTarget.Valid != ntExpected.Valid {
								t.Errorf("json.Unmarshal: pointer validity mismatch: expected %#v, got %#v", tt.expectedVal, targetVal)
								return
							}
							if ntTarget.Valid { // Both valid, compare pointed-to values
								if (ntTarget.V == nil && ntExpected.V != nil) || (ntTarget.V != nil && ntExpected.V == nil) || (ntTarget.V != nil && ntExpected.V != nil && *ntTarget.V != *ntExpected.V) {
									t.Errorf("json.Unmarshal: pointer value mismatch: expected %#v, got %#v", tt.expectedVal, targetVal)
									return
								}
							}
							return // Skip default check
						}
					}
					// Default comparison
					t.Errorf("json.Unmarshal: expected %#v, got %#v", tt.expectedVal, targetVal)
				}
			}
		})
	}
}

func TestNullSQL(t *testing.T) {
	// --- Scan Tests ---
	t.Run("Scan", func(t *testing.T) {
		t.Run("Int From Int64", func(t *testing.T) {
			var ni Null[int]
			val := int64(987)
			err := ni.Scan(val)
			assert.NoError(t, err, "Scan int64")
			if !ni.Valid {
				t.Errorf("Scan(%v): Expected Valid=true, got false", val)
			}
			if ni.V != int(val) {
				t.Errorf("Scan(%v): Expected V=%d, got %d", val, int(val), ni.V)
			}
		})

		t.Run("String From String", func(t *testing.T) {
			var ns Null[string]
			val := "db string"
			err := ns.Scan(val)
			assert.NoError(t, err, "Scan string")
			if !ns.Valid {
				t.Errorf("Scan(%q): Expected Valid=true, got false", val)
			}
			if ns.V != val {
				t.Errorf("Scan(%q): Expected V=%q, got %q", val, val, ns.V)
			}
		})

		t.Run("String From Bytes", func(t *testing.T) {
			var ns Null[string]
			val := []byte("db bytes")
			expected := "db bytes"
			err := ns.Scan(val)
			assert.NoError(t, err, "Scan string from bytes")
			if !ns.Valid {
				t.Errorf("Scan(%v): Expected Valid=true, got false", val)
			}
			if ns.V != expected {
				t.Errorf("Scan(%v): Expected V=%q, got %q", val, expected, ns.V)
			}
		})

		t.Run("Time From Time", func(t *testing.T) {
			var nt Null[time.Time]
			val := time.Date(2025, 4, 22, 10, 0, 0, 0, time.UTC)
			err := nt.Scan(val)
			assert.NoError(t, err, "Scan time")
			if !nt.Valid {
				t.Errorf("Scan(%v): Expected Valid=true, got false", val)
			}
			if !nt.V.Equal(val) {
				t.Errorf("Scan(%v): Expected V=%v, got %v", val, val, nt.V)
			}
		})

		t.Run("Any From Nil", func(t *testing.T) {
			ni := New(100) // Start valid
			err := ni.Scan(nil)
			assert.NoError(t, err, "Scan nil")
			if ni.Valid {
				t.Errorf("Scan(nil): Expected Valid=false, got true")
			}
			if ni.V != 0 { // Should be zero value for int
				t.Errorf("Scan(nil): Expected V=0, got %d", ni.V)
			}

			ns := New("initial") // Start valid
			err = ns.Scan(nil)
			assert.NoError(t, err, "Scan nil")
			if ns.Valid {
				t.Errorf("Scan(nil): Expected Valid=false, got true")
			}
			if ns.V != "" { // Should be zero value for string
				t.Errorf("Scan(nil): Expected V=\"\", got %q", ns.V)
			}
		})
	})

	// --- Value Tests ---
	t.Run("Value", func(t *testing.T) {
		t.Run("Valid Int", func(t *testing.T) {
			n := New(int(55))
			v, err := n.Value()
			assert.NoError(t, err, "Value valid int")
			if dv, ok := v.(int64); !ok || dv != int64(55) { // Drivers often expect basic types
				// Check original type if direct match fails
				if val, ok := v.(int); !ok || val != 55 {
					t.Errorf("Value(): Expected driver.Value=%d (or int64), got %T(%v)", 55, v, v)
				}
			}
		})

		t.Run("Invalid Int", func(t *testing.T) {
			var n Null[int]
			v, err := n.Value()
			assert.NoError(t, err, "Value invalid int")
			if v != nil {
				t.Errorf("Value(): Expected driver.Value=nil, got %T(%v)", v, v)
			}
		})

		t.Run("Valid String", func(t *testing.T) {
			n := New("sql value")
			v, err := n.Value()
			assert.NoError(t, err, "Value valid string")
			if dv, ok := v.(string); !ok || dv != "sql value" {
				t.Errorf("Value(): Expected driver.Value=%q, got %T(%v)", "sql value", v, v)
			}
		})

		t.Run("Invalid String", func(t *testing.T) {
			var n Null[string]
			v, err := n.Value()
			assert.NoError(t, err, "Value invalid string")
			if v != nil {
				t.Errorf("Value(): Expected driver.Value=nil, got %T(%v)", v, v)
			}
		})

		t.Run("Valid Time", func(t *testing.T) {
			tval := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			n := New(tval)
			v, err := n.Value()
			assert.NoError(t, err, "Value valid time")
			if dv, ok := v.(time.Time); !ok || !dv.Equal(tval) {
				t.Errorf("Value(): Expected driver.Value=%v, got %T(%v)", tval, v, v)
			}
		})

		t.Run("Invalid Time", func(t *testing.T) {
			var n Null[time.Time]
			v, err := n.Value()
			assert.NoError(t, err, "Value invalid time")
			if v != nil {
				t.Errorf("Value(): Expected driver.Value=nil, got %T(%v)", v, v)
			}
		})

		t.Run("Invalid Struct", func(t *testing.T) {
			type testStruct struct{ V int }
			var n Null[testStruct]
			v, err := n.Value()
			assert.NoError(t, err, "Value invalid struct")
			if v != nil {
				t.Errorf("Value(): Expected driver.Value=nil, got %T(%v)", v, v)
			}
		})
	})
}
