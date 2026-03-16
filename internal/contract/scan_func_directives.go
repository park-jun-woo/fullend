//ff:func feature=contract type=walker control=iteration dimension=1
//ff:what Go 소스에서 함수별 fullend 디렉티브를 추출한다
package contract

import (
	"go/ast"
	"go/parser"
	"go/token"
)

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
