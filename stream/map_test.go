package stream

import (
	"fmt"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestStream_Map(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Map(func(i int) string {
			return fmt.Sprint(i)
		}).
		Slice()
	assert.Equal(t, []string{"1", "2", "3"}, s)
}

func TestStream_Map_Empty(t *testing.T) {
	s := OfSlice([]int{}).
		Map(func(i int) string {
			return fmt.Sprint(i)
		}).
		Slice()
	assert.Empty(t, s)
}

func TestStream_Map_SingleElement(t *testing.T) {
	s := OfSlice([]int{42}).
		Map(func(i int) string {
			return fmt.Sprint(i)
		}).
		Slice()
	assert.Equal(t, []string{"42"}, s)
}

func TestStream_Map_Identity(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Map(func(i int) int {
			return i
		}).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Map_ChainMap(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
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
	s := OfSlice([]int{1, 2, 3}).
		FlatMap(func(i int) []int {
			return []int{i, i}
		}).
		Slice()
	assert.Equal(t, []int{1, 1, 2, 2, 3, 3}, s)
}

func TestStream_FlatMap_Empty(t *testing.T) {
	s := OfSlice([]int{}).
		FlatMap(func(i int) []int {
			return []int{i, i}
		}).
		Slice()
	assert.Empty(t, s)
}

func TestStream_FlatMap_EmptyResults(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		FlatMap(func(i int) []int {
			return []int{}
		}).
		Slice()
	assert.Empty(t, s)
}

func TestStream_FlatMap_MixedResults(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
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
	s := OfSlice([]int{1, 2, 3}).
		FlatMap(func(i int) []int {
			return []int{i}
		}).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Reduce(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Reduce(func(total, i int) int {
			return total + i
		})
	assert.Equal(t, 6, s)
}

func TestStream_Reduce_Empty(t *testing.T) {
	s := OfSlice([]int{}).
		Reduce(func(total, i int) int {
			return total + i
		})
	assert.Equal(t, 0, s)
}

func TestStream_Reduce_SingleElement(t *testing.T) {
	s := OfSlice([]int{42}).
		Reduce(func(total, i int) int {
			return total + i
		})
	assert.Equal(t, 42, s)
}

func TestStream_Reduce_Concat(t *testing.T) {
	s := OfSlice([]string{"a", "b", "c"}).
		Reduce(func(total, i string) string {
			return total + i
		})
	assert.Equal(t, "abc", s)
}

func TestStream_Chain_MapFilterReduce(t *testing.T) {
	s := OfSlice([]int{1, 2, 3, 4, 5, 6}).
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

func TestStream_Chain_FilterMapCollect(t *testing.T) {
	s := OfSlice([]int{1, 2, 3, 4, 5}).
		Filter(func(i int) bool {
			return i >= 3
		}).
		Map(func(i int) string {
			return fmt.Sprintf("item-%d", i)
		}).
		Slice()
	assert.Equal(t, []string{"item-3", "item-4", "item-5"}, s)
}

func TestStream_Chain_FlatMapFilterLimit(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		FlatMap(func(i int) []int {
			return []int{i, i * 10}
		}).
		Filter(func(i int) bool {
			return i > 5
		}).
		Limit(2).
		Slice()
	assert.Equal(t, []int{10, 20}, s)
}
