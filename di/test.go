package di

import (
	"os"
	"testing"
)

func TestResetDefaultProvider() {
	defaultProvider = NewDependencyProvider()
}

func TestMain(m *testing.M) {
	TestResetDefaultProvider()

	code := m.Run()

	os.Exit(code)
}
