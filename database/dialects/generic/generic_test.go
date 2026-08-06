package generic_test

import (
	"strings"

	"github.com/abibby/salusa/database/dialects"
)

type testCore struct{}

func (*testCore) Identifier(s string) string {
	if s == "*" {
		return s
	}
	parts := strings.Split(s, ".")
	for i, p := range parts {
		parts[i] = "`" + p + "`"
	}
	return strings.Join(parts, ".")
}

func (*testCore) DataType(t dialects.DataType) string {
	return string(t)
}

func (*testCore) CurrentTime() string {
	return "CURRENT_TIMESTAMP"
}

func (*testCore) AutoIncrement() string {
	return "PRIMARY KEY AUTO_INCREMENT"
}

func (*testCore) Escape(v any) string {
	return ""
}

func (*testCore) Binding() string {
	return "?"
}

func (*testCore) Features() dialects.Features {
	return dialects.Features{}
}
