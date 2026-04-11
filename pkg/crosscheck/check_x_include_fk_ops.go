//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkXIncludeFKOps — 단일 경로의 operation별 x-include FK 검증
package crosscheck

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func checkXIncludeFKOps(path string, item *openapi3.PathItem, fkCols rule.StringSet) []CrossError {
	var errs []CrossError
	for _, op := range item.Operations() {
		if op.Extensions == nil {
			continue
		}
		raw, ok := op.Extensions["x-include"]
		if !ok {
			continue
		}
		var xInc struct {
			Allowed []string `json:"allowed"`
		}
		data, _ := json.Marshal(raw)
		if json.Unmarshal(data, &xInc) != nil {
			continue
		}
		errs = append(errs, checkXIncludeAllowed(path, xInc.Allowed, fkCols)...)
	}
	return errs
}
