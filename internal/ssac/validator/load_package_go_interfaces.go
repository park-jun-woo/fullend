//ff:func feature=symbol type=loader control=sequence topic=go-interface
//ff:what 디렉토리에서 Go interface를 파싱하여 "pkg.Model" 키로 등록한다
package validator

import (
	"go/token"
	"os"
)

// loadPackageGoInterfaces는 디렉토리에서 Go interface를 파싱하여 "pkg.Model" 키로 등록한다.
// 또한 {Method}Request struct를 파싱하여 ParamTypes에 필드 타입을 저장한다.
func (st *SymbolTable) loadPackageGoInterfaces(pkgName, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	fset := token.NewFileSet()
	requestStructs := collectRequestStructs(fset, entries, dir)
	st.parsePackageInterfaces(fset, entries, dir, pkgName, requestStructs)
	st.parseStandaloneFuncs(fset, entries, dir, pkgName, requestStructs)
}
