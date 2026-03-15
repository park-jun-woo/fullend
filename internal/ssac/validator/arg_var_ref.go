//ff:func feature=ssac-validate type=util
//ff:what Arg가 변수 참조인 경우 루트 변수명을 반환한다
package validator

import "github.com/geul-org/fullend/internal/ssac/parser"

// argVarRef는 Arg가 변수 참조인 경우 루트 변수명을 반환한다.
func argVarRef(a parser.Arg) string {
	if a.Literal != "" {
		return ""
	}
	if a.Source == "request" || a.Source == "currentUser" || a.Source == "query" || a.Source == "" {
		return ""
	}
	return a.Source
}
