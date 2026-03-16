//ff:func feature=symbol type=util control=iteration dimension=1
//ff:what CommentGroup에서 지정한 annotation 문자열이 포함되어 있는지 검사한다
package validator

import (
	"go/ast"
	"strings"
)

// hasAnnotation은 CommentGroup에서 지정한 annotation 문자열이 포함되어 있는지 검사한다.
func hasAnnotation(cg *ast.CommentGroup, tag string) bool {
	if cg == nil {
		return false
	}
	for _, c := range cg.List {
		if strings.Contains(c.Text, tag) {
			return true
		}
	}
	return false
}
