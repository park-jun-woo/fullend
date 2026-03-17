//ff:func feature=genapi type=util control=sequence
//ff:what 설정에 따라 백엔드 코드젠 구현체를 선택한다
package gen

import (
	"github.com/park-jun-woo/fullend/internal/gen/gogin"
	"github.com/park-jun-woo/fullend/internal/genapi"
	"github.com/park-jun-woo/fullend/internal/projectconfig"
)

func selectBackend(cfg *projectconfig.ProjectConfig) genapi.Backend {
	// Future: branch on cfg.Backend field.
	return &gogin.GoGin{}
}
