//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-openapi
//ff:what PathItem에서 operationId로 필드 타입 조회
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// findOperationFieldType finds a field type from a path item's operations.
func findOperationFieldType(pathItem *openapi3.PathItem, operationID, fieldName string) string {
	for _, op := range pathItemOperations(pathItem) {
		if op == nil || op.OperationID != operationID {
			continue
		}
		return resolveFieldFromRequestBody(op, fieldName)
	}
	return ""
}
