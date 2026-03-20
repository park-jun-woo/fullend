//ff:func feature=funcspec type=util control=sequence
//ff:what 함수 본문이 스텁(TODO + return)인지 확인한다
package funcspec

import (
	"go/ast"
	"go/token"
)

// isStubBody checks if function body only contains a return or panic("TODO").
func isStubBody(fset *token.FileSet, body *ast.BlockStmt) bool {
	if len(body.List) == 0 {
		return true
	}
	if len(body.List) > 1 {
		return false
	}
	// Single statement: check if it's a return.
	if _, isReturn := body.List[0].(*ast.ReturnStmt); isReturn {
		return true
	}
	// Single statement: check if it's a panic(...) call.
	exprStmt, ok := body.List[0].(*ast.ExprStmt)
	if !ok {
		return false
	}
	callExpr, ok := exprStmt.X.(*ast.CallExpr)
	if !ok {
		return false
	}
	ident, ok := callExpr.Fun.(*ast.Ident)
	return ok && ident.Name == "panic"
}
