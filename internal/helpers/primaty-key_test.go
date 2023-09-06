package helpers_test

import (
	"testing"

	"github.com/abibby/salusa/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestPrimaryKey(t *testing.T) {

	type NoTag struct {
		NoTag int
	}
	type WithTag struct {
		Primary    int `db:"with_tag"`
		NotPrimary int `db:"not_primary"`
	}
	type WithTagAndPrimary struct {
		NotPrimary int `db:"not_primary"`
		Primary    int `db:"with_tag_and_primary,primary"`
	}
	type Composite struct {
		NotPrimary int `db:"not_primary"`
		Primary1   int `db:"primary1,primary"`
		Primary2   int `db:"primary2,primary"`
	}

	testCases := []struct {
		name               string
		model              any
		expectedPrimaryKey []string
	}{
		{"No Tag", &NoTag{}, []string{"NoTag"}},
		{"With Tag", &WithTag{}, []string{"with_tag"}},
		{"With Tag And Primary", &WithTagAndPrimary{}, []string{"with_tag_and_primary"}},
		{"Composite", &Composite{}, []string{"primary1", "primary2"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedPrimaryKey, helpers.PrimaryKey(tc.model))
		})
	}
}
