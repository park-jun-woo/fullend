//ff:type feature=genapi type=model
//ff:what 모든 SSOT 파싱 결과를 보관하는 타입
package genapi

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/geul-org/fullend/internal/funcspec"
	"github.com/geul-org/fullend/internal/policy"
	"github.com/geul-org/fullend/internal/projectconfig"
	"github.com/geul-org/fullend/internal/statemachine"
	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
	ssacvalidator "github.com/geul-org/fullend/internal/ssac/validator"
	stmlparser "github.com/geul-org/fullend/internal/stml/parser"
)

// ParsedSSOTs holds all SSOT parsing results.
// orchestrator.ParseAll() populates this; crosscheck and gen consume it.
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
