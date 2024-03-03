package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTable(t *testing.T) {
	type Simple struct{}
	type TwoWords struct{}
	type local struct{}
	type Already_Snake struct{}
	type Singular struct{}
	type Plurals struct{}

	testCases := []struct {
		model        any
		expectedName string
	}{
		{Simple{}, "simples"},
		{TwoWords{}, "two_words"},
		{local{}, "locals"},
		{Already_Snake{}, "already_snakes"},
		{Singular{}, "singulars"},
		{Plurals{}, "plurals"},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedName, func(t *testing.T) {
			assert.Equal(t, tc.expectedName, GetTable(tc.model))
		})
	}
}

func TestGetTableSingular(t *testing.T) {
	type Simple struct{}
	type TwoWords struct{}
	type local struct{}
	type Already_Snake struct{}
	type Singular struct{}
	type Plurals struct{}

	testCases := []struct {
		model        any
		expectedName string
	}{
		{Simple{}, "simple"},
		{TwoWords{}, "two_word"},
		{local{}, "local"},
		{Already_Snake{}, "already_snake"},
		{Singular{}, "singular"},
		{Plurals{}, "plural"},
	}

	for _, tc := range testCases {
		t.Run(tc.expectedName, func(t *testing.T) {
			assert.Equal(t, tc.expectedName, GetTableSingular(tc.model))
		})
	}
}
