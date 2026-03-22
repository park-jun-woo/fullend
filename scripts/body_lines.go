//ff:func feature=stat type=util control=sequence
//ff:what BlockStmt 내부 라인 수 계산 (Lbrace~Rbrace 사이)
package main

import (
	"go/ast"
	"go/token"
)

func bodyLines(fset *token.FileSet, block *ast.BlockStmt) int {
	if block == nil {
		return 0
	}
	start := fset.Position(block.Lbrace).Line
	end := fset.Position(block.Rbrace).Line
	lines := end - start - 1
	if lines < 0 {
		return 0
	}
	return lines
}
