package main

import (
	"os"
	"testing"
)

func TestMainFunc(t *testing.T) {
	oldArgs := os.Args
	os.Args = []string{"spice"}
	t.Cleanup(func() { os.Args = oldArgs })

	main()
}
