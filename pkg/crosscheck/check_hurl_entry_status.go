//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkHurlEntryStatus — 단일 Hurl entry의 status code가 OpenAPI 응답에 있는지 검증 (X-37)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/parser/hurl"
)

func checkHurlEntryStatus(entry hurl.HurlEntry, fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	for _, item := range fs.OpenAPIDoc.Paths.Map() {
		errs = append(errs, checkHurlEntryOps(entry, item.Operations())...)
	}
	return errs
}
