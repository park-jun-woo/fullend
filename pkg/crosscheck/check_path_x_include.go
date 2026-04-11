//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkPathXInclude — 개별 path operation의 x-include 검증
package crosscheck

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkPathXInclude(g *rule.Ground, path string, ops map[string]*openapi3.Operation) []CrossError {
	var errs []CrossError
	for _, op := range ops {
		if op.Extensions == nil {
			continue
		}
		raw, ok := op.Extensions["x-include"]
		if !ok {
			continue
		}
		var xInc struct{ Allowed []string `json:"allowed"` }
		data, _ := json.Marshal(raw)
		if json.Unmarshal(data, &xInc) != nil {
			continue
		}
		errs = append(errs, checkXIncludeEntries(g, path, xInc.Allowed)...)
	}
	return errs
}
