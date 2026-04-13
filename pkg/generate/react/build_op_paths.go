//ff:func feature=gen-react type=util control=iteration dimension=1
//ff:what OpenAPI 문서에서 operationID -> path 매핑을 구축한다

package react

import "github.com/getkin/kin-openapi/openapi3"

// buildOpPaths builds a mapping from operationID to OpenAPI path.
func buildOpPaths(doc *openapi3.T) map[string]string {
	opPaths := make(map[string]string)
	if doc == nil || doc.Paths == nil {
		return opPaths
	}
	for path, pi := range doc.Paths.Map() {
		addOpPaths(opPaths, path, pi)
	}
	return opPaths
}
