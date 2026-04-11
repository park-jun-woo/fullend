//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkCursorOps — 하나의 path item에서 operation별 cursor 검증 위임
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

func checkCursorOps(item *openapi3.PathItem, path string) []CrossError {
	var errs []CrossError
	for _, op := range item.Operations() {
		if op.Extensions == nil {
			continue
		}
		errs = append(errs, checkOpCursor(op, path)...)
	}
	return errs
}
