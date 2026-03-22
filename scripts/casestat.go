// casestat: switch/select의 case 절 내부 라인 수 통계
//
// Usage: go run scripts/casestat.go [dir]
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

type caseStat struct {
	File  string
	Func  string
	Lines int
}

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	var stats []caseStat
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

	// 상위 출력
	fmt.Printf("%-60s %5s\n", "LOCATION", "LINES")
	fmt.Println(strings.Repeat("-", 67))
	for _, s := range stats {
		if s.Lines >= 10 {
			fmt.Printf("%-60s %5d\n", s.File+":"+s.Func, s.Lines)
		}
	}

	// 분포
	buckets := [6]int{} // 0:1줄, 1:2-3, 2:4-5, 3:6-10, 4:11-20, 5:21+
	for _, s := range stats {
		switch {
		case s.Lines <= 1:
			buckets[0]++
		case s.Lines <= 3:
			buckets[1]++
		case s.Lines <= 5:
			buckets[2]++
		case s.Lines <= 10:
			buckets[3]++
		case s.Lines <= 20:
			buckets[4]++
		default:
			buckets[5]++
		}
	}

	fmt.Println()
	fmt.Println("=== case 절 내부 라인 수 분포 ===")
	fmt.Printf("  1줄:     %d개\n", buckets[0])
	fmt.Printf("  2-3줄:   %d개\n", buckets[1])
	fmt.Printf("  4-5줄:   %d개\n", buckets[2])
	fmt.Printf("  6-10줄:  %d개\n", buckets[3])
	fmt.Printf("  11-20줄: %d개\n", buckets[4])
	fmt.Printf("  21+줄:   %d개\n", buckets[5])
	fmt.Printf("  합계:    %d개\n", len(stats))
}

func parseFile(path string) []caseStat {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil
	}

	var stats []caseStat
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		funcName := fn.Name.Name
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			funcName = exprName(fn.Recv.List[0].Type) + "." + funcName
		}

		ast.Inspect(fn.Body, func(n ast.Node) bool {
			cc, ok := n.(*ast.CaseClause)
			if ok {
				lines := caseLines(fset, cc)
				stats = append(stats, caseStat{File: path, Func: funcName, Lines: lines})
			}
			return true
		})
	}
	return stats
}

func caseLines(fset *token.FileSet, cc *ast.CaseClause) int {
	if len(cc.Body) == 0 {
		return 0
	}
	start := fset.Position(cc.Colon).Line
	last := cc.Body[len(cc.Body)-1]
	end := fset.Position(last.End()).Line
	lines := end - start
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
