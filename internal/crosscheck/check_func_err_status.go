//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=ssac-openapi
//ff:what 단일 SSaC 함수의 ErrStatus 코드가 OpenAPI에 정의되어 있는지 검증
package crosscheck

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// checkFuncErrStatus validates ErrStatus codes for a single SSaC function.
func checkFuncErrStatus(fn ssacparser.ServiceFunc, opMap map[string]*openapi3.Operation) []CrossError {
	op := opMap[fn.Name]
	if op == nil || op.Responses == nil {
		return nil
	}

	var errs []CrossError
	for seqIdx, seq := range fn.Sequences {
		defaultStatus, ok := errStatusTypes[seq.Type]
		if !ok {
			continue
		}

		statusCode := defaultStatus
		if seq.ErrStatus != 0 {
			statusCode = seq.ErrStatus
		}

		codeStr := fmt.Sprintf("%d", statusCode)
		resp := op.Responses.Status(statusCode)
		if resp == nil {
			errs = append(errs, CrossError{
				Rule:       "SSaC @" + seq.Type + " → OpenAPI",
				Context:    fmt.Sprintf("%s:%s seq[%d]", fn.FileName, fn.Name, seqIdx),
				Message:    fmt.Sprintf("SSaC @%s uses HTTP %s but OpenAPI %s has no %s response defined", seq.Type, codeStr, fn.Name, codeStr),
				Suggestion: fmt.Sprintf("OpenAPI %s responses에 %s 응답을 추가하세요", fn.Name, codeStr),
			})
		}
	}

	return errs
}
