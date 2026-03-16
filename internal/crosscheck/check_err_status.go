//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=ssac-openapi
//ff:what SSaC ErrStatus 코드가 OpenAPI 응답에 정의되어 있는지 검증
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// errStatusTypes are the SSaC sequence types that support custom ErrStatus.
var errStatusTypes = map[string]int{
	"empty":  404,
	"exists": 409,
	"state":  409,
	"auth":   403,
}

// checkErrStatus validates that SSaC ErrStatus codes are defined in OpenAPI responses.
func checkErrStatus(funcs []ssacparser.ServiceFunc, doc *openapi3.T) []CrossError {
	var errs []CrossError

	opMap := buildOperationMap(doc)

	for _, fn := range funcs {
		errs = append(errs, checkFuncErrStatus(fn, opMap)...)
	}

	return errs
}
