package stream

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestStream_Filter(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Filter(func(i int) bool {
			return i > 1
		}).
		Slice()
	assert.Equal(t, []int{2, 3}, s)
}

func TestStream_Filter_Empty(t *testing.T) {
	s := Of([]int{}).
		Filter(func(i int) bool {
			return true
		}).
		Slice()
	assert.Empty(t, s)
}

func TestStream_Filter_AllFiltered(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Filter(func(i int) bool {
			return i > 10
		}).
		Slice()
	assert.Empty(t, s)
}

func TestStream_Filter_NoneFiltered(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Filter(func(i int) bool {
			return true
		}).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Find_Found(t *testing.T) {
	v, ok := Of([]int{1, 2, 3}).
		Find(func(i int) bool {
			return i == 2
		})
	assert.Equal(t, 2, v)
	assert.True(t, ok)
}

func TestStream_Find_NotFound(t *testing.T) {
	v, ok := Of([]int{1, 2, 3}).
		Find(func(i int) bool {
			return i == 10
		})
	assert.Equal(t, 0, v)
	assert.False(t, ok)
}

func TestStream_Find_Empty(t *testing.T) {
	v, ok := Of([]int{}).
		Find(func(i int) bool {
			return true
		})
	assert.Equal(t, 0, v)
	assert.False(t, ok)
}

func TestStream_Find_First(t *testing.T) {
	v, ok := Of([]int{1, 2, 3}).
		Find(func(i int) bool {
			return true
		})
	assert.Equal(t, 1, v)
	assert.True(t, ok)
}
