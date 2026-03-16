//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what Hurl 시나리오 파일이 유효한 OpenAPI 경로를 참조하는지 검증
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// CheckHurlFiles validates that .hurl scenario files reference valid OpenAPI paths.
func CheckHurlFiles(hurlFiles []string, doc *openapi3.T) []CrossError {
	var errs []CrossError

	routes := buildHurlRoutes(doc)

	for _, f := range hurlFiles {
		entries := parseHurlFile(f)
		for _, e := range entries {
			errs = append(errs, validateHurlEntry(e, routes)...)
		}
	}

	return errs
}
