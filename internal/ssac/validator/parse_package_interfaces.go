//ff:func feature=symbol type=loader control=iteration dimension=5
//ff:what Go 파일에서 interface를 파싱하여 SymbolTable.Models에 등록한다
package validator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// parsePackageInterfaces는 디렉토리의 Go 파일에서 interface를 파싱하여 SymbolTable에 등록한다.
func (st *SymbolTable) parsePackageInterfaces(fset *token.FileSet, entries []os.DirEntry, dir, pkgName string, requestStructs map[string]map[string]string) {
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
				iface, ok := ts.Type.(*ast.InterfaceType)
				if !ok {
					continue
				}
				ms := ModelSymbol{Methods: make(map[string]MethodInfo)}
				for _, method := range iface.Methods.List {
					if len(method.Names) == 0 {
						continue
					}
					methodName := method.Names[0].Name
					var params []string
					ft, ok := method.Type.(*ast.FuncType)
					if !ok || ft.Params == nil {
						ms.Methods[methodName] = MethodInfo{Params: params}
						continue
					}
					for _, param := range ft.Params.List {
						if isContextType(param.Type) {
							continue
						}
						for _, name := range param.Names {
							params = append(params, name.Name)
						}
					}
					mi := MethodInfo{Params: params}
					// Request struct 매칭: {MethodName}Request
					reqStructName := methodName + "Request"
					if fields, ok := requestStructs[reqStructName]; ok {
						mi.ParamTypes = fields
					}
					ms.Methods[methodName] = mi
				}
				if len(ms.Methods) > 0 {
					// "Model" suffix 제거: "SessionModel" → "Session"
					modelName := ts.Name.Name
					if strings.HasSuffix(modelName, "Model") {
						modelName = modelName[:len(modelName)-5]
					}
					key := pkgName + "." + modelName
					st.Models[key] = ms
				}
			}
		}
	}
}
