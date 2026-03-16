//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=scenario-check
//ff:what OpenAPI에서 Hurl 검증용 정규화 라우트 목록 생성
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

// buildHurlRoutes builds normalized OpenAPI routes for hurl validation.
func buildHurlRoutes(doc *openapi3.T) []apiRoute {
	var routes []apiRoute
	if doc.Paths == nil {
		return routes
	}
	for path, pi := range doc.Paths.Map() {
		routes = append(routes, buildPathRoutes(path, pi)...)
	}
	return routes
}
