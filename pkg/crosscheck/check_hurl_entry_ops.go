//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkHurlEntryOps — operation 목록에서 Hurl entry의 status code 응답 존재 여부 검증 (X-37)
package crosscheck

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/park-jun-woo/fullend/pkg/parser/hurl"
)

func checkHurlEntryOps(entry hurl.HurlEntry, ops map[string]*openapi3.Operation) []CrossError {
	var errs []CrossError
	for _, op := range ops {
		if op.OperationID == "" || op.Responses == nil {
			continue
		}
		if op.Responses.Value(entry.StatusCode) == nil {
			errs = append(errs, CrossError{Rule: "X-37", Context: fmt.Sprintf("%s:%d", entry.File, entry.Line),
				Level: "WARNING", Message: "Hurl status " + entry.StatusCode + " not defined in OpenAPI responses"})
		}
	}
	return errs
}
