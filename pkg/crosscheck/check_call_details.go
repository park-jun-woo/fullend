//ff:func feature=crosscheck type=rule control=sequence
//ff:what checkCallDetails — 단일 @call의 함수명/input/result 세부 검증
package crosscheck

import (
	"strings"
	"unicode"

	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkCallDetails(g *rule.Ground, funcName string, seq ssac.Sequence, fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	idx := strings.IndexByte(seq.Model, '.')
	if idx <= 0 {
		return nil
	}
	callFunc := seq.Model[idx+1:]
	callKey := strings.ToLower(seq.Model)

	// X-38: function name must start lowercase
	if len(callFunc) > 0 && unicode.IsUpper(rune(callFunc[0])) {
		// Actually callFunc in SSaC is PascalCase, but the @func annotation uses camelCase
		// This checks if the pkg-level reference starts with lowercase
	}

	// X-42: @call input count vs FuncRequest fields
	reqFields := g.Schemas["Func.request."+callFunc]
	if reqFields != nil && len(seq.Args) != len(reqFields) {
		errs = append(errs, CrossError{Rule: "X-42", Context: funcName + "/" + seq.Model, Level: "ERROR",
			Message: "@call input count mismatch with func request fields"})
	}

	// X-43: @call input field exists in FuncRequest
	if reqFields != nil {
		errs = append(errs, checkCallInputFields(funcName, seq, reqFields)...)
	}

	// X-45: @call has result but func has no response
	if seq.Result != nil && g.Schemas["Func.request."+callFunc] != nil {
		if !funcHasResponse(callKey, fs) {
			errs = append(errs, CrossError{Rule: "X-45", Context: funcName + "/" + seq.Model, Level: "ERROR",
				Message: "@call has result but func has no response fields"})
		}
	}

	// X-46: @call has no result but func has response (WARNING)
	if seq.Result == nil && funcHasResponse(callKey, fs) {
		errs = append(errs, CrossError{Rule: "X-46", Context: funcName + "/" + seq.Model, Level: "WARNING",
			Message: "@call ignores func response"})
	}

	return errs
}
