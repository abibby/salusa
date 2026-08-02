package relationship

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type trackingRel struct {
	initialized bool
	err         error
	self        any
	field       reflect.StructField
}

func (r *trackingRel) Initialize(self any, field reflect.StructField) error {
	r.initialized = true
	r.self = self
	r.field = field
	return r.err
}
func (r *trackingRel) Loaded() bool { return r.initialized }

type valueRel struct {
	called bool
}

func (v *valueRel) Initialize(self any, field reflect.StructField) error {
	v.called = true
	return nil
}
func (v *valueRel) Loaded() bool { return v.called }

type errorRel struct{}

func (e errorRel) Initialize(self any, field reflect.StructField) error {
	return errors.New("boom")
}
func (e errorRel) Loaded() bool { return false }

type Embedded struct {
	Rel *trackingRel
}

type withRelationships struct {
	Embedded
	Ptr   *trackingRel
	Value *valueRel
	Plain int
}

type alreadySet struct {
	Rel trackingRel
}

func TestInitializeRelationships(t *testing.T) {
	t.Run("pointer parent", func(t *testing.T) {
		m := &withRelationships{}
		err := InitializeRelationships(m)
		assert.NoError(t, err)
		assert.True(t, m.Ptr.initialized)
		assert.True(t, m.Value.called)
		assert.Equal(t, m.Ptr, m.Ptr.self.(*withRelationships).Ptr)
		assert.NotNil(t, m.Rel.self)
	})

	t.Run("value parent", func(t *testing.T) {
		m := withRelationships{}
		err := InitializeRelationships(&m)
		assert.NoError(t, err)
		assert.True(t, m.Ptr.initialized)
		assert.True(t, m.Value.called)
	})

	t.Run("slice of models", func(t *testing.T) {
		ms := []*withRelationships{{}, {}}
		err := InitializeRelationships(ms)
		assert.NoError(t, err)
		assert.True(t, ms[0].Ptr.initialized)
		assert.True(t, ms[1].Ptr.initialized)
	})

	t.Run("already initialized field is skipped", func(t *testing.T) {
		m := &alreadySet{}
		initialized := &trackingRel{}
		m.Rel = *initialized
		err := InitializeRelationships(m)
		assert.NoError(t, err)
		assert.False(t, m.Rel.initialized, "existing relationship should not be re-initialized")
	})

	t.Run("relationship error propagates", func(t *testing.T) {
		m := &struct {
			Rel errorRel
		}{}
		err := InitializeRelationships(m)
		assert.EqualError(t, err, "boom")
	})

	t.Run("plain value", func(t *testing.T) {
		err := InitializeRelationships(5)
		assert.NoError(t, err)
	})
}
