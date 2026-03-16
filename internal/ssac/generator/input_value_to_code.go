//ff:func feature=ssac-gen type=generator control=selection
//ff:what inputs 값에 예약 소스 변환(query, request 등)을 적용
package generator

import (
	"strings"

	"github.com/ettle/strcase"
)

// inputValueToCode는 inputs 값에 argToCode와 동일한 예약 소스 변환을 적용한다.
func inputValueToCode(val string) string {
	switch {
	case val == "query":
		return "opts"
	case strings.HasPrefix(val, "request."):
		return strcase.ToGoCamel(val[len("request."):])
	default:
		return val
	}
}
