//ff:func feature=funcspec type=parser control=iteration dimension=1
//ff:what 코멘트 그룹에서 @func, @error, @description 어노테이션을 추출한다
package funcspec

import (
	"go/ast"
	"strings"
)

func parseCommentGroup(cg *ast.CommentGroup, spec *FuncSpec) {
	for _, c := range cg.List {
		line := strings.TrimPrefix(c.Text, "//")
		line = strings.TrimSpace(line)
		applyAnnotation(line, spec)
	}
}
