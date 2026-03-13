package contract

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SpliceResult holds the merged content and any warnings.
type SpliceResult struct {
	Content  string
	Warnings []Warning
}

// Warning records a contract mismatch between preserved body and regenerated contract.
type Warning struct {
	File        string
	Function    string
	OldContract string
	NewContract string
}

// PreservedFunc holds a preserved function's directive and body text.
type PreservedFunc struct {
	Directive Directive
	BodyText  string // raw source between { and }, excluding the braces themselves
}

// PreserveSnapshot captures all preserved functions/files before code generation.
type PreserveSnapshot struct {
	FilePreserves map[string]string                    // path → saved whole file content
	FuncPreserves map[string]map[string]*PreservedFunc // path → funcName → preserved body
}

// ScanPreserveSnapshot walks a directory and captures all preserved content.
func ScanPreserveSnapshot(dir string) *PreserveSnapshot {
	snap := &PreserveSnapshot{
		FilePreserves: make(map[string]string),
		FuncPreserves: make(map[string]map[string]*PreservedFunc),
	}

	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		src, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		content := string(src)

		// Check file-level preserve.
		if hasFilePreserve(content) {
			snap.FilePreserves[path] = content
			return nil
		}

		// Check function-level preserves.
		funcs := scanPreservedFromSource(content)
		if len(funcs) > 0 {
			snap.FuncPreserves[path] = funcs
		}
		return nil
	})

	return snap
}

// RestorePreserved restores all preserved content after code generation.
func RestorePreserved(snap *PreserveSnapshot) []Warning {
	var allWarnings []Warning

	// Restore file-level preserves.
	for path, content := range snap.FilePreserves {
		if _, err := os.Stat(path); err == nil {
			os.WriteFile(path, []byte(content), 0644)
		}
	}

	// Restore function-level preserves.
	for path, funcs := range snap.FuncPreserves {
		src, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		result, err := SpliceWithPreserved(string(src), funcs, path)
		if err != nil {
			continue
		}
		os.WriteFile(path, []byte(result.Content), 0644)
		allWarnings = append(allWarnings, result.Warnings...)

		// Write .new file for contract changes.
		for _, w := range result.Warnings {
			newPath := path + ".new"
			os.WriteFile(newPath, src, 0644)
			_ = w // warning already recorded
			break  // one .new per file is enough
		}
	}

	return allWarnings
}

// CountPreserveFuncs counts all preserved functions and files in a directory.
func CountPreserveFuncs(dir string) int {
	count := 0
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		src, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		content := string(src)
		if hasFilePreserve(content) {
			count++
			return nil
		}
		funcs := scanPreservedFromSource(content)
		count += len(funcs)
		return nil
	})
	return count
}

// SpliceWithPreserved merges preserved function bodies into new content.
func SpliceWithPreserved(newContent string, preserved map[string]*PreservedFunc, filePath string) (*SpliceResult, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", newContent, parser.ParseComments)
	if err != nil {
		return &SpliceResult{Content: newContent}, nil
	}

	type replacement struct {
		start int
		end   int
		text  string
	}

	var replacements []replacement
	var warnings []Warning

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}

		pf, ok := preserved[fd.Name.Name]
		if !ok {
			continue
		}

		// Check contract change.
		newD := extractDirectiveFromDoc(fd.Doc)
		if newD != nil && newD.Contract != pf.Directive.Contract {
			warnings = append(warnings, Warning{
				File:        filePath,
				Function:    fd.Name.Name,
				OldContract: pf.Directive.Contract,
				NewContract: newD.Contract,
			})
		}

		// Replace body.
		bodyStart := fset.Position(fd.Body.Lbrace).Offset
		bodyEnd := fset.Position(fd.Body.Rbrace).Offset
		replacements = append(replacements, replacement{
			start: bodyStart + 1,
			end:   bodyEnd,
			text:  pf.BodyText,
		})

		// Update directive from gen to preserve.
		if newD != nil {
			for _, c := range fd.Doc.List {
				if _, err := Parse(c.Text); err == nil {
					lineStart := fset.Position(c.Pos()).Offset
					lineEnd := lineStart + len(c.Text)
					preserveD := &Directive{
						Ownership: "preserve",
						SSOT:      newD.SSOT,
						Contract:  pf.Directive.Contract,
					}
					replacements = append(replacements, replacement{
						start: lineStart,
						end:   lineEnd,
						text:  preserveD.String(),
					})
					break
				}
			}
		}
	}

	if len(replacements) == 0 {
		return &SpliceResult{Content: newContent}, nil
	}

	// Apply replacements in reverse order to preserve byte offsets.
	sort.Slice(replacements, func(i, j int) bool {
		return replacements[i].start > replacements[j].start
	})

	result := newContent
	for _, r := range replacements {
		result = result[:r.start] + r.text + result[r.end:]
	}

	return &SpliceResult{Content: result, Warnings: warnings}, nil
}

// scanPreservedFromSource extracts all preserved functions from Go source.
func scanPreservedFromSource(src string) map[string]*PreservedFunc {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil
	}

	result := make(map[string]*PreservedFunc)

	for _, decl := range f.Decls {
		fd, ok := decl.(*ast.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}

		d := extractDirectiveFromDoc(fd.Doc)
		if d == nil || d.Ownership != "preserve" {
			continue
		}

		bodyStart := fset.Position(fd.Body.Lbrace).Offset
		bodyEnd := fset.Position(fd.Body.Rbrace).Offset
		bodyText := src[bodyStart+1 : bodyEnd]

		result[fd.Name.Name] = &PreservedFunc{
			Directive: *d,
			BodyText:  bodyText,
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}

// extractDirectiveFromDoc finds a //fullend: directive in a doc comment group.
func extractDirectiveFromDoc(doc *ast.CommentGroup) *Directive {
	if doc == nil {
		return nil
	}
	for _, c := range doc.List {
		if d, err := Parse(c.Text); err == nil {
			return d
		}
	}
	return nil
}

// hasFilePreserve checks if source has a file-level //fullend:preserve directive.
func hasFilePreserve(src string) bool {
	lines := strings.SplitN(src, "\n", 10)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "//") {
			if d, err := Parse(line); err == nil && d.Ownership == "preserve" {
				return true
			}
			continue
		}
		if strings.HasPrefix(line, "package ") {
			break
		}
	}
	return false
}
