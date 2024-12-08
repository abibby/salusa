package maps

import (
	"iter"
	"sync"
)

type Sync[K, V any] struct {
	m sync.Map
}

var _ Map[any, any] = (*Sync[any, any])(nil)

// Clear deletes all the entries, resulting in an empty Map.
func (m *Sync[K, V]) Clear() {
	m.m.Clear()
}

// CompareAndDelete deletes the entry for key if its value is equal to old. The
// old value must be of a comparable type.
//
// If there is no current value for key in the map, CompareAndDelete returns
// false (even if the old value is the nil interface value).
func (m *Sync[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.m.CompareAndDelete(key, old)
}

// CompareAndSwap swaps the old and new values for key if the value stored in
// the map is equal to old. The old value must be of a comparable type.
func (m *Sync[K, V]) CompareAndSwap(key K, old V, new V) (swapped bool) {
	return m.m.CompareAndSwap(key, old, new)
}

// Delete deletes the value for a key.
func (m *Sync[K, V]) Delete(key K) {
	m.m.Delete(key)
}

// Load returns the value stored in the map for a key, or nil if no value is
// present. The ok result indicates whether value was found in the map.
func (m *Sync[K, V]) Load(key K) (value V, ok bool) {
	iv, ok := m.m.Load(key)
	v, _ := iv.(V)
	return v, ok
}

// LoadAndDelete deletes the value for a key, returning the previous value if
// any. The loaded result reports whether the key was present.
func (m *Sync[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	iv, ok := m.m.LoadAndDelete(key)
	v, _ := iv.(V)
	return v, ok
}

// LoadOrStore returns the existing value for the key if present. Otherwise, it
// stores and returns the given value. The loaded result is true if the value
// was loaded, false if stored.
func (m *Sync[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	iv, ok := m.m.LoadOrStore(key, value)
	v, _ := iv.(V)
	return v, ok
}

// Range calls f sequentially for each key and value present in the map. If f
// returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *Sync[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

// Store sets the value for a key.
func (m *Sync[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

// Swap swaps the value for a key and returns the previous value if any. The
// loaded result reports whether the key was present.
func (m *Sync[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	iv, ok := m.m.Swap(key, value)
	v, _ := iv.(V)
	return v, ok
}

func (m *Sync[K, V]) Set(key K, value V) {
	m.Store(key, value)
}
func (m *Sync[K, V]) Get(key K) (V, bool) {
	return m.Load(key)
}

func (m *Sync[K, V]) Remove(key K) {
	m.m.Delete(key)
}

func (m *Sync[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		m.Range(yield)
	}
}
