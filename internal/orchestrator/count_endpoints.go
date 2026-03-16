//ff:func feature=orchestrator type=util control=iteration dimension=1
//ff:what OpenAPI 문서에서 총 엔드포인트 수를 센다

package orchestrator

import "github.com/getkin/kin-openapi/openapi3"

// countEndpoints counts total operations across all paths in an OpenAPI document.
func countEndpoints(doc *openapi3.T) int {
	count := 0
	for _, pi := range doc.Paths.Map() {
		for range pi.Operations() {
			count++
		}
	}
	return count
}
