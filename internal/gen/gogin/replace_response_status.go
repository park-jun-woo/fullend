//ff:func feature=gen-gogin type=generator control=sequence
//ff:what __RESPONSE_STATUS__ 플레이스홀더를 OpenAPI 성공 상태 코드로 치환한다

package gogin

import (
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// replaceResponseStatus replaces __RESPONSE_STATUS__ with OpenAPI success code.
func replaceResponseStatus(src string, doc *openapi3.T, operationID string) string {
	if !strings.Contains(src, "__RESPONSE_STATUS__") || doc == nil || operationID == "" {
		return src
	}
	statusConst := resolveSuccessStatus(doc, operationID)
	if statusConst == "" {
		return src
	}
	if statusConst == "http.StatusNoContent" {
		re := regexp.MustCompile(`c\.JSON\(__RESPONSE_STATUS__,\s*[^)]+\)`)
		return re.ReplaceAllString(src, "c.Status(http.StatusNoContent)")
	}
	return strings.ReplaceAll(src, "__RESPONSE_STATUS__", statusConst)
}
