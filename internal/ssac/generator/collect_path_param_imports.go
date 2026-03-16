//ff:func feature=ssac-gen type=util control=iteration dimension=1
//ff:what 경로 파라미터에 비문자열 타입이 있으면 strconv import 추가
package generator

import "github.com/geul-org/fullend/internal/ssac/validator"

func collectPathParamImports(pathParams []validator.PathParam, seen map[string]bool) {
	for _, pp := range pathParams {
		if pp.GoType != "string" {
			seen["strconv"] = true
			return
		}
	}
}
