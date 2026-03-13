package orchestrator

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

// ParsedSSOTs holds the results of parsing all detected SSOTs once.
// Consumers receive pointers to these results instead of re-parsing.
type ParsedSSOTs struct {
	Config       *projectconfig.ProjectConfig
	OpenAPIDoc   *openapi3.T
	SymbolTable  *ssacvalidator.SymbolTable
	ServiceFuncs []ssacparser.ServiceFunc
	STMLPages    []stmlparser.PageSpec
	States       []*statemachine.StateDiagram
	Policies     []*policy.Policy
	FuncSpecs    []funcspec.FuncSpec // project func/
	PkgFuncSpecs []funcspec.FuncSpec // fullend pkg/
	HurlFiles    []string
	ModelDir     string
}

// ParseAll parses all detected SSOTs once and returns the cached results.
// Errors during parsing are recorded in the returned slices (nil values
// indicate parse failure). Skipped kinds are not parsed.
func ParseAll(root string, detected []DetectedSSOT, skip map[SSOTKind]bool) *ParsedSSOTs {
	p := &ParsedSSOTs{}

	has := make(map[SSOTKind]DetectedSSOT)
	for _, d := range detected {
		has[d.Kind] = d
	}

	if _, ok := has[KindConfig]; ok && !skip[KindConfig] {
		cfg, err := projectconfig.Load(root)
		if err == nil {
			p.Config = cfg
		}
	}

	if d, ok := has[KindOpenAPI]; ok && !skip[KindOpenAPI] {
		doc, err := openapi3.NewLoader().LoadFromFile(d.Path)
		if err == nil {
			p.OpenAPIDoc = doc
		}
	}

	if _, ok := has[KindDDL]; ok && !skip[KindDDL] {
		st, err := ssacvalidator.LoadSymbolTable(root)
		if err == nil {
			p.SymbolTable = st
		}
	}

	if d, ok := has[KindSSaC]; ok && !skip[KindSSaC] {
		funcs, err := ssacparser.ParseDir(d.Path)
		if err == nil {
			p.ServiceFuncs = funcs
		}
	}

	if d, ok := has[KindSTML]; ok && !skip[KindSTML] {
		pages, err := stmlparser.ParseDir(d.Path)
		if err == nil {
			p.STMLPages = pages
		}
	}

	if d, ok := has[KindStates]; ok && !skip[KindStates] {
		diagrams, err := statemachine.ParseDir(d.Path)
		if err == nil {
			p.States = diagrams
		}
	}

	if d, ok := has[KindPolicy]; ok && !skip[KindPolicy] {
		policies, err := policy.ParseDir(d.Path)
		if err == nil {
			p.Policies = policies
		}
	}

	if d, ok := has[KindFunc]; ok && !skip[KindFunc] {
		specs, err := funcspec.ParseDir(d.Path)
		if err == nil {
			p.FuncSpecs = specs
		}
	}

	// fullend built-in pkg/ specs.
	if pkgRoot := findFullendPkgRoot(); pkgRoot != "" {
		if specs, err := funcspec.ParseDir(pkgRoot); err == nil {
			p.PkgFuncSpecs = specs
		}
	}

	if d, ok := has[KindModel]; ok {
		p.ModelDir = d.Path
	}

	return p
}
