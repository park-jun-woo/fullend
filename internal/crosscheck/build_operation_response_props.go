//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what OpenAPI에서 operationId별 응답 스키마 속성명 수집
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// buildOperationResponseProps collects response schema property names per operationId from the OpenAPI doc.
func buildOperationResponseProps(doc *openapi3.T) map[string]map[string]bool {
	result := make(map[string]map[string]bool)
	if doc == nil || doc.Paths == nil {
		return result
	}

	for _, pathItem := range doc.Paths.Map() {
		addPathResponseProps(result, pathItem)
	}

	return result
}
