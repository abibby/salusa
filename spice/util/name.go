package util

import (
	"fmt"
	"strings"
	"time"
)

func MigrationName(args []string) string {
	return fmt.Sprintf("%s-%s", time.Now().Format("20060102_150405"), strings.ReplaceAll(strings.Join(args, "_"), " ", "_"))
}
