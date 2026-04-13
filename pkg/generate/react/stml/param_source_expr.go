//ff:func feature=stml-gen type=util control=sequence
//ff:what ParamBind의 Source를 JSX 표현식으로 변환한다
package stml

import (
	"strings"

	stmlparser "github.com/park-jun-woo/fullend/pkg/parser/stml"
)

func paramSourceExpr(p stmlparser.ParamBind) string {
	if strings.HasPrefix(p.Source, "route.") {
		return strings.TrimPrefix(p.Source, "route.")
	}
	return p.Source
}
