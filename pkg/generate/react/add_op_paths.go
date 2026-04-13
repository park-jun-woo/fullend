//ff:func feature=gen-react type=util control=iteration dimension=1
//ff:what 단일 PathItem의 operationID를 경로 매핑에 추가한다

package react

import "github.com/getkin/kin-openapi/openapi3"

// addOpPaths adds operationID->path mappings from a single path item.
func addOpPaths(opPaths map[string]string, path string, pi *openapi3.PathItem) {
	for _, op := range pi.Operations() {
		if op == nil || op.OperationID == "" {
			continue
		}
		opPaths[op.OperationID] = path
	}
}
