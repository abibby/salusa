package stream

import (
	"fmt"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestStream_Map(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Map(func(i int) string {
			return fmt.Sprint(i)
		}).
		Slice()
	assert.Equal(t, []string{"1", "2", "3"}, s)
}

func TestStream_Map_Empty(t *testing.T) {
	s := Of([]int{}).
		Map(func(i int) string {
			return fmt.Sprint(i)
		}).
		Slice()
	assert.Empty(t, s)
}

func TestStream_Map_SingleElement(t *testing.T) {
	s := Of([]int{42}).
		Map(func(i int) string {
			return fmt.Sprint(i)
		}).
		Slice()
	assert.Equal(t, []string{"42"}, s)
}

func TestStream_Map_Identity(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Map(func(i int) int {
			return i
		}).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Map_ChainMap(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		Map(func(i int) int {
			return i * 10
		}).
		Map(func(i int) string {
			return fmt.Sprintf("#%d", i)
		}).
		Slice()
	assert.Equal(t, []string{"#10", "#20", "#30"}, s)
}

func TestStream_FlatMap(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		FlatMap(func(i int) []int {
			return []int{i, i}
		}).
		Slice()
	assert.Equal(t, []int{1, 1, 2, 2, 3, 3}, s)
}

func TestStream_FlatMap_Empty(t *testing.T) {
	s := Of([]int{}).
		FlatMap(func(i int) []int {
			return []int{i, i}
		}).
		Slice()
	assert.Empty(t, s)
}

func TestStream_FlatMap_EmptyResults(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		FlatMap(func(i int) []int {
			return []int{}
		}).
		Slice()
	assert.Empty(t, s)
}

func TestStream_FlatMap_MixedResults(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		FlatMap(func(i int) []int {
			if i%2 == 0 {
				return []int{i}
			}
			return []int{}
		}).
		Slice()
	assert.Equal(t, []int{2}, s)
}

func TestStream_FlatMap_SingleResult(t *testing.T) {
	s := Of([]int{1, 2, 3}).
		FlatMap(func(i int) []int {
			return []int{i}
		}).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}
