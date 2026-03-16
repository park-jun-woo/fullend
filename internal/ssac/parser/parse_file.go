//ff:func feature=ssac-parse type=parser control=iteration dimension=1
//ff:what 단일 .ssac 파일을 파싱하여 []ServiceFunc 반환
package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// ParseFile은 단일 .ssac 파일을 파싱하여 []ServiceFunc를 반환한다.
func ParseFile(path string) ([]ServiceFunc, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("Go 파싱 실패: %w", err)
	}

	imports := collectImports(f)
	structs := collectStructs(f)
	var funcs []ServiceFunc

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		sf, err := parseFuncDecl(fn, f, path, imports, structs)
		if err != nil {
			return nil, err
		}
		if sf != nil {
			funcs = append(funcs, *sf)
		}
	}
	return funcs, nil
}
