//ff:func feature=funcspec type=util control=sequence
//ff:what 함수 본문이 스텁(TODO + return)인지 확인한다
package funcspec

import (
	"go/ast"
	"go/token"
)

// isStubBody checks if function body only contains "// TODO: implement" and a return.
func isStubBody(fset *token.FileSet, body *ast.BlockStmt) bool {
	if len(body.List) == 0 {
		return true
	}
	if len(body.List) > 1 {
		return false
	}
	// Single statement: check if it's a return.
	_, isReturn := body.List[0].(*ast.ReturnStmt)
	return isReturn
}
