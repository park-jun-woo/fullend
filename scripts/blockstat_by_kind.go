// blockstat_by_kind: if/for/range/switch 종류별 body 라인 수 분포
//
// Usage: go run scripts/blockstat_by_kind.go [dir]
package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type kindBucket struct {
	b1_5   int
	b6_10  int
	b11_20 int
	b21_50 int
	b51    int
	total  int
	max    int
}

func (b *kindBucket) add(lines int) {
	b.total++
	if lines > b.max {
		b.max = lines
	}
	switch {
	case lines <= 5:
		b.b1_5++
	case lines <= 10:
		b.b6_10++
	case lines <= 20:
		b.b11_20++
	case lines <= 50:
		b.b21_50++
	default:
		b.b51++
	}
}

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	buckets := map[string]*kindBucket{
		"if":     {},
		"for":    {},
		"range":  {},
		"switch": {},
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.Contains(path, "vendor/") || strings.Contains(path, "_test.go") {
			return nil
		}
		collectFromFile(path, buckets)
		return nil
	})

	fmt.Printf("%-8s %6s  %6s  %6s  %6s  %6s  %6s  %6s\n", "KIND", "TOTAL", "1-5", "6-10", "11-20", "21-50", "51+", "MAX")
	fmt.Println(strings.Repeat("-", 62))
	allTotal := 0
	for _, kind := range []string{"if", "for", "range", "switch"} {
		b := buckets[kind]
		fmt.Printf("%-8s %6d  %6d  %6d  %6d  %6d  %6d  %6d\n",
			kind, b.total, b.b1_5, b.b6_10, b.b11_20, b.b21_50, b.b51, b.max)
		allTotal += b.total
	}
	fmt.Println(strings.Repeat("-", 62))
	fmt.Printf("%-8s %6d\n", "TOTAL", allTotal)
}

func collectFromFile(path string, buckets map[string]*kindBucket) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return
	}
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		for _, stmt := range fn.Body.List {
			collectBlock(fset, stmt, buckets)
		}
	}
}

func collectBlock(fset *token.FileSet, stmt ast.Stmt, buckets map[string]*kindBucket) {
	switch s := stmt.(type) {
	case *ast.IfStmt:
		buckets["if"].add(bodyLines(fset, s.Body))
	case *ast.ForStmt:
		buckets["for"].add(bodyLines(fset, s.Body))
	case *ast.RangeStmt:
		buckets["range"].add(bodyLines(fset, s.Body))
	case *ast.SwitchStmt:
		buckets["switch"].add(bodyLines(fset, s.Body))
	case *ast.TypeSwitchStmt:
		buckets["switch"].add(bodyLines(fset, s.Body))
	}
}

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
