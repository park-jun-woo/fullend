//ff:func feature=gen-hurl type=generator control=selection
//ff:what OpenAPI 스키마에서 더미 값을 생성한다
package hurl

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// generateDummyValue returns a dummy value for a schema field based on type, format, and field name hints.
func generateDummyValue(fieldName string, schema *openapi3.Schema, checkEnums map[string][]string) interface{} {
	if schema == nil {
		return "test_string"
	}

	// Use first enum value if available (OpenAPI enum).
	if len(schema.Enum) > 0 {
		return fmt.Sprint(schema.Enum[0])
	}

	// Use DDL CHECK constraint enum values.
	if vals, ok := checkEnums[fieldName]; ok && len(vals) > 0 {
		return vals[0]
	}

	// Field name hints (checked before type-based defaults).
	lower := strings.ToLower(fieldName)
	switch {
	case strings.Contains(lower, "password"):
		return "Password1234!"
	case strings.Contains(lower, "price") || strings.Contains(lower, "amount"):
		return 10000
	case strings.Contains(lower, "rating"):
		return 5
	case strings.Contains(lower, "url"):
		return "https://example.com/test"
	}

	switch schema.Type.Slice()[0] {
	case "string":
		switch schema.Format {
		case "email":
			return "test@example.com"
		case "date-time":
			return "2025-01-01T00:00:00Z"
		default:
			return "test_string"
		}
	case "integer":
		return 1
	case "number":
		return 1.0
	case "boolean":
		return true
	case "object":
		return map[string]interface{}{}
	case "array":
		return []interface{}{}
	default:
		return "test_string"
	}
}
