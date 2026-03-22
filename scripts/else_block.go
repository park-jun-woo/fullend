//ff:func feature=stat type=util control=sequence
//ff:what if문의 else 절에서 BlockStmt 추출
package main

import "go/ast"

func elseBlock(stmt ast.Stmt) *ast.BlockStmt {
	if stmt == nil {
		return nil
	}
	if b, ok := stmt.(*ast.BlockStmt); ok {
		return b
	}
	if ifStmt, ok := stmt.(*ast.IfStmt); ok {
		return ifStmt.Body
	}
	return nil
}
