package sets

import (
	"cmp"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Factory[T cmp.Ordered] struct {
	name string
	fn   FactoryFunc[T]
}

type FactoryFunc[T cmp.Ordered] func(values ...T) Set[T]

func makeFactories[T cmp.Ordered]() []Factory[T] {
	return []Factory[T]{
		{
			name: "Map",
			fn: func(values ...T) Set[T] {
				return NewMapSet(values...)
			},
		},
		{
			name: "Slice",
			fn: func(values ...T) Set[T] {
				return NewSliceSet(values...)
			},
		},
		{
			name: "Ordered",
			fn: func(values ...T) Set[T] {
				return NewOrderedSet(values...)
			},
		},
		{
			name: "Default",
			fn: func(values ...T) Set[T] {
				return New(values...)
			},
		},
	}
}

func TestNew(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		expectedLen int
	}{
		{
			name:        "one value",
			args:        []string{"a"},
			expectedLen: 1,
		},
		{
			name:        "zero values",
			args:        []string{},
			expectedLen: 0,
		},
		{
			name:        "two value",
			args:        []string{"a", "b"},
			expectedLen: 2,
		},
		{
			name:        "duplicate",
			args:        []string{"a", "a"},
			expectedLen: 1,
		},
	}
	for _, fac := range makeFactories[string]() {
		for _, tc := range testCases {
			t.Run(tc.name+" "+fac.name, func(t *testing.T) {
				s := fac.fn(tc.args...)
				assert.Equal(t, tc.expectedLen, s.Len())
				for _, arg := range tc.args {
					assert.True(t, s.Has(arg))
				}
			})
		}
	}
}
func TestAdd(t *testing.T) {
	testCases := []struct {
		name          string
		initialValues []string
		args          []string
		expectedLen   int
	}{
		{
			name:        "zero values",
			args:        []string{},
			expectedLen: 0,
		},
		{
			name:          "one value",
			initialValues: []string{"a"},
			args:          []string{"b"},
			expectedLen:   2,
		},
		{
			name:          "existing value",
			initialValues: []string{"a"},
			args:          []string{"a"},
			expectedLen:   1,
		},
		{
			name:          "duplicate value",
			initialValues: []string{"a"},
			args:          []string{"b", "b"},
			expectedLen:   2,
		},
	}
	for _, fac := range makeFactories[string]() {
		for _, tc := range testCases {
			t.Run(tc.name+" "+fac.name, func(t *testing.T) {
				s := fac.fn(tc.initialValues...)
				s.Add(tc.args...)
				assert.Equal(t, tc.expectedLen, s.Len())
				for _, arg := range tc.args {
					assert.True(t, s.Has(arg))
				}
			})
		}
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name           string
		initialValues  []string
		args           []string
		expectedValues []string
	}{
		{
			name:           "zero values",
			args:           []string{},
			expectedValues: []string{},
		},
		{
			name:           "one value",
			initialValues:  []string{"a", "b"},
			args:           []string{"b"},
			expectedValues: []string{"a"},
		},
		{
			name:           "missing value",
			initialValues:  []string{"a", "b"},
			args:           []string{"c"},
			expectedValues: []string{"a", "b"},
		},
		{
			name:           "multiple values",
			initialValues:  []string{"a", "b"},
			args:           []string{"a", "b"},
			expectedValues: []string{},
		},
	}
	for _, fac := range makeFactories[string]() {
		for _, tc := range testCases {
			t.Run(tc.name+" "+fac.name, func(t *testing.T) {
				s := fac.fn(tc.initialValues...)
				s.Delete(tc.args...)
				assert.Equal(t, len(tc.expectedValues), s.Len())
				for _, arg := range tc.expectedValues {
					assert.True(t, s.Has(arg))
				}
			})
		}
	}
}

func TestClone(t *testing.T) {
	for _, fac := range makeFactories[string]() {
		t.Run(fac.name, func(t *testing.T) {
			s := fac.fn("a", "b", "c")
			c := s.Clone()
			assert.Equal(t, s.Len(), c.Len())
			for _, v := range []string{"a", "b", "c"} {
				assert.True(t, c.Has(v))
			}
			c.Add("d")
			assert.False(t, s.Has("d"))
			assert.Equal(t, 3, s.Len())
		})
	}
}

func TestAll(t *testing.T) {
	for _, fac := range makeFactories[string]() {
		t.Run(fac.name, func(t *testing.T) {
			s := fac.fn("a", "b", "c")
			got := map[string]bool{}
			for v := range s.All() {
				got[v] = true
			}
			assert.Equal(t, map[string]bool{"a": true, "b": true, "c": true}, got)
		})
	}
}

func TestAllEarlyExit(t *testing.T) {
	for _, fac := range makeFactories[string]() {
		t.Run(fac.name, func(t *testing.T) {
			s := fac.fn("a", "b", "c")
			count := 0
			for range s.All() {
				count++
				if count == 1 {
					break
				}
			}
			assert.Equal(t, 1, count)
		})
	}
}

func TestGet(t *testing.T) {
	t.Run("Slice", func(t *testing.T) {
		s := NewSliceSet("a", "b", "c")
		assert.Equal(t, "a", s.Get(0))
		assert.Equal(t, "c", s.Get(2))
	})
	t.Run("Ordered", func(t *testing.T) {
		s := NewOrderedSet("a", "b", "c")
		assert.Equal(t, "a", s.Get(0))
		assert.Equal(t, "c", s.Get(2))
	})
}
