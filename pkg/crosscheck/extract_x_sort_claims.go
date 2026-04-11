//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what extractXSortClaims — x-sort 확장에서 컬럼 검증 대상 추출
package crosscheck

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
)

func extractXSortClaims(op *openapi3.Operation, path string) []xSortFilterClaim {
	raw, ok := op.Extensions["x-sort"]
	if !ok {
		return nil
	}
	var xSort struct {
		Allowed []string `json:"allowed"`
	}
	data, _ := json.Marshal(raw)
	if json.Unmarshal(data, &xSort) != nil {
		return nil
	}
	var claims []xSortFilterClaim
	for _, col := range xSort.Allowed {
		claims = append(claims, xSortFilterClaim{
			ruleID: "X-1", col: col, lookupKey: lookupKeyForPath(op),
			context: path + " x-sort", message: "x-sort column not found in DDL",
		})
	}
	return claims
}
