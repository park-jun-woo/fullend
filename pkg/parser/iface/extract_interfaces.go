//ff:func feature=iface-parse type=parser control=iteration dimension=2
//ff:what extractInterfaces — ast.File 에서 인터페이스 선언 목록 추출
package iface

import "go/ast"

func extractInterfaces(f *ast.File) []Interface {
	var ifaces []Interface
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
			methods := collectMethodNames(iface)
			if len(methods) == 0 {
				continue
			}
			ifaces = append(ifaces, Interface{Name: ts.Name.Name, Methods: methods})
		}
	}
	return ifaces
}
