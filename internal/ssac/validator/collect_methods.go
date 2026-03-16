//ff:func feature=symbol type=util control=iteration dimension=1
//ff:what InterfaceType에서 메서드 목록을 추출하여 ModelSymbol을 생성한다
package validator

import "go/ast"

func collectMethods(iface *ast.InterfaceType) ModelSymbol {
	ms := ModelSymbol{Methods: make(map[string]MethodInfo)}
	for _, method := range iface.Methods.List {
		if len(method.Names) == 0 {
			continue
		}
		ms.Methods[method.Names[0].Name] = MethodInfo{}
	}
	return ms
}
