//ff:func feature=symbol type=loader control=iteration dimension=2
//ff:what 서비스 파일의 import 경로에서 패키지 접두사 모델의 Go interface를 파싱한다
package validator

import (
	"os"
	"path/filepath"
	"strings"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// LoadPackageInterfaces는 서비스 파일의 import 경로에서 패키지 접두사 모델의 Go interface를 파싱한다.
// 패키지명 → import 경로 매핑 후, 해당 경로에서 interface를 찾아 st.Models["pkg.Model"]에 등록.
func (st *SymbolTable) LoadPackageInterfaces(funcs []ssacparser.ServiceFunc, projectRoot string) {
	// 1. 모든 서비스 파일에서 패키지 접두사 모델 수집
	pkgModels := map[string]bool{} // "session" → true
	for _, sf := range funcs {
		for _, seq := range sf.Sequences {
			if seq.Package != "" {
				pkgModels[seq.Package] = true
			}
		}
	}
	if len(pkgModels) == 0 {
		return
	}

	// 2. 서비스 파일 import에서 패키지명 → import 경로 매핑
	pkgPaths := map[string]string{} // "session" → "myapp/session"
	for _, sf := range funcs {
		for _, imp := range sf.Imports {
			// import 경로의 마지막 segment가 패키지명
			segments := strings.Split(imp, "/")
			pkgName := segments[len(segments)-1]
			if !pkgModels[pkgName] {
				continue
			}
			pkgPaths[pkgName] = imp
		}
	}

	// 3. 각 패키지 경로에서 Go interface 파싱
	for pkgName, impPath := range pkgPaths {
		// projectRoot 기준으로 경로 탐색 (상대 경로)
		dir := filepath.Join(projectRoot, impPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		st.loadPackageGoInterfaces(pkgName, dir)
	}
}
