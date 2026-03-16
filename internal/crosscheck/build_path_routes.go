//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=scenario-check
//ff:what 단일 OpenAPI 경로의 메서드별 정규화 라우트 생성
package crosscheck

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// buildPathRoutes builds normalized routes for a single OpenAPI path.
func buildPathRoutes(path string, pi *openapi3.PathItem) []apiRoute {
	segs := normalizeOpenAPIPath(path)
	var routes []apiRoute
	for method, op := range pi.Operations() {
		responseCodes := collectResponseCodes(op)
		routes = append(routes, apiRoute{
			Method:    strings.ToUpper(method),
			Segments:  segs,
			Responses: responseCodes,
		})
	}
	return routes
}
