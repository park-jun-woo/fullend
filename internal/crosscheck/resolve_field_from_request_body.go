//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what Operation의 요청 본문에서 필드의 Go 타입 추출
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// resolveFieldFromRequestBody extracts a field's Go type from an operation's request body.
func resolveFieldFromRequestBody(op *openapi3.Operation, fieldName string) string {
	if op.RequestBody == nil || op.RequestBody.Value == nil {
		return ""
	}
	for _, mt := range op.RequestBody.Value.Content {
		if mt.Schema == nil || mt.Schema.Value == nil {
			continue
		}
		propRef := findSchemaProperty(mt.Schema.Value, fieldName)
		if propRef != nil && propRef.Value != nil {
			return openAPITypeToGo(propRef.Value)
		}
	}
	return ""
}
