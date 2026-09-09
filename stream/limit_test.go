package stream

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestStream_Limit(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Limit(2).
		Slice()
	assert.Equal(t, []int{1, 2}, s)
}

func TestStream_Limit_Zero(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Limit(0).
		Slice()
	assert.Empty(t, s)
}

func TestStream_Limit_GreaterThanLength(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Limit(10).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Limit_EqualsLength(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Limit(3).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Limit_One(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Limit(1).
		Slice()
	assert.Equal(t, []int{1}, s)
}

func TestStream_Skip(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Skip(1).
		Slice()
	assert.Equal(t, []int{2, 3}, s)
}

func TestStream_Skip_Zero(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Skip(0).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Skip_All(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Skip(3).
		Slice()
	assert.Empty(t, s)
}

func TestStream_Skip_GreaterThanLength(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Skip(10).
		Slice()
	assert.Empty(t, s)
}

func TestStream_Skip_Empty(t *testing.T) {
	s := Of([]int{}).
		Skip(1).
		Slice()
	assert.Empty(t, s)
}
