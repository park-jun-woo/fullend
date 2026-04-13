//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=func-check
//ff:what matchFuncSpec — 단일 funcspec 리스트에서 pkg.name 일치 항목 조회

package crosscheck

import (
	"strings"

	"github.com/park-jun-woo/fullend/pkg/parser/funcspec"
)

func matchFuncSpec(list []funcspec.FuncSpec, pkg, name string) *funcspec.FuncSpec {
	for i := range list {
		fs := &list[i]
		if fs.Package == pkg && strings.EqualFold(fs.Name, name) {
			return fs
		}
	}
	return nil
}
