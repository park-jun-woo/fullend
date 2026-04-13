//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what openAPIHasResponse — operationId에 해당 status code 응답이 있는지 확인
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/fullend"

func openAPIHasResponse(fs *fullend.Fullstack, opID, code string) bool {
	for _, item := range fs.OpenAPIDoc.Paths.Map() {
		if opHasResponse(item.Operations(), opID, code) {
			return true
		}
	}
	return false
}
