//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what openAPIHas2xx — operationId에 2xx 응답이 있는지 확인
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/fullend"

func openAPIHas2xx(fs *fullend.Fullstack, opID string) bool {
	for _, item := range fs.OpenAPIDoc.Paths.Map() {
		if opHas2xx(item.Operations(), opID) {
			return true
		}
	}
	return false
}
