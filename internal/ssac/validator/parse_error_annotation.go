//ff:func feature=symbol type=util control=iteration dimension=1
//ff:what FuncDecl의 Doc에서 @error 어노테이션의 HTTP 상태 코드를 추출한다
package validator

import (
	"go/ast"
	"strconv"
	"strings"
)

// parseErrorAnnotation은 FuncDecl의 Doc에서 @error 어노테이션의 HTTP 상태 코드를 추출한다.
func parseErrorAnnotation(doc *ast.CommentGroup) int {
	if doc == nil {
		return 0
	}
	for _, comment := range doc.List {
		text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
		if !strings.HasPrefix(text, "@error ") {
			continue
		}
		if code, err := strconv.Atoi(strings.TrimSpace(text[7:])); err == nil {
			return code
		}
	}
	return 0
}
