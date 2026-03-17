//ff:func feature=funcspec type=parser control=iteration dimension=1
//ff:what 디렉토리 내 모든 .go 파일에서 구조체 타입과 필드를 수집한다
package funcspec

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// collectPackageTypes parses all .go files in dir (non-recursive)
// and returns a map of struct name to fields.
func collectPackageTypes(dir string) map[string][]Field {
	result := make(map[string][]Field)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return result
	}
	fset := token.NewFileSet()
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(dir, e.Name())
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			continue
		}
		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}
				result[ts.Name.Name] = extractFields(st)
			}
		}
	}
	return result
}
