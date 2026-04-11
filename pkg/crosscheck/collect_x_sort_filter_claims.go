//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectXSortFilterClaims — OpenAPI x-sort, x-filter 확장에서 검증 대상 수집
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
)

func collectXSortFilterClaims(fs *fullend.Fullstack) []xSortFilterClaim {
	var claims []xSortFilterClaim
	for path, item := range fs.OpenAPIDoc.Paths.Map() {
		claims = append(claims, collectPathXClaims(item, path)...)
	}
	return claims
}
