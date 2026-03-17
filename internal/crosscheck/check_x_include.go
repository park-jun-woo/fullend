//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=openapi-ddl
//ff:what x-include FK 조인 필드가 DDL 테이블과 일치하는지 검증
package crosscheck

import (
	"github.com/getkin/kin-openapi/openapi3"

	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func checkXInclude(op *openapi3.Operation, st *ssacvalidator.SymbolTable, ctx string, primaryTable string) []CrossError {
	var errs []CrossError

	raw, ok := op.Extensions["x-include"]
	if !ok {
		return errs
	}

	var includeExt struct {
		Allowed []string `json:"allowed"`
	}
	if err := unmarshalExt(raw, &includeExt); err != nil {
		return errs
	}

	for _, spec := range includeExt.Allowed {
		errs = append(errs, validateXIncludeSpec(spec, ctx, primaryTable, st)...)
	}

	return errs
}
