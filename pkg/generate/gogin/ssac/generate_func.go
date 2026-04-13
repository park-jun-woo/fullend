//ff:func feature=ssac-gen type=generator control=sequence topic=output
//ff:what 단일 ServiceFunc의 Go 코드를 생성하는 래퍼
package ssac

import (
	"github.com/park-jun-woo/fullend/pkg/parser/funcspec"
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// GenerateFunc는 단일 ServiceFunc의 Go 코드를 생성한다.
func GenerateFunc(sf ssacparser.ServiceFunc, st *rule.Ground, funcSpecs []funcspec.FuncSpec) ([]byte, error) {
	return (&GoTarget{FuncSpecs: funcSpecs}).GenerateFunc(sf, st)
}
