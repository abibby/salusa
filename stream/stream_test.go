package stream

import (
	"iter"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
)

func TestOf(t *testing.T) {
	var seq iter.Seq[int] = func(yield func(int) bool) {
		yield(10)
		yield(20)
		yield(30)
	}
	s := Of(seq).Slice()
	assert.Equal(t, []int{10, 20, 30}, s)
}

func TestOfSlice(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestOfSlice_Empty(t *testing.T) {
	s := OfSlice([]int{}).Slice()
	assert.Empty(t, s)
}

func TestStream_All(t *testing.T) {
	var collected []int
	for v := range OfSlice([]int{1, 2, 3}).All() {
		collected = append(collected, v)
	}
	assert.Equal(t, []int{1, 2, 3}, collected)
}

func TestStream_Slice(t *testing.T) {
	s := OfSlice([]int{4, 5, 6}).Slice()
	assert.Equal(t, []int{4, 5, 6}, s)
}

func TestStream_Filter(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Filter(func(i int) bool {
			return i > 1
		}).
		Slice()
	assert.Equal(t, []int{2, 3}, s)
}

func TestStream_Filter_Empty(t *testing.T) {
	s := OfSlice([]int{}).
		Filter(func(i int) bool {
			return true
		}).
		Slice()
	assert.Empty(t, s)
}

func TestStream_Filter_AllFiltered(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Filter(func(i int) bool {
			return i > 10
		}).
		Slice()
	assert.Empty(t, s)
}

func TestStream_Filter_NoneFiltered(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Filter(func(i int) bool {
			return true
		}).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Limit(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Limit(2).
		Slice()
	assert.Equal(t, []int{1, 2}, s)
}

func TestStream_Limit_Zero(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Limit(0).
		Slice()
	assert.Empty(t, s)
}

func TestStream_Limit_GreaterThanLength(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Limit(10).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Limit_EqualsLength(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Limit(3).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Limit_One(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Limit(1).
		Slice()
	assert.Equal(t, []int{1}, s)
}

func TestStream_Skip(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Skip(1).
		Slice()
	assert.Equal(t, []int{2, 3}, s)
}

func TestStream_Skip_Zero(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Skip(0).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Skip_All(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Skip(3).
		Slice()
	assert.Empty(t, s)
}

func TestStream_Skip_GreaterThanLength(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Skip(10).
		Slice()
	assert.Empty(t, s)
}

func TestStream_Skip_Empty(t *testing.T) {
	s := OfSlice([]int{}).
		Skip(1).
		Slice()
	assert.Empty(t, s)
}

func TestStream_Find_Found(t *testing.T) {
	v, ok := OfSlice([]int{1, 2, 3}).
		Find(func(i int) bool {
			return i == 2
		})
	assert.Equal(t, 2, v)
	assert.True(t, ok)
}

func TestStream_Find_NotFound(t *testing.T) {
	v, ok := OfSlice([]int{1, 2, 3}).
		Find(func(i int) bool {
			return i == 10
		})
	assert.Equal(t, 0, v)
	assert.False(t, ok)
}

func TestStream_Find_Empty(t *testing.T) {
	v, ok := OfSlice([]int{}).
		Find(func(i int) bool {
			return true
		})
	assert.Equal(t, 0, v)
	assert.False(t, ok)
}

func TestStream_Find_First(t *testing.T) {
	v, ok := OfSlice([]int{1, 2, 3}).
		Find(func(i int) bool {
			return true
		})
	assert.Equal(t, 1, v)
	assert.True(t, ok)
}

func TestStream_Sort(t *testing.T) {
	s := OfSlice([]int{3, 1, 2}).
		Sort(func(a, b int) int {
			return a - b
		}).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Sort_AlreadySorted(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Sort(func(a, b int) int {
			return a - b
		}).
		Slice()
	assert.Equal(t, []int{1, 2, 3}, s)
}

func TestStream_Sort_Descending(t *testing.T) {
	s := OfSlice([]int{1, 2, 3}).
		Sort(func(a, b int) int {
			return b - a
		}).
		Slice()
	assert.Equal(t, []int{3, 2, 1}, s)
}

func TestStream_Sort_SingleElement(t *testing.T) {
	s := OfSlice([]int{42}).
		Sort(func(a, b int) int {
			return a - b
		}).
		Slice()
	assert.Equal(t, []int{42}, s)
}

func TestStream_Chain(t *testing.T) {
	s := OfSlice([]int{1, 2, 3, 4, 5, 6}).
		Filter(func(i int) bool {
			return i%2 == 0
		}).
		Limit(2).
		Slice()
	assert.Equal(t, []int{2, 4}, s)
}

func TestStream_Chain_FilterSkipMap(t *testing.T) {
	s := OfSlice([]int{1, 2, 3, 4, 5}).
		Filter(func(i int) bool {
			return i > 2
		}).
		Skip(1).
		Map(func(i int) string {
			return "v"
		}).
		Slice()
	assert.Equal(t, []string{"v", "v"}, s)
}
