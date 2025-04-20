package nulls

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	n := New(10)
	assert.Equal(t, 10, n.Val())
	assert.True(t, n.valid)
	v, ok := n.Ok()
	assert.Equal(t, 10, v)
	assert.True(t, ok)
}

func TestNull(t *testing.T) {
	n := Null[int]()
	assert.Equal(t, 0, n.Val()) // zero value of int
	assert.False(t, n.valid)
	v, ok := n.Ok()
	assert.Equal(t, 0, v) // zero value of int
	assert.False(t, ok)
}

func TestString(t *testing.T) {
	n := New(10)
	assert.Equal(t, "10", n.String())

	nullN := Null[int]()
	assert.Equal(t, "", nullN.String())
}

func TestIsNull(t *testing.T) {
	n := New(10)
	assert.False(t, n.IsNull())

	nullN := Null[int]()
	assert.True(t, nullN.IsNull())
}

func TestMarshalJSON(t *testing.T) {
	n := New(10)
	b, err := json.Marshal(n)
	assert.NoError(t, err)
	assert.Equal(t, "10", string(b))

	nullN := Null[int]()
	b, err = json.Marshal(nullN)
	assert.NoError(t, err)
	assert.Equal(t, "null", string(b))

	type TestStruct struct {
		Value Nullable[int] `json:"value"`
	}

	testStruct := TestStruct{Value: New(5)}
	b, err = json.Marshal(testStruct)
	assert.NoError(t, err)
	assert.Equal(t, `{"value":5}`, string(b))

	testStructNull := TestStruct{Value: Null[int]()}
	b, err = json.Marshal(testStructNull)
	assert.NoError(t, err)
	assert.Equal(t, `{"value":null}`, string(b))

	testString := New("test")
	b, err = json.Marshal(testString)
	assert.NoError(t, err)
	assert.Equal(t, `"test"`, string(b))

	testStringNull := Null[string]()
	b, err = json.Marshal(testStringNull)
	assert.NoError(t, err)
	assert.Equal(t, `null`, string(b))
}

func TestUnmarshalJSON(t *testing.T) {
	var n Nullable[int]
	err := json.Unmarshal([]byte("10"), &n)
	assert.NoError(t, err)
	assert.Equal(t, 10, n.Val())
	assert.True(t, n.valid)

	err = json.Unmarshal([]byte("null"), &n)
	assert.NoError(t, err)
	assert.False(t, n.valid)

	type TestStruct struct {
		Value Nullable[int] `json:"value"`
	}

	var testStruct TestStruct
	err = json.Unmarshal([]byte(`{"value":5}`), &testStruct)
	assert.NoError(t, err)
	assert.Equal(t, 5, testStruct.Value.Val())
	assert.True(t, testStruct.Value.valid)

	var testStructNull TestStruct
	err = json.Unmarshal([]byte(`{"value":null}`), &testStructNull)
	assert.NoError(t, err)
	assert.False(t, testStructNull.Value.valid)

	var stringNullable Nullable[string]
	err = json.Unmarshal([]byte(`"test"`), &stringNullable)
	assert.NoError(t, err)
	assert.Equal(t, "test", stringNullable.Val())
	assert.True(t, stringNullable.valid)

	var stringNullableNull Nullable[string]
	err = json.Unmarshal([]byte(`null`), &stringNullableNull)
	assert.NoError(t, err)
	assert.False(t, stringNullableNull.valid)
}

func TestNullable_Scan(t *testing.T) {
	var n Nullable[int]

	err := n.Scan(int64(10))
	assert.NoError(t, err)
	assert.Equal(t, 10, n.Val())
	assert.True(t, n.valid)
	//
	err = n.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, n.valid)

	err = n.Scan("invalid")
	assert.Error(t, err)

	var ns Nullable[string]

	err = ns.Scan("test")
	assert.NoError(t, err)
	assert.Equal(t, "test", ns.Val())
	assert.True(t, ns.valid)

	err = ns.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, ns.valid)

	err = ns.Scan(int64(123))
	assert.Error(t, err)

	var nb Nullable[bool]

	err = nb.Scan(true)
	assert.NoError(t, err)
	assert.Equal(t, true, nb.Val())
	assert.True(t, nb.valid)

	err = nb.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, nb.valid)

	err = nb.Scan("invalid")
	assert.Error(t, err)

	var nf Nullable[float32]

	err = nf.Scan(float64(1.23))
	assert.NoError(t, err)
	assert.Equal(t, float32(1.23), nf.Val())
	assert.True(t, nf.valid)

	err = nf.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, nf.valid)

	err = nf.Scan("invalid")
	assert.Error(t, err)
}

func TestNullable_Value(t *testing.T) {
	n := New(10)
	val, err := n.Value()
	assert.NoError(t, err)
	assert.Equal(t, 10, val)

	nullN := Null[int]()
	val, err = nullN.Value()
	assert.NoError(t, err)
	assert.Equal(t, nil, val)

	ns := New("test")
	val, err = ns.Value()
	assert.NoError(t, err)
	assert.Equal(t, "test", val)

	nullNs := Null[string]()
	val, err = nullNs.Value()
	assert.NoError(t, err)
	assert.Equal(t, nil, val)

	nb := New(true)
	val, err = nb.Value()
	assert.NoError(t, err)
	assert.Equal(t, true, val)

	nullNb := Null[bool]()
	val, err = nullNb.Value()
	assert.NoError(t, err)
	assert.Equal(t, nil, val)

	nf := New(1.23)
	val, err = nf.Value()
	assert.NoError(t, err)
	assert.Equal(t, 1.23, val)

	nullNf := Null[float64]()
	val, err = nullNf.Value()
	assert.NoError(t, err)
	assert.Equal(t, nil, val)
}
