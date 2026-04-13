//ff:func feature=ssac-gen type=generator control=selection topic=output
//ff:what ServiceFunc의 subscribe 여부에 따라 HTTP/Subscribe 코드 생성을 분기
package ssac

import (
	ssacparser "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

// GenerateFunc는 단일 ServiceFunc의 Go 코드를 생성한다.
func (g *GoTarget) GenerateFunc(sf ssacparser.ServiceFunc, st *rule.Ground) ([]byte, error) {
	if sf.Subscribe != nil {
		return g.generateSubscribeFunc(sf, st)
	}
	return g.generateHTTPFunc(sf, st)
}
