//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what OpenAPI 요청 스키마에서 필드의 Go 타입 조회
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// resolveOpenAPIFieldType looks up a field's Go type from the OpenAPI request schema.
func resolveOpenAPIFieldType(doc *openapi3.T, operationID, fieldName string) string {
	if doc == nil || doc.Paths == nil {
		return ""
	}
	for _, pathItem := range doc.Paths.Map() {
		if result := findOperationFieldType(pathItem, operationID, fieldName); result != "" {
			return result
		}
	}
	return ""
}
