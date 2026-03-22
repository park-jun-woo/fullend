// patternstat: 함수 depth-1 제어문 패턴 분석
// 함수 body의 루트 레벨에서 if/for/switch 개수를 세어 패턴 분류
//
// Usage: go run scripts/patternstat.go [dir]
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

type pattern struct {
	Ifs      int
	Loops    int // for + range
	Switches int
}

func (p pattern) String() string {
	return fmt.Sprintf("if=%d loop=%d switch=%d", p.Ifs, p.Loops, p.Switches)
}

func (p pattern) Category() string {
	hasIf := p.Ifs > 0
	hasLoop := p.Loops > 0
	hasSwitch := p.Switches > 0

	parts := 0
	if hasIf {
		parts++
	}
	if hasLoop {
		parts++
	}
	if hasSwitch {
		parts++
	}

	if parts == 0 {
		return "pure-sequence"
	}
	if parts > 1 {
		return "mixed"
	}
	if hasIf {
		if p.Ifs == 1 {
			return "single-if"
		}
		return "multi-if"
	}
	if hasLoop {
		if p.Loops == 1 {
			return "single-loop"
		}
		return "multi-loop"
	}
	if p.Switches == 1 {
		return "single-switch"
	}
	return "multi-switch"
}

type funcInfo struct {
	File    string
	Name    string
	Pattern pattern
	Lines   int
}

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}

	var funcs []funcInfo
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.Contains(path, "vendor/") || strings.Contains(path, "_test.go") {
			return nil
		}
		funcs = append(funcs, parseFile(path)...)
		return nil
	})

	// 카테고리별 집계
	catCount := make(map[string]int)
	catLines := make(map[string][]int)
	for _, f := range funcs {
		cat := f.Pattern.Category()
		catCount[cat]++
		catLines[cat] = append(catLines[cat], f.Lines)
	}

	// 카테고리 정렬 출력
	cats := []string{"pure-sequence", "single-if", "multi-if", "single-loop", "multi-loop", "single-switch", "multi-switch", "mixed"}
	fmt.Println("=== 함수 패턴 분포 ===")
	fmt.Printf("%-18s %6s %8s %8s %8s\n", "PATTERN", "COUNT", "AVG", "MED", "MAX")
	fmt.Println(strings.Repeat("-", 52))
	for _, cat := range cats {
		n := catCount[cat]
		if n == 0 {
			continue
		}
		lines := catLines[cat]
		sort.Ints(lines)
		avg := sum(lines) / n
		med := lines[n/2]
		max := lines[n-1]
		fmt.Printf("%-18s %6d %8d %8d %8d\n", cat, n, avg, med, max)
	}
	fmt.Printf("%-18s %6d\n", "TOTAL", len(funcs))

	// mixed 상세
	fmt.Println()
	fmt.Println("=== mixed 패턴 상세 (if+loop+switch 혼합) ===")
	fmt.Printf("%-55s %s\n", "FUNCTION", "PATTERN")
	fmt.Println(strings.Repeat("-", 85))
	for _, f := range funcs {
		if f.Pattern.Category() == "mixed" {
			fmt.Printf("%-55s %s\n", f.File+":"+f.Name, f.Pattern)
		}
	}

	// multi-if 상세 (orchestrator 후보)
	fmt.Println()
	fmt.Println("=== multi-if 패턴 (orchestrator 후보) ===")
	fmt.Printf("%-55s %4s %s\n", "FUNCTION", "LINE", "PATTERN")
	fmt.Println(strings.Repeat("-", 85))
	var multiIfs []funcInfo
	for _, f := range funcs {
		if f.Pattern.Category() == "multi-if" {
			multiIfs = append(multiIfs, f)
		}
	}
	sort.Slice(multiIfs, func(i, j int) bool { return multiIfs[i].Pattern.Ifs > multiIfs[j].Pattern.Ifs })
	for _, f := range multiIfs {
		fmt.Printf("%-55s %4d %s\n", f.File+":"+f.Name, f.Lines, f.Pattern)
	}
}

func parseFile(path string) []funcInfo {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil
	}

	var funcs []funcInfo
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		name := fn.Name.Name
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			name = exprName(fn.Recv.List[0].Type) + "." + name
		}

		start := fset.Position(fn.Body.Lbrace).Line
		end := fset.Position(fn.Body.Rbrace).Line
		lines := end - start - 1
		if lines < 0 {
			lines = 0
		}

		p := pattern{}
		for _, stmt := range fn.Body.List {
			switch stmt.(type) {
			case *ast.IfStmt:
				p.Ifs++
			case *ast.ForStmt, *ast.RangeStmt:
				p.Loops++
			case *ast.SwitchStmt, *ast.TypeSwitchStmt:
				p.Switches++
			}
		}

		funcs = append(funcs, funcInfo{
			File:    path,
			Name:    name,
			Pattern: p,
			Lines:   lines,
		})
	}
	return funcs
}

func sum(a []int) int {
	s := 0
	for _, v := range a {
		s += v
	}
	return s
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
