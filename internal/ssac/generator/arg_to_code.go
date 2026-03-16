//ff:func feature=ssac-gen type=generator control=selection
//ff:what 단일 Arg를 Go 코드 표현으로 변환
package generator

import (
	"github.com/ettle/strcase"
	"github.com/geul-org/fullend/internal/ssac/parser"
)

func argToCode(a parser.Arg) string {
	switch {
	case a.Literal != "":
		return `"` + a.Literal + `"`
	case a.Source == "query":
		return "opts"
	case a.Source == "request":
		return strcase.ToGoCamel(a.Field)
	case a.Source == "currentUser":
		return a.Source + "." + a.Field
	case a.Source != "":
		if a.Field == "" {
			return a.Source
		}
		return a.Source + "." + a.Field
	default:
		return a.Field
	}
}
