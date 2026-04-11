//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkHurlStatus — Hurl status code → OpenAPI 응답 정의 존재 WARNING (X-37)
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/fullend"

func checkHurlStatus(fs *fullend.Fullstack) []CrossError {
	if len(fs.HurlEntries) == 0 || fs.OpenAPIDoc == nil {
		return nil
	}
	var errs []CrossError
	for _, entry := range fs.HurlEntries {
		if entry.StatusCode == "" {
			continue
		}
		errs = append(errs, checkHurlEntryStatus(entry, fs)...)
	}
	return errs
}
