//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what OpenAPI에서 operationId별 Operation 맵 생성
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// buildOperationMap builds an operationId -> Operation map from the OpenAPI doc.
func buildOperationMap(doc *openapi3.T) map[string]*openapi3.Operation {
	opMap := make(map[string]*openapi3.Operation)
	if doc.Paths == nil {
		return opMap
	}
	for _, pathItem := range doc.Paths.Map() {
		addPathItemOperations(opMap, pathItem)
	}
	return opMap
}
