//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkCursor — cursor pagination x-sort 제약 검증 (X-7, X-8)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
)

func checkCursor(fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil {
		return nil
	}
	var errs []CrossError
	for path, item := range fs.OpenAPIDoc.Paths.Map() {
		errs = append(errs, checkCursorOps(item, path)...)
	}
	return errs
}
