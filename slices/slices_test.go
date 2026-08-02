package slices_test

import (
	"testing"

	"github.com/abibby/salusa/slices"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	result := slices.Filter([]int{1, 2, 3, 4, 5}, func(v int) bool {
		return v%2 == 0
	})
	assert.Equal(t, []int{2, 4}, result)

	assert.Equal(t, []int{}, slices.Filter([]int{}, func(v int) bool { return true }))
}

func TestFind(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		v, ok := slices.Find([]string{"a", "b", "c"}, func(s string) bool {
			return s == "b"
		})
		assert.True(t, ok)
		assert.Equal(t, "b", v)
	})

	t.Run("not found", func(t *testing.T) {
		v, ok := slices.Find([]string{"a", "b"}, func(s string) bool {
			return s == "z"
		})
		assert.False(t, ok)
		assert.Equal(t, "", v)
	})
}

func TestMap(t *testing.T) {
	result := slices.Map([]int{1, 2, 3}, func(v int) string {
		return string(rune('a' + v - 1))
	})
	assert.Equal(t, []string{"a", "b", "c"}, result)

	assert.Equal(t, []string{}, slices.Map([]int{}, func(v int) string { return "" }))
}

type equaler struct {
	v int
}

func (e equaler) Equal(o equaler) bool {
	return e.v == o.v
}

func TestHas(t *testing.T) {
	t.Run("comparable found", func(t *testing.T) {
		assert.True(t, slices.Has([]int{1, 2, 3}, 2))
	})
	t.Run("comparable not found", func(t *testing.T) {
		assert.False(t, slices.Has([]int{1, 2, 3}, 4))
	})
	t.Run("equaler found", func(t *testing.T) {
		assert.True(t, slices.Has([]equaler{{v: 1}, {v: 2}}, equaler{v: 2}))
	})
	t.Run("equaler not found", func(t *testing.T) {
		assert.False(t, slices.Has([]equaler{{v: 1}, {v: 2}}, equaler{v: 3}))
	})
}
