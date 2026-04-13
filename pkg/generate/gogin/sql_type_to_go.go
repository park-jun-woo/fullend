//ff:func feature=gen-gogin type=util control=selection
//ff:what maps a SQL type to a Go type

package gogin

import "strings"

// sqlTypeToGo maps a SQL type to a Go type.
func sqlTypeToGo(sqlType string) string {
	// Normalize: strip parenthesized args.
	upper := strings.ToUpper(sqlType)
	if idx := strings.Index(upper, "("); idx > 0 {
		upper = upper[:idx]
	}

	switch upper {
	case "BIGSERIAL", "BIGINT":
		return "int64"
	case "INT", "INTEGER":
		return "int64"
	case "VARCHAR", "TEXT":
		return "string"
	case "BOOLEAN", "BOOL":
		return "bool"
	case "TIMESTAMPTZ", "TIMESTAMP":
		return "time.Time"
	case "JSONB", "JSON":
		return "json.RawMessage"
	default:
		return "string"
	}
}
