package builder_test

import (
	"encoding/json"
	"testing"

	"github.com/abibby/salusa/database/builder"
	"github.com/abibby/salusa/database/insert"
	"github.com/abibby/salusa/database/internal/helpers"
	"github.com/abibby/salusa/database/internal/test"
	"github.com/abibby/salusa/database/models"
	"github.com/stretchr/testify/assert"
)

func NewTestBuilder() *builder.Builder[*test.Foo] {
	return builder.From[*test.Foo]()
}

func MustSave(tx helpers.QueryExecer, v models.Model) {
	err := insert.Save(tx, v)
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
