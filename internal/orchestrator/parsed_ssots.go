//ff:type feature=orchestrator type=model
//ff:what 모든 SSOT 파싱 결과를 보관하는 타입 (Phase012: genapi 의존 제거 위해 인라인)
package orchestrator

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/internal/funcspec"
	"github.com/park-jun-woo/fullend/internal/policy"
	"github.com/park-jun-woo/fullend/internal/projectconfig"
	"github.com/park-jun-woo/fullend/internal/statemachine"
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
	ssacvalidator "github.com/park-jun-woo/fullend/internal/ssac/validator"
	stmlparser "github.com/park-jun-woo/fullend/internal/stml/parser"
)

// ParsedSSOTs holds all SSOT parsing results.
// ParseAll() populates this; crosscheck consume it.
type ParsedSSOTs struct {
	Config           *projectconfig.ProjectConfig
	OpenAPIDoc       *openapi3.T
	SymbolTable      *ssacvalidator.SymbolTable
	ServiceFuncs     []ssacparser.ServiceFunc
	STMLPages        []stmlparser.PageSpec
	StateDiagrams    []*statemachine.StateDiagram
	Policies         []*policy.Policy
	ProjectFuncSpecs []funcspec.FuncSpec
	FullendPkgSpecs  []funcspec.FuncSpec
	HurlFiles        []string
	ModelDir         string
	StatesErr        error
}
