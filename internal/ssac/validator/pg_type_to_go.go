//ff:func feature=symbol type=util control=selection topic=ddl
//ff:what PostgreSQL 타입을 Go 타입으로 매핑한다
package validator

import "strings"

// pgTypeToGo는 PostgreSQL 타입을 Go 타입으로 매핑한다.
func pgTypeToGo(pgType string) string {
	switch pgType {
	case "BIGINT", "BIGSERIAL", "INTEGER", "SERIAL", "INT", "SMALLINT":
		return "int64"
	case "VARCHAR", "TEXT", "UUID", "CHAR":
		return "string"
	case "BOOLEAN", "BOOL":
		return "bool"
	case "TIMESTAMPTZ", "TIMESTAMP", "DATE":
		return "time.Time"
	case "NUMERIC", "DECIMAL", "REAL", "FLOAT", "DOUBLE":
		return "float64"
	case "JSONB", "JSON":
		return "json.RawMessage"
	default:
		// VARCHAR(255) 같은 경우
		if strings.HasPrefix(pgType, "VARCHAR") || strings.HasPrefix(pgType, "CHAR") {
			return "string"
		}
		return "string"
	}
}
