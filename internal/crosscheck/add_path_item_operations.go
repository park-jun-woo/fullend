//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-openapi
//ff:what PathItem의 Operation을 operationId 맵에 등록
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// addPathItemOperations adds all operations from a path item to the operation map.
func addPathItemOperations(opMap map[string]*openapi3.Operation, pathItem *openapi3.PathItem) {
	for _, op := range pathItemOperations(pathItem) {
		if op != nil && op.OperationID != "" {
			opMap[op.OperationID] = op
		}
	}
}
