package contract

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// FuncStatus describes the contract status of a single function.
type FuncStatus struct {
	File      string // relative path from artifacts dir
	Function  string
	Directive Directive
	Status    string // "gen", "preserve", "broken", "orphan"
	Detail    string // violation detail
}

// ScanDir scans artifacts directory for all Go files with //fullend: directives.
func ScanDir(artifactsDir string) ([]FuncStatus, error) {
	var results []FuncStatus

	backendDir := filepath.Join(artifactsDir, "backend")
	err := filepath.WalkDir(backendDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(artifactsDir, path)
		content := string(src)

		// Check file-level directive.
		if d := findFileLevelDirective(content); d != nil {
			results = append(results, FuncStatus{
				File:      relPath,
				Function:  "(file)",
				Directive: *d,
				Status:    d.Ownership,
			})
			return nil
		}

		// Check function-level directives.
		funcStatuses := scanFuncDirectives(content, relPath)
		results = append(results, funcStatuses...)
		return nil
	})

	return results, err
}

// findFileLevelDirective finds a file-level //fullend: directive in source.
func findFileLevelDirective(src string) *Directive {
	lines := strings.SplitN(src, "\n", 10)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "//") {
			if d, err := Parse(line); err == nil {
				return d
			}
			continue
		}
		if strings.HasPrefix(line, "package ") {
			break
		}
	}
	return nil
}

// scanFuncDirectives parses Go source and extracts function-level directives.
func scanFuncDirectives(src, relPath string) []FuncStatus {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil
	}

	var results []FuncStatus
	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		d := extractDirectiveFromDoc(fd.Doc)
		if d == nil {
			continue
		}

		results = append(results, FuncStatus{
			File:      relPath,
			Function:  fd.Name.Name,
			Directive: *d,
			Status:    d.Ownership,
		})
	}

	return results
}
