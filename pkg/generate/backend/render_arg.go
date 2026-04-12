//ff:func feature=rule type=util control=selection
//ff:what renderArg — 단일 Arg를 Go 식으로 렌더
package backend

import parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func renderArg(arg parsessac.Arg) string {
	switch {
	case arg.Literal != "":
		return arg.Literal
	case arg.Source != "" && arg.Field != "":
		return arg.Source + "." + arg.Field
	case arg.Source != "":
		return arg.Source
	}
	return ""
}
