//ff:func feature=contract type=walker control=iteration dimension=1
//ff:what Go 소스에서 보존된 함수들을 추출한다
package contract

import (
	"go/ast"
	"go/parser"
	"go/token"
)

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
		if !ok {
			continue
		}
		name, pf, found := extractPreservedFunc(fd, src, fset)
		if found {
			result[name] = pf
		}
	}

	if len(result) == 0 {
		return nil
	}
	return result
}
