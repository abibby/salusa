package builder_test

import (
	"encoding/json"
	"testing"

	"github.com/abibby/salusa/database"
	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/model"
	"github.com/abibby/salusa/internal/test"
	"github.com/stretchr/testify/assert"
)

func NewTestBuilder() *builder.ModelBuilder[*test.Foo] {
	return builder.From[*test.Foo]()
}

func MustSave(tx database.DB, v model.Model) {
	err := model.Save(tx, v)
	if err != nil {
		panic(err)
	}
}

func assertJsonEqual(t *testing.T, rawJson string, v any) bool {
	b, err := json.Marshal(v)
	if !assert.NoError(t, err) {
		return false
	}
	var data any
	err = json.Unmarshal([]byte(rawJson), &data)
	if !assert.NoError(t, err) {
		return false
	}
	formattedJson, err := json.Marshal(data)
	if !assert.NoError(t, err) {
		return false
	}

	return assert.JSONEq(t, string(formattedJson), string(b))
}
