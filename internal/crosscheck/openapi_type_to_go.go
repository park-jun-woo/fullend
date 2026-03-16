//ff:func feature=crosscheck type=util control=selection
//ff:what OpenAPI 스키마 타입을 Go 타입으로 변환
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// openAPITypeToGo converts an OpenAPI schema type to a Go type.
func openAPITypeToGo(schema *openapi3.Schema) string {
	switch schema.Type.Slice()[0] {
	case "string":
		if schema.Format == "date-time" {
			return "time.Time"
		}
		return "string"
	case "integer":
		if schema.Format == "int32" {
			return "int32"
		}
		return "int64"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "array":
		if schema.Items != nil && schema.Items.Value != nil {
			return "[]" + openAPITypeToGo(schema.Items.Value)
		}
		return "[]interface{}"
	default:
		return "interface{}"
	}
}
