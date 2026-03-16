//ff:func feature=symbol type=loader control=iteration dimension=3 topic=go-interface
//ff:what model/*.go에서 interface(component)와 func을 추출한다
package validator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// loadGoInterfaces는 model/*.go에서 interface(component)와 func을 추출한다.
func (st *SymbolTable) loadGoInterfaces(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	fset := token.NewFileSet()
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") {
			continue
		}

		f, err := parser.ParseFile(fset, filepath.Join(dir, entry.Name()), nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("%s 파싱 실패: %w", entry.Name(), err)
		}

		for _, decl := range f.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			// GenDecl의 Doc에서 @dto 감지
			hasDtoTag := hasAnnotation(gd.Doc, "@dto")

			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				// @dto 태그가 있으면 DTO로 등록 (GenDecl Doc 또는 TypeSpec Doc)
				if hasDtoTag || hasAnnotation(ts.Doc, "@dto") {
					st.DTOs[ts.Name.Name] = true
					hasDtoTag = false
				}

				// interface → Models에 등록
				iface, ok := ts.Type.(*ast.InterfaceType)
				if !ok {
					continue
				}
				ms := collectMethods(iface)
				if len(ms.Methods) == 0 {
					continue
				}
				st.Models[ts.Name.Name] = ms
			}
		}

		// ast.Decls에서 FuncDecl은 GenDecl과 별개로 순회
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || fd.Recv != nil {
				continue
			}
			st.Funcs[fd.Name.Name] = true
		}
	}

	return nil
}
