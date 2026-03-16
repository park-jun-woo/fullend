//ff:func feature=crosscheck type=util control=sequence topic=openapi-ddl
//ff:what 스키마에서 속성명 또는 camelCase 대체 조회
package crosscheck

import (
	"github.com/ettle/strcase"
	"github.com/getkin/kin-openapi/openapi3"
)

// findSchemaProperty looks up a property by name or camelCase fallback.
func findSchemaProperty(schema *openapi3.Schema, fieldName string) *openapi3.SchemaRef {
	if propRef, ok := schema.Properties[fieldName]; ok {
		return propRef
	}
	if propRef, ok := schema.Properties[strcase.ToGoCamel(fieldName)]; ok {
		return propRef
	}
	return nil
}
