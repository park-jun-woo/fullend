//ff:func feature=stat type=util control=selection
//ff:what AST 수신자/식별자에서 타입 이름 추출
package main

import "go/ast"

func exprName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return exprName(e.X)
	}
	return ""
}
