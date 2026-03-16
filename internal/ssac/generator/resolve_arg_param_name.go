//ff:func feature=ssac-gen type=util control=selection topic=args-inputs
//ff:what Arg에서 Go 파라미터 이름을 추론
package generator

import (
	"github.com/ettle/strcase"
	"github.com/geul-org/fullend/internal/ssac/parser"
)

func resolveArgParamName(a parser.Arg) string {
	switch {
	case a.Literal != "":
		return strcase.ToGoCamel(a.Literal)
	case a.Source == "request" || a.Source == "currentUser":
		return strcase.ToGoCamel(a.Field)
	case a.Source != "":
		return a.Source + strcase.ToGoPascal(a.Field)
	default:
		return strcase.ToGoCamel(a.Field)
	}
}
