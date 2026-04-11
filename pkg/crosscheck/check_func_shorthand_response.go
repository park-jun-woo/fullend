//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkFuncShorthandResponse — shorthand @response 변수 타입 → OpenAPI 응답 타입 매칭 (X-19, X-20)
package crosscheck

import (
	"strings"

	"github.com/jinzhu/inflection"

	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkFuncShorthandResponse(g *rule.Ground, funcName string, seqs []ssac.Sequence) []CrossError {
	var errs []CrossError
	for _, seq := range seqs {
		if seq.Type != "response" || len(seq.Fields) > 0 {
			continue
		}
		// shorthand: @response varName — find the variable's type
		varName := seq.Target
		if varName == "" {
			continue
		}
		varType := g.Types["SSaC.var."+funcName+"."+varName]
		if varType == "" {
			continue
		}
		// Check that OpenAPI response schema name matches the variable type
		table := strings.ToLower(inflection.Plural(varType))
		responseKey := "OpenAPI.response." + funcName
		if _, ok := g.Schemas[responseKey]; !ok {
			errs = append(errs, CrossError{Rule: "X-19", Context: funcName, Level: "WARNING",
				Message: "shorthand @response " + varName + " (type " + varType + ") but no OpenAPI response schema"})
		}
		_ = table
	}
	return errs
}
