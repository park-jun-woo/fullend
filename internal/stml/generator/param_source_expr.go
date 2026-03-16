//ff:func feature=stml-gen type=util control=sequence
//ff:what ParamBind의 Source를 JSX 표현식으로 변환한다
package generator

import (
	"strings"

	"github.com/geul-org/fullend/internal/stml/parser"
)

func paramSourceExpr(p parser.ParamBind) string {
	if strings.HasPrefix(p.Source, "route.") {
		return strings.TrimPrefix(p.Source, "route.")
	}
	return p.Source
}
