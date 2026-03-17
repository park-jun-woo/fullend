//ff:func feature=stml-gen type=util control=iteration dimension=1
//ff:what ParamBind 슬라이스에서 소스 표현식 값 목록을 생성한다
package generator

import (
	"strings"

	"github.com/park-jun-woo/fullend/internal/stml/parser"
)

func renderParamValues(params []parser.ParamBind) string {
	var parts []string
	for _, p := range params {
		parts = append(parts, paramSourceExpr(p))
	}
	return strings.Join(parts, ", ")
}
