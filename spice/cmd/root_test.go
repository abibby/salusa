package cmd

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	rootCmd.SetArgs([]string{"--help"})
	Execute()
}

func TestExecuteError(t *testing.T) {
	if os.Getenv("SPICE_CMD_EXECUTE_ERROR") == "1" {
		rootCmd.SetArgs([]string{"nonexistent-subcommand"})
		Execute()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExecuteError")
	cmd.Env = append(os.Environ(), "SPICE_CMD_EXECUTE_ERROR=1")
	err := cmd.Run()

	var exitErr *exec.ExitError
	if assert.ErrorAs(t, err, &exitErr) {
		assert.Equal(t, 1, exitErr.ExitCode())
	}
}
