package stream

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestStream_Reduce(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Reduce(func(total, i int) int {
			return total + i
		})
	assert.Equal(t, 6, s)
}

func TestStream_Reduce_Empty(t *testing.T) {
	s := Of([]int{}).
		Reduce(func(total, i int) int {
			return total + i
		})
	assert.Equal(t, 0, s)
}

func TestStream_Reduce_SingleElement(t *testing.T) {
	s := Of([]int{42}).
		Reduce(func(total, i int) int {
			return total + i
		})
	assert.Equal(t, 42, s)
}

func TestStream_Reduce_Concat(t *testing.T) {
	s := Of([]string{"a", "b", "c"}).
		Reduce(func(total, i string) string {
			return total + i
		})
	assert.Equal(t, "abc", s)
}

func TestStream_Chain_MapFilterReduce(t *testing.T) {
	s := Of([]int{1, 2, 3, 4, 5, 6}).
		Filter(func(i int) bool {
			return i%2 == 0
		}).
		Map(func(i int) int {
			return i * i
		}).
		Reduce(func(total, i int) int {
			return total + i
		})
	assert.Equal(t, 56, s) // 4 + 16 + 36
}
