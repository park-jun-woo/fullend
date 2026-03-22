// blockstat: if/for/switch 블록 내부 라인 수 통계
// depth-1 제어 구조 블록의 내부 body 라인 수를 측정
//
// Usage: go run scripts/blockstat.go [dir]
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

type blockStat struct {
	File  string
	Func  string
	Kind  string // "if", "for", "range", "switch"
	Lines int    // 블록 내부 라인 수
}

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	var stats []blockStat
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

	sort.Slice(stats, func(i, j int) bool { return stats[i].Lines > stats[j].Lines })

	fmt.Printf("%-55s %-8s %5s\n", "LOCATION", "KIND", "LINES")
	fmt.Println(strings.Repeat("-", 70))

	for _, s := range stats {
		if s.Lines >= 10 {
			fmt.Printf("%-55s %-8s %5d\n", s.File+":"+s.Func, s.Kind, s.Lines)
		}
	}

	buckets := map[string]int{
		"1-5":   0,
		"6-10":  0,
		"11-20": 0,
		"21-50": 0,
		"51+":   0,
	}

	for _, s := range stats {
		switch {
		case s.Lines <= 5:
			buckets["1-5"]++
		case s.Lines <= 10:
			buckets["6-10"]++
		case s.Lines <= 20:
			buckets["11-20"]++
		case s.Lines <= 50:
			buckets["21-50"]++
		default:
			buckets["51+"]++
		}
	}

	fmt.Println()
	fmt.Println("=== if/for/switch 블록 내부 라인 수 분포 ===")
	fmt.Printf("  1-5줄:   %d개\n", buckets["1-5"])
	fmt.Printf("  6-10줄:  %d개\n", buckets["6-10"])
	fmt.Printf("  11-20줄: %d개\n", buckets["11-20"])
	fmt.Printf("  21-50줄: %d개\n", buckets["21-50"])
	fmt.Printf("  51+줄:   %d개\n", buckets["51+"])
	fmt.Printf("  합계:    %d개\n", len(stats))
}

func parseFile(path string) []blockStat {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil
	}

	var stats []blockStat
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		funcName := fn.Name.Name
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			funcName = exprName(fn.Recv.List[0].Type) + "." + funcName
		}

		for _, stmt := range fn.Body.List {
			stats = append(stats, collectBlocks(fset, path, funcName, stmt)...)
		}
	}
	return stats
}

func collectBlocks(fset *token.FileSet, path, funcName string, stmt ast.Stmt) []blockStat {
	var stats []blockStat
	switch s := stmt.(type) {
	case *ast.IfStmt:
		stats = append(stats, blockStat{File: path, Func: funcName, Kind: "if", Lines: bodyLines(fset, s.Body)})
	case *ast.ForStmt:
		stats = append(stats, blockStat{File: path, Func: funcName, Kind: "for", Lines: bodyLines(fset, s.Body)})
	case *ast.RangeStmt:
		stats = append(stats, blockStat{File: path, Func: funcName, Kind: "range", Lines: bodyLines(fset, s.Body)})
	case *ast.SwitchStmt:
		stats = append(stats, blockStat{File: path, Func: funcName, Kind: "switch", Lines: bodyLines(fset, s.Body)})
	case *ast.TypeSwitchStmt:
		stats = append(stats, blockStat{File: path, Func: funcName, Kind: "switch", Lines: bodyLines(fset, s.Body)})
	case *ast.SelectStmt:
		stats = append(stats, blockStat{File: path, Func: funcName, Kind: "select", Lines: bodyLines(fset, s.Body)})
	}
	return stats
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

func exprName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return exprName(e.X)
	}
	return ""
}
