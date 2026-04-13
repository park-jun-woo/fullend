//ff:func feature=iface-parse type=util control=iteration dimension=1
//ff:what collectMethodNames — InterfaceType 의 메서드 이름을 선언 순서로 추출
package iface

import "go/ast"

func collectMethodNames(iface *ast.InterfaceType) []string {
	var names []string
	for _, method := range iface.Methods.List {
		if len(method.Names) == 0 {
			continue
		}
		names = append(names, method.Names[0].Name)
	}
	return names
}
