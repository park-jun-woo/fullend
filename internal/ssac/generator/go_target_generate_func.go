//ff:func feature=ssac-gen type=generator control=selection
//ff:what ServiceFunc의 subscribe 여부에 따라 HTTP/Subscribe 코드 생성을 분기
package generator

import (
	"github.com/geul-org/fullend/internal/ssac/parser"
	"github.com/geul-org/fullend/internal/ssac/validator"
)

// GenerateFunc는 단일 ServiceFunc의 Go 코드를 생성한다.
func (g *GoTarget) GenerateFunc(sf parser.ServiceFunc, st *validator.SymbolTable) ([]byte, error) {
	if sf.Subscribe != nil {
		return g.generateSubscribeFunc(sf, st)
	}
	return g.generateHTTPFunc(sf, st)
}
