package stream

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestStream_Sort(t *testing.T) {
	s := Of([]int{3, 1, 2}).
		Sort(func(a, b int) int {
			return a - b
		}).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Sort_AlreadySorted(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Sort(func(a, b int) int {
			return a - b
		}).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Sort_Descending(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Sort(func(a, b int) int {
			return b - a
		}).
		Slice()
	assert.Equal(t, []int{3, 2, 1}, s)
}

func TestStream_Sort_SingleElement(t *testing.T) {
	s := Of([]int{42}).
		Sort(func(a, b int) int {
			return a - b
		}).
		Slice()
	assert.Equal(t, []int{42}, s)
}
