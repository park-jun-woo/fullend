// structural_metrics: Phase011 구조 건전성 지표 측정.
// internal/gen 과 pkg/generate 를 비교해 리포트 생성.
//
// 지표:
//  - 함수별 매개변수 수 분포 (평균, 중앙값, 최대, 8+ 함수 수)
//  - 파일 수, 평균/중앙값 파일 줄 수
//  - *WithDomains 중복 함수 쌍
//  - Decide* 순수 판정 함수 (Phase010 대체 지표. Toulmin 미채택.)
//  - toulmin.NewGraph / g.Rule / g.Counter / g.Except 호출 수 (fullend 전체)
//
// Usage: go run scripts/structural_metrics.go
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

type fileMetric struct {
	Path       string
	Lines      int
	FuncCount  int
	ParamCount []int // per func
	DeciderFns int   // func name starts with "Decide"
	WithDomFns int   // ends with "WithDomains"
	TulGraphs  int   // toulmin.NewGraph calls (for counting)
}

type areaMetric struct {
	Name     string
	Files    int
	Lines    []int
	Params   []int
	DecFns   int
	WithDFns int
	TulN     int
}

func (a *areaMetric) add(m fileMetric) {
	a.Files++
	a.Lines = append(a.Lines, m.Lines)
	a.Params = append(a.Params, m.ParamCount...)
	a.DecFns += m.DeciderFns
	a.WithDFns += m.WithDomFns
	a.TulN += m.TulGraphs
}

func main() {
	internal := &areaMetric{Name: "internal/gen"}
	pkg := &areaMetric{Name: "pkg/generate"}

	walk("internal/gen", internal)
	walk("pkg/generate", pkg)

	fmt.Println("# Structural Metrics Report (Phase011)")
	fmt.Println()
	fmt.Println("2026-04-13 · pkg/generate vs internal/gen")
	fmt.Println()
	fmt.Println("## 파일·줄 수")
	fmt.Printf("                    %-12s  %-12s\n", internal.Name, pkg.Name)
	fmt.Printf("파일 수             %-12d  %-12d\n", internal.Files, pkg.Files)
	fmt.Printf("평균 줄 수          %-12.1f  %-12.1f\n", mean(internal.Lines), mean(pkg.Lines))
	fmt.Printf("중앙값 줄 수        %-12.0f  %-12.0f\n", median(internal.Lines), median(pkg.Lines))
	fmt.Printf("최대 줄 수          %-12d  %-12d\n", maxInt(internal.Lines), maxInt(pkg.Lines))
	fmt.Println()
	fmt.Println("## 함수 매개변수 분포")
	fmt.Printf("                    %-12s  %-12s\n", internal.Name, pkg.Name)
	fmt.Printf("총 함수 수          %-12d  %-12d\n", len(internal.Params), len(pkg.Params))
	fmt.Printf("평균 매개변수       %-12.2f  %-12.2f\n", mean(internal.Params), mean(pkg.Params))
	fmt.Printf("중앙값              %-12.0f  %-12.0f\n", median(internal.Params), median(pkg.Params))
	fmt.Printf("최대                %-12d  %-12d\n", maxInt(internal.Params), maxInt(pkg.Params))
	fmt.Printf("8+ params 함수      %-12d  %-12d\n", countGE(internal.Params, 8), countGE(pkg.Params, 8))
	fmt.Printf("5+ params 함수      %-12d  %-12d\n", countGE(internal.Params, 5), countGE(pkg.Params, 5))
	fmt.Println()
	fmt.Println("## 중복 패턴 *WithDomains")
	fmt.Printf("                    %-12s  %-12s\n", internal.Name, pkg.Name)
	fmt.Printf("*WithDomains 함수   %-12d  %-12d\n", internal.WithDFns, pkg.WithDFns)
	fmt.Println()
	fmt.Println("## Decide* 순수 판정 함수 (Phase010 구조 정비 지표)")
	fmt.Printf("                    %-12s  %-12s\n", internal.Name, pkg.Name)
	fmt.Printf("Decide* 함수 수     %-12d  %-12d\n", internal.DecFns, pkg.DecFns)
	fmt.Println()
	fmt.Println("## Toulmin 사용 (참고)")
	fmt.Printf("                    %-12s  %-12s\n", internal.Name, pkg.Name)
	fmt.Printf("toulmin.NewGraph    %-12d  %-12d\n", internal.TulN, pkg.TulN)
	fmt.Println()
	fmt.Println("> Phase010 결정: 2-depth 이내 if-else 로 해결되어 Toulmin 미채택. 대신 Decide* 순수 함수 3곳 수렴.")
	fmt.Println("> fullend 전체 Toulmin 사용:", countToulminAll())
}

func walk(root string, a *areaMetric) {
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}
		m := analyze(path)
		a.add(m)
		return nil
	})
}

func analyze(path string) fileMetric {
	m := fileMetric{Path: path}
	src, err := os.ReadFile(path)
	if err != nil {
		return m
	}
	m.Lines = strings.Count(string(src), "\n") + 1

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, src, parser.SkipObjectResolution)
	if err != nil {
		return m
	}
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		m.FuncCount++
		params := 0
		if fn.Type.Params != nil {
			for _, field := range fn.Type.Params.List {
				n := len(field.Names)
				if n == 0 {
					n = 1
				}
				params += n
			}
		}
		m.ParamCount = append(m.ParamCount, params)
		if strings.HasPrefix(fn.Name.Name, "Decide") {
			m.DeciderFns++
		}
		if strings.HasSuffix(fn.Name.Name, "WithDomains") {
			m.WithDomFns++
		}
	}
	m.TulGraphs = strings.Count(string(src), "toulmin.NewGraph(")
	return m
}

func countToulminAll() int {
	roots := []string{"pkg", "internal", "cmd"}
	total := 0
	for _, r := range roots {
		total += walkAndCount(r, "toulmin.NewGraph(")
	}
	return total
}

func walkAndCount(root, needle string) int {
	n := 0
	_ = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		b, _ := os.ReadFile(path)
		n += strings.Count(string(b), needle)
		return nil
	})
	return n
}

func mean(xs []int) float64 {
	if len(xs) == 0 {
		return 0
	}
	total := 0
	for _, x := range xs {
		total += x
	}
	return float64(total) / float64(len(xs))
}

func median(xs []int) float64 {
	if len(xs) == 0 {
		return 0
	}
	sorted := append([]int(nil), xs...)
	sort.Ints(sorted)
	mid := len(sorted) / 2
	if len(sorted)%2 == 0 {
		return float64(sorted[mid-1]+sorted[mid]) / 2
	}
	return float64(sorted[mid])
}

func maxInt(xs []int) int {
	if len(xs) == 0 {
		return 0
	}
	m := xs[0]
	for _, x := range xs[1:] {
		if x > m {
			m = x
		}
	}
	return m
}

func countGE(xs []int, threshold int) int {
	n := 0
	for _, x := range xs {
		if x >= threshold {
			n++
		}
	}
	return n
}
