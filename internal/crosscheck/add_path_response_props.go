//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=ssac-openapi
//ff:what 단일 PathItem의 모든 Operation 응답 속성을 결과 맵에 추가
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// addPathResponseProps adds response property names from a path item to the result map.
func addPathResponseProps(result map[string]map[string]bool, pathItem *openapi3.PathItem) {
	for _, op := range pathItemOperations(pathItem) {
		if op == nil || op.OperationID == "" || op.Responses == nil {
			continue
		}
		props := collectOperationResponseProps(op)
		if len(props) > 0 {
			result[op.OperationID] = props
		}
	}
}
