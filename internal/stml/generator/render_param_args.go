//ff:func feature=stml-gen type=util control=iteration dimension=1
//ff:what ParamBind 슬라이스에서 { key: value } 형태의 인자 문자열을 생성한다
package generator

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/fullend/internal/stml/parser"
)

func renderParamArgs(params []parser.ParamBind) string {
	if len(params) == 0 {
		return ""
	}
	var parts []string
	for _, p := range params {
		parts = append(parts, fmt.Sprintf("%s: %s", p.Name, paramSourceExpr(p)))
	}
	return "{ " + strings.Join(parts, ", ") + " }"
}
