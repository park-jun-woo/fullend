//ff:func feature=stml-gen type=util control=sequence
//ff:what ParamBind에서 route. 접두사가 있으면 파라미터 이름을 반환한다
package generator

import (
	"strings"

	"github.com/geul-org/fullend/internal/stml/parser"
)

func extractRouteParamName(p parser.ParamBind) string {
	if strings.HasPrefix(p.Source, "route.") {
		return strings.TrimPrefix(p.Source, "route.")
	}
	return ""
}
