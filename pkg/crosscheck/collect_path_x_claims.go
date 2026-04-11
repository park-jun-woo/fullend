//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectPathXClaims — 하나의 path item에서 x-sort/x-filter 검증 대상 수집
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

func collectPathXClaims(item *openapi3.PathItem, path string) []xSortFilterClaim {
	var claims []xSortFilterClaim
	for _, op := range item.Operations() {
		if op.Extensions == nil {
			continue
		}
		claims = append(claims, extractXSortClaims(op, path)...)
		claims = append(claims, extractXFilterClaims(op, path)...)
	}
	return claims
}
