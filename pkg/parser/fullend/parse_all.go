//ff:func feature=orchestrator type=command control=iteration dimension=1
//ff:what 탐지된 모든 SSOT를 1회 파싱하여 Fullstack에 담아 반환
package fullend

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/pkg/parser/funcspec"
	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
	"github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/parser/statemachine"
	"github.com/park-jun-woo/fullend/pkg/parser/stml"
)

// ParseAll parses all detected SSOTs once and returns the results.
// Skipped kinds are not parsed. Parse errors result in nil/zero fields.
func ParseAll(root string, detected []DetectedSSOT, skip map[SSOTKind]bool) *Fullstack {
	fs := &Fullstack{}

	has := make(map[SSOTKind]DetectedSSOT)
	for _, d := range detected {
		has[d.Kind] = d
	}

	if _, ok := has[KindConfig]; ok && !skip[KindConfig] {
		cfg, err := manifest.Load(root)
		if err == nil {
			fs.Config = cfg
		}
	}

	if d, ok := has[KindOpenAPI]; ok && !skip[KindOpenAPI] {
		doc, err := openapi3.NewLoader().LoadFromFile(d.Path)
		if err == nil {
			fs.OpenAPIDoc = doc
		}
	}

	if d, ok := has[KindSSaC]; ok && !skip[KindSSaC] {
		funcs, err := ssac.ParseDir(d.Path)
		if err == nil {
			fs.ServiceFuncs = funcs
		}
	}

	if d, ok := has[KindSTML]; ok && !skip[KindSTML] {
		pages, err := stml.ParseDir(d.Path)
		if err == nil {
			fs.STMLPages = pages
		}
	}

	if d, ok := has[KindStates]; ok && !skip[KindStates] {
		diagrams, err := statemachine.ParseDir(d.Path)
		if err == nil {
			fs.StateDiagrams = diagrams
		} else {
			fs.StatesErr = err
		}
	}

	if d, ok := has[KindDDL]; ok && !skip[KindDDL] {
		fs.DDLResults = parseDDLDir(d.Path)
	}

	if d, ok := has[KindPolicy]; ok && !skip[KindPolicy] {
		fs.Policies = parseRegoDir(d.Path)
	}

	if d, ok := has[KindFunc]; ok && !skip[KindFunc] {
		specs, err := funcspec.ParseDir(d.Path)
		if err == nil {
			fs.ProjectFuncSpecs = specs
		}
	}

	if d, ok := has[KindModel]; ok {
		fs.ModelDir = d.Path
	}

	return fs
}
