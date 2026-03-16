//ff:func feature=symbol type=util control=selection
//ff:what ast.Expr를 Go 타입 문자열로 변환한다
package validator

import "go/ast"

// exprToGoType는 ast.Expr를 Go 타입 문자열로 변환한다.
func exprToGoType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		if x, ok := t.X.(*ast.Ident); ok {
			return x.Name + "." + t.Sel.Name
		}
	case *ast.StarExpr:
		return "*" + exprToGoType(t.X)
	case *ast.ArrayType:
		return "[]" + exprToGoType(t.Elt)
	}
	return "interface{}"
}
