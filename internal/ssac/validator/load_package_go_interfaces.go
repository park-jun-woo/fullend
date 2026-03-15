//ff:func feature=symbol type=loader
//ff:what 디렉토리에서 Go interface를 파싱하여 "pkg.Model" 키로 등록한다
package validator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// loadPackageGoInterfaces는 디렉토리에서 Go interface를 파싱하여 "pkg.Model" 키로 등록한다.
// 또한 {Method}Request struct를 파싱하여 ParamTypes에 필드 타입을 저장한다.
func (st *SymbolTable) loadPackageGoInterfaces(pkgName, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	fset := token.NewFileSet()
	// 1차: Request struct 수집
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

	// 2차: interface 파싱
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
					if len(method.Names) > 0 {
						methodName := method.Names[0].Name
						var params []string
						if ft, ok := method.Type.(*ast.FuncType); ok && ft.Params != nil {
							for _, param := range ft.Params.List {
								if isContextType(param.Type) {
									continue
								}
								for _, name := range param.Names {
									params = append(params, name.Name)
								}
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

	// 3차: standalone function 파싱 (@call 대상)
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
			errStatus := 0
			if fd.Doc != nil {
				for _, comment := range fd.Doc.List {
					text := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
					if strings.HasPrefix(text, "@error ") {
						if code, err := strconv.Atoi(strings.TrimSpace(text[7:])); err == nil {
							errStatus = code
						}
					}
				}
			}

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
