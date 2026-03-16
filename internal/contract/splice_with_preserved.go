//ff:func feature=contract type=util control=iteration dimension=1
//ff:what 새 콘텐츠에 보존 함수 본문을 병합한다
package contract

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
)

// SpliceWithPreserved merges preserved function bodies into new content.
func SpliceWithPreserved(newContent string, preserved map[string]*PreservedFunc, filePath string) (*SpliceResult, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", newContent, parser.ParseComments)
	if err != nil {
		return &SpliceResult{Content: newContent}, nil
	}

	var allReplacements []spliceReplacement
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

		replacements, warn := spliceFunc(fd, fset, pf, filePath)
		allReplacements = append(allReplacements, replacements...)
		if warn != nil {
			warnings = append(warnings, *warn)
		}
	}

	if len(allReplacements) == 0 {
		return &SpliceResult{Content: newContent}, nil
	}

	// Apply replacements in reverse order to preserve byte offsets.
	sort.Slice(allReplacements, func(i, j int) bool {
		return allReplacements[i].start > allReplacements[j].start
	})

	result := newContent
	for _, r := range allReplacements {
		result = result[:r.start] + r.text + result[r.end:]
	}

	return &SpliceResult{Content: result, Warnings: warnings}, nil
}
