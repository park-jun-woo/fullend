//ff:func feature=orchestrator type=rule control=iteration dimension=1
//ff:what 프로젝트 func이 built-in 패키지를 override하는지 검증
package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/funcspec"
)

func checkBuiltinOverride(projectSpecs, fullendSpecs []funcspec.FuncSpec) []string {
	if len(fullendSpecs) == 0 {
		return nil
	}
	builtinPkgs := make(map[string]bool)
	for _, s := range fullendSpecs {
		builtinPkgs[s.Package] = true
	}
	var errs []string
	for _, s := range projectSpecs {
		if builtinPkgs[s.Package] {
			errs = append(errs, fmt.Sprintf(
				"func/%s: built-in 패키지 %q를 override할 수 없습니다. 커스텀 패키지명을 사용하세요 (예: func/my%s/)",
				s.Package, s.Package, s.Package))
		}
	}
	return errs
}
