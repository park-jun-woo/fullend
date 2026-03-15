//ff:func feature=symbol type=loader
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

			// GenDecl 또는 TypeSpec의 Doc에서 @dto 감지
			hasDtoTag := false
			if gd.Doc != nil {
				for _, c := range gd.Doc.List {
					if strings.Contains(c.Text, "@dto") {
						hasDtoTag = true
						break
					}
				}
			}

			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				// TypeSpec 자체의 Doc도 확인
				if !hasDtoTag && ts.Doc != nil {
					for _, c := range ts.Doc.List {
						if strings.Contains(c.Text, "@dto") {
							hasDtoTag = true
							break
						}
					}
				}

				// @dto 태그가 있으면 DTO로 등록
				if hasDtoTag {
					st.DTOs[ts.Name.Name] = true
					hasDtoTag = false // 다음 spec을 위해 리셋
				}

				// interface → Models에 등록
				if _, ok := ts.Type.(*ast.InterfaceType); ok {
					// interface의 메서드도 Models에 등록
					ms := ModelSymbol{Methods: make(map[string]MethodInfo)}
					iface := ts.Type.(*ast.InterfaceType)
					for _, method := range iface.Methods.List {
						if len(method.Names) > 0 {
							ms.Methods[method.Names[0].Name] = MethodInfo{}
						}
					}
					if len(ms.Methods) > 0 {
						st.Models[ts.Name.Name] = ms
					}
				}
			}

			// 패키지 레벨 func → Funcs로 등록
			fd, ok := decl.(*ast.FuncDecl)
			if ok && fd.Recv == nil {
				st.Funcs[fd.Name.Name] = true
			}
		}

		// ast.Decls에서 FuncDecl은 GenDecl과 별개로 순회
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if ok && fd.Recv == nil {
				st.Funcs[fd.Name.Name] = true
			}
		}
	}

	return nil
}
