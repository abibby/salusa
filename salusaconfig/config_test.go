package salusaconfig_test

import (
	"testing"

	"github.com/abibby/salusa/salusaconfig"
)

// mockConfig implements salusaconfig.Config.
type mockConfig struct{}

func (mockConfig) GetHTTPPort() int    { return 8080 }
func (mockConfig) GetBaseURL() string  { return "https://example.com" }

var _ salusaconfig.Config = (*mockConfig)(nil)

func TestConfigInterface(t *testing.T) {
	var c salusaconfig.Config = &mockConfig{}
	if c.GetHTTPPort() != 8080 {
		t.Fatal("expected port 8080")
	}
	if c.GetBaseURL() != "https://example.com" {
		t.Fatal("unexpected base url")
	}
}
