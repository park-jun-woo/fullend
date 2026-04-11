//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what extractXFilterClaims — x-filter 확장에서 컬럼 검증 대상 추출
package crosscheck

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
)

func extractXFilterClaims(op *openapi3.Operation, path string) []xSortFilterClaim {
	raw, ok := op.Extensions["x-filter"]
	if !ok {
		return nil
	}
	var xFilter struct {
		Allowed []string `json:"allowed"`
	}
	data, _ := json.Marshal(raw)
	if json.Unmarshal(data, &xFilter) != nil {
		return nil
	}
	var claims []xSortFilterClaim
	for _, col := range xFilter.Allowed {
		claims = append(claims, xSortFilterClaim{
			ruleID: "X-3", col: col, lookupKey: lookupKeyForPath(op),
			context: path + " x-filter", message: "x-filter column not found in DDL",
		})
	}
	return claims
}
