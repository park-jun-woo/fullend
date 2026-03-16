//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=config-check
//ff:what OpenAPI 엔드포인트의 보안 참조가 미들웨어에 존재하는지 검증
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

func checkEndpointSecurity(doc *openapi3.T, mwSet map[string]bool) []CrossError {
	var errs []CrossError
	for pathStr, pathItem := range doc.Paths.Map() {
		for method, op := range pathItem.Operations() {
			errs = append(errs, checkOpSecurity(op, mwSet, method, pathStr)...)
		}
	}
	return errs
}
