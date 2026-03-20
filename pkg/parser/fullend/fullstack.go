//ff:type feature=orchestrator type=model
//ff:what 모든 SSOT 파싱 결과를 담는 풀스택 컨테이너
package fullend

import (
	pg_query "github.com/pganalyze/pg_query_go/v5"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/open-policy-agent/opa/ast"

	"github.com/park-jun-woo/fullend/pkg/diagnostic"
	"github.com/park-jun-woo/fullend/pkg/parser/funcspec"
	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
	"github.com/park-jun-woo/fullend/pkg/parser/hurl"
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/parser/statemachine"
	"github.com/park-jun-woo/fullend/pkg/parser/stml"
)

// Fullstack holds all SSOT parsing results.
// ParseAll() populates this; crosscheck and gen consume it.
type Fullstack struct {
	Config           *manifest.ProjectConfig
	OpenAPIDoc       *openapi3.T
	DDLResults       []*pg_query.ParseResult
	Policies         []*ast.Module
	ServiceFuncs     []ssac.ServiceFunc
	STMLPages        []stml.PageSpec
	StateDiagrams    []*statemachine.StateDiagram
	HurlEntries      []hurl.HurlEntry
	ProjectFuncSpecs []funcspec.FuncSpec
	FullendPkgSpecs  []funcspec.FuncSpec
	HurlFiles        []string
	ModelDir         string
	StatesDiags      []diagnostic.Diagnostic
}
