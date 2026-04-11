//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what funcHasResponse — func spec에 response 필드가 있는지 확인
package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
)

func funcHasResponse(callKey string, fs *fullend.Fullstack) bool {
	allSpecs := append(fs.ProjectFuncSpecs, fs.FullendPkgSpecs...)
	for _, sp := range allSpecs {
		key := strings.ToLower(sp.Package + "." + sp.Name)
		if key == callKey && len(sp.ResponseFields) > 0 {
			return true
		}
	}
	return false
}
