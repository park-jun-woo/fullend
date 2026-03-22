//ff:func feature=stat type=util control=sequence
//ff:what BlockStmt의 { 부터 } 까지 전체 라인 수 계산
package main

import (
	"go/ast"
	"go/token"
)

func blockLines(fset *token.FileSet, block *ast.BlockStmt) int {
	if block == nil {
		return 0
	}
	start := fset.Position(block.Lbrace).Line
	end := fset.Position(block.Rbrace).Line
	return end - start + 1
}
