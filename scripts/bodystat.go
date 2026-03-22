// bodystat: 함수별 순수 body 라인 수 통계
// 단일 if 블록, for {} 블록의 라인을 제외한 순수 시퀀스 라인 수를 계산
//
// Usage: go run scripts/bodystat.go [dir]
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type funcStat struct {
	File      string
	Name      string
	TotalBody int // 전체 body 라인
	PureBody  int // if/for 블록 제외 순수 라인
}

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	var stats []funcStat
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.Contains(path, "vendor/") || strings.Contains(path, "_test.go") {
			return nil
		}
		stats = append(stats, parseFile(path)...)
		return nil
	})

	sort.Slice(stats, func(i, j int) bool { return stats[i].PureBody > stats[j].PureBody })

	fmt.Printf("%-60s %6s %6s\n", "FUNCTION", "TOTAL", "PURE")
	fmt.Println(strings.Repeat("-", 74))

	buckets := map[string]int{
		"1-20":   0,
		"21-50":  0,
		"51-100": 0,
		"101+":   0,
	}

	for _, s := range stats {
		if s.PureBody > 30 {
			fmt.Printf("%-60s %6d %6d\n", s.File+":"+s.Name, s.TotalBody, s.PureBody)
		}
		switch {
		case s.PureBody <= 20:
			buckets["1-20"]++
		case s.PureBody <= 50:
			buckets["21-50"]++
		case s.PureBody <= 100:
			buckets["51-100"]++
		default:
			buckets["101+"]++
		}
	}

	fmt.Println()
	fmt.Println("=== 분포 ===")
	fmt.Printf("  1-20줄:   %d개\n", buckets["1-20"])
	fmt.Printf("  21-50줄:  %d개\n", buckets["21-50"])
	fmt.Printf("  51-100줄: %d개\n", buckets["51-100"])
	fmt.Printf("  101+줄:   %d개\n", buckets["101+"])
	fmt.Printf("  합계:     %d개\n", len(stats))
}

func parseFile(path string) []funcStat {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil
	}

	var stats []funcStat
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		name := fn.Name.Name
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			name = exprName(fn.Recv.List[0].Type) + "." + name
		}

		totalStart := fset.Position(fn.Body.Lbrace).Line
		totalEnd := fset.Position(fn.Body.Rbrace).Line
		totalBody := totalEnd - totalStart - 1
		if totalBody < 0 {
			totalBody = 0
		}

		// depth-1 if/for 블록 라인 수 합산
		controlLines := 0
		for _, stmt := range fn.Body.List {
			controlLines += countControlLines(fset, stmt)
		}

		pureBody := totalBody - controlLines
		if pureBody < 0 {
			pureBody = 0
		}

		stats = append(stats, funcStat{
			File:      path,
			Name:      name,
			TotalBody: totalBody,
			PureBody:  pureBody,
		})
	}
	return stats
}

func countControlLines(fset *token.FileSet, stmt ast.Stmt) int {
	switch s := stmt.(type) {
	case *ast.IfStmt:
		return blockLines(fset, s.Body) + blockLines(fset, elseBlock(s.Else))
	case *ast.ForStmt:
		return blockLines(fset, s.Body)
	case *ast.RangeStmt:
		return blockLines(fset, s.Body)
	case *ast.SwitchStmt:
		return blockLines(fset, s.Body)
	case *ast.TypeSwitchStmt:
		return blockLines(fset, s.Body)
	case *ast.SelectStmt:
		return blockLines(fset, s.Body)
	}
	return 0
}

func blockLines(fset *token.FileSet, block *ast.BlockStmt) int {
	if block == nil {
		return 0
	}
	start := fset.Position(block.Lbrace).Line
	end := fset.Position(block.Rbrace).Line
	return end - start + 1 // { 부터 } 까지 전체
}

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

func exprName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return exprName(e.X)
	}
	return ""
}
