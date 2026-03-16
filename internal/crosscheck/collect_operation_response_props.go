//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what 단일 Operation의 2xx 응답에서 스키마 속성명 수집
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// collectOperationResponseProps collects 2xx response schema property names from an operation.
func collectOperationResponseProps(op *openapi3.Operation) map[string]bool {
	props := make(map[string]bool)
	for code, respRef := range op.Responses.Map() {
		if len(code) != 3 || code[0] != '2' {
			continue
		}
		addResponseSchemaProps(props, respRef)
	}
	return props
}
