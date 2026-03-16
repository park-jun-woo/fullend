//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what 단일 OpenAPI 오퍼레이션의 보안 요구사항이 미들웨어에 존재하는지 검증
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"
)

func checkOpSecurity(op *openapi3.Operation, mwSet map[string]bool, method, pathStr string) []CrossError {
	if op.Security == nil {
		return nil
	}
	var errs []CrossError
	for _, req := range *op.Security {
		errs = append(errs, checkSecurityReqNames(req, mwSet, method, pathStr)...)
	}
	return errs
}
