//ff:func feature=ssac-gen type=util control=sequence topic=type-resolve
//ff:what @call의 에러 상태 코드를 해석하여 Go 상수 문자열로 반환
package ssac

import "github.com/park-jun-woo/fullend/pkg/rule"

func resolveCallErrStatus(errStatus int, st *rule.Ground, model string) string {
	if errStatus != 0 {
		return httpStatusConst(errStatus)
	}
	if code := lookupCallErrStatus(st, model); code != 0 {
		return httpStatusConst(code)
	}
	return "http.StatusInternalServerError"
}
