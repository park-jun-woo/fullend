//ff:func feature=crosscheck type=rule control=iteration dimension=2 topic=openapi-ddl
//ff:what SSaC request 참조 필드가 OpenAPI required에 포함되는지 검증
package crosscheck

import (
	"fmt"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
)

func checkC1RequiredMissing(funcs []ssacparser.ServiceFunc, st *ssacvalidator.SymbolTable) []CrossError {
	var errs []CrossError
	for _, fn := range funcs {
		rs, ok := st.RequestSchemas[fn.Name]
		if !ok {
			continue
		}
		usedFields := collectRequestFields(fn)
		for _, field := range sortedStrings(usedFields) {
			fc, exists := rs.Fields[field]
			if !exists {
				continue
			}
			if !fc.Required {
				errs = append(errs, CrossError{
					Rule:       "OpenAPI Constraints C1",
					Context:    fn.Name,
					Message:    fmt.Sprintf("field %q is used in SSaC sequences (request.%s) but not marked required in OpenAPI requestBody", field, field),
					Suggestion: fmt.Sprintf("OpenAPI requestBody의 required 배열에 %q 추가", field),
				})
			}
		}
	}
	return errs
}
