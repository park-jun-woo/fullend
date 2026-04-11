//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what pathHasSecurity — path item의 operation 중 security 설정이 있는지 확인
package crosscheck

import "github.com/getkin/kin-openapi/openapi3"

func pathHasSecurity(item *openapi3.PathItem) bool {
	for _, op := range item.Operations() {
		if op.Security != nil {
			return true
		}
	}
	return false
}
