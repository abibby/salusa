package config_test

import (
	"os"
	"testing"

	"github.com/abibby/salusa/database/dialects/sqlite"
	"github.com/abibby/salusa/email/emailtest"
	"github.com/abibby/salusa/event"
	"github.com/abibby/salusa/static/template/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Setenv("PORT", "8080")
	t.Setenv("BASE_PATH", "/base")
	t.Setenv("DATABASE_PATH", "/tmp/db.sqlite")
	t.Setenv("MAIL_FROM", "from@example.com")
	t.Setenv("MAIL_HOST", "smtp.example.com")
	t.Setenv("MAIL_PORT", "2526")
	t.Setenv("MAIL_USERNAME", "user2")
	t.Setenv("MAIL_PASSWORD", "pass2")

	c := config.Load()
	require.NotNil(t, c)

	assert.Equal(t, 8080, c.Port)
	assert.Equal(t, "/base", c.BasePath)
}

func TestConfig_Methods(t *testing.T) {
	c := &config.Config{
		Port:     8080,
		BasePath: "/base",
		Database: sqlite.NewConfig(":memory:"),
		Mail:     emailtest.NewTestMailerConfig(),
		Queue:    event.NewChannelQueueConfig(),
	}

	assert.Equal(t, 8080, c.GetHTTPPort())
	assert.Equal(t, "/base", c.GetBaseURL())
	assert.Equal(t, c.Database, c.DBConfig())
	assert.Equal(t, c.Mail, c.MailConfig())
	assert.Equal(t, c.Queue, c.QueueConfig())
}

func TestLoad_missingEnvFile(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("BASE_PATH", "")
	t.Setenv("DATABASE_PATH", "")

	c := config.Load()
	require.NotNil(t, c)
	assert.Equal(t, 2303, c.Port)
}

func TestLoad_malformedEnvFile(t *testing.T) {
	t.Chdir(t.TempDir())
	require.NoError(t, os.WriteFile(".env", []byte("NOT_A_VALID_LINE\n"), 0o644))

	assert.Panics(t, func() {
		config.Load()
	})
}
