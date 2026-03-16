//ff:func feature=symbol type=util control=sequence
//ff:what ast.Expr이 context.Context 타입인지 확인한다
package validator

import "go/ast"

// isContextType는 ast.Expr이 context.Context 타입인지 확인한다.
func isContextType(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "context" && sel.Sel.Name == "Context"
}
