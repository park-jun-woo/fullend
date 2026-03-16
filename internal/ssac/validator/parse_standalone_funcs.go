//ff:func feature=symbol type=loader control=iteration dimension=2 topic=go-interface
//ff:what Go 파일에서 standalone 함수를 파싱하여 SymbolTable.Models에 등록한다
package validator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// parseStandaloneFuncs는 디렉토리의 Go 파일에서 standalone 함수(@call 대상)를 파싱하여 SymbolTable에 등록한다.
func (st *SymbolTable) parseStandaloneFuncs(fset *token.FileSet, entries []os.DirEntry, dir, pkgName string, requestStructs map[string]map[string]string) {
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}
		f, err := parser.ParseFile(fset, filepath.Join(dir, entry.Name()), nil, parser.ParseComments)
		if err != nil {
			continue
		}
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Recv != nil {
				continue
			}
			funcName := fd.Name.Name
			reqStructName := funcName + "Request"
			if _, ok := requestStructs[reqStructName]; !ok {
				continue
			}

			// @error 어노테이션 파싱
			errStatus := parseErrorAnnotation(fd.Doc)

			modelKey := pkgName + "._func"
			ms, exists := st.Models[modelKey]
			if !exists {
				ms = ModelSymbol{Methods: make(map[string]MethodInfo)}
			}
			mi := MethodInfo{
				ParamTypes: requestStructs[reqStructName],
				ErrStatus:  errStatus,
			}
			ms.Methods[funcName] = mi
			st.Models[modelKey] = ms
		}
	}
}
