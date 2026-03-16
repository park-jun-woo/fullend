//ff:func feature=symbol type=loader control=iteration dimension=4
//ff:what Go 파일에서 *Request struct를 수집하여 필드 맵을 반환한다
package validator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// collectRequestStructs는 디렉토리의 Go 파일에서 *Request struct를 수집하여 필드 맵을 반환한다.
func collectRequestStructs(fset *token.FileSet, entries []os.DirEntry, dir string) map[string]map[string]string {
	requestStructs := map[string]map[string]string{} // "VerifyPasswordRequest" → {"Email": "string", "Password": "string"}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, entry.Name()), nil, parser.ParseComments)
		if err != nil {
			continue
		}
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				if !strings.HasSuffix(ts.Name.Name, "Request") {
					continue
				}
				st2, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				fields := map[string]string{}
				for _, field := range st2.Fields.List {
					typeName := exprToGoType(field.Type)
					for _, name := range field.Names {
						fields[name.Name] = typeName
					}
				}
				if len(fields) > 0 {
					requestStructs[ts.Name.Name] = fields
				}
			}
		}
	}
	return requestStructs
}
