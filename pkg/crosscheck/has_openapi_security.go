//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what hasOpenAPISecurity — OpenAPI 경로에 security 설정이 하나라도 있는지 확인
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/fullend"

func hasOpenAPISecurity(fs *fullend.Fullstack) bool {
	for _, item := range fs.OpenAPIDoc.Paths.Map() {
		if pathHasSecurity(item) {
			return true
		}
	}
	return false
}
