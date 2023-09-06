package postgres

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/abibby/salusa/database/dialects"
)

type Posgtgres struct {
	bindingNumber int
}

func (*Posgtgres) Identifier(s string) string {
	parts := strings.Split(s, ".")
	for i, p := range parts {
		parts[i] = "\"" + p + "\""
	}
	return strings.Join(parts, ".")
}

func (*Posgtgres) DataType(t dialects.DataType) string {
	switch t {
	case dialects.DataTypeBlob:
		return "BYTEA"
	case dialects.DataTypeString:
		return "VARCHAR(255)"
	case dialects.DataTypeEnum:
		panic("not implemented")

	case dialects.DataTypeBoolean:
		return "BOOLEAN"

	case dialects.DataTypeDate, dialects.DataTypeDateTime:
		return "TIMESTAMP"

	case dialects.DataTypeFloat32:
		return "REAL"
	case dialects.DataTypeFloat64:
		return "DOUBLE PRECISION"

	case dialects.DataTypeInt8, dialects.DataTypeInt16, dialects.DataTypeUInt8, dialects.DataTypeUInt16:
		return "SMALLINT"
	case dialects.DataTypeInt32, dialects.DataTypeUInt32:
		return "INTEGER"
	case dialects.DataTypeInt64, dialects.DataTypeUInt64:
		return "BIGINT"

	case dialects.DataTypeJSON:
		return "JSON"
	}
	return string(t)
}

func (*Posgtgres) CurrentTime() string {
	return "CURRENT_TIMESTAMP()"
}

func (s *Posgtgres) Escape(v any) string {
	switch v := v.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return fmt.Sprint(v)
	case bool:
		if v {
			return "TRUE"
		} else {
			return "FALSE"
		}
	default:
		b, err := json.Marshal(v)
		if err != nil {
			panic(fmt.Errorf("failed to escape value: %w", err))
		}
		return s.Escape(string(b))
	}
}

func (p *Posgtgres) Binding() string {
	p.bindingNumber++
	return fmt.Sprintf("$%d", p.bindingNumber)
}

func UsePostgres() {
	dialects.SetDefaultDialect(func() dialects.Dialect {
		return &Posgtgres{}
	})
}
func init() {
	UsePostgres()
}
