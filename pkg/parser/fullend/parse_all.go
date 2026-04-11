//ff:func feature=orchestrator type=command control=iteration dimension=1
//ff:what 탐지된 모든 SSOT를 1회 파싱하여 Fullstack에 담아 반환
package fullend

import (
	"github.com/getkin/kin-openapi/openapi3"

	"github.com/park-jun-woo/fullend/pkg/parser/ddl"
	"github.com/park-jun-woo/fullend/pkg/parser/funcspec"
	"github.com/park-jun-woo/fullend/pkg/parser/hurl"
	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
	"github.com/park-jun-woo/fullend/pkg/parser/rego"
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
		cfg, diags := manifest.Load(root)
		if len(diags) == 0 {
			fs.Manifest = cfg
		}
	}

	if d, ok := has[KindOpenAPI]; ok && !skip[KindOpenAPI] {
		doc, err := openapi3.NewLoader().LoadFromFile(d.Path)
		if err == nil {
			fs.OpenAPIDoc = doc
		}
	}

	if d, ok := has[KindSSaC]; ok && !skip[KindSSaC] {
		funcs, diags := ssac.ParseDir(d.Path)
		if len(diags) == 0 {
			fs.ServiceFuncs = funcs
		}
	}

	if d, ok := has[KindSTML]; ok && !skip[KindSTML] {
		pages, diags := stml.ParseDir(d.Path)
		if len(diags) == 0 {
			fs.STMLPages = pages
		}
	}

	if d, ok := has[KindStates]; ok && !skip[KindStates] {
		diagrams, diags := statemachine.ParseDir(d.Path)
		if len(diags) == 0 {
			fs.StateDiagrams = diagrams
		} else {
			fs.StatesDiags = diags
		}
	}

	if d, ok := has[KindDDL]; ok && !skip[KindDDL] {
		results, diags := ddl.ParseDir(d.Path)
		if len(diags) == 0 {
			fs.DDLResults = results
		}
		tables, tdiags := ddl.ParseTables(d.Path)
		if len(tdiags) == 0 {
			fs.DDLTables = tables
		}
	}

	if d, ok := has[KindPolicy]; ok && !skip[KindPolicy] {
		modules, diags := rego.ParseDir(d.Path)
		if len(diags) == 0 {
			fs.Policies = modules
		}
		policies, pdiags := rego.ParsePolicies(d.Path)
		if len(pdiags) == 0 {
			fs.ParsedPolicies = policies
		}
	}

	if d, ok := has[KindFunc]; ok && !skip[KindFunc] {
		specs, diags := funcspec.ParseDir(d.Path)
		if len(diags) == 0 {
			fs.ProjectFuncSpecs = specs
		}
	}

	if d, ok := has[KindScenario]; ok && !skip[KindScenario] {
		fs.HurlFiles = hurl.CollectFiles(d.Path)
		for _, hf := range fs.HurlFiles {
			entries, _ := hurl.ParseFile(hf)
			fs.HurlEntries = append(fs.HurlEntries, entries...)
		}
	}

	if d, ok := has[KindToulmin]; ok && !skip[KindToulmin] {
		fs.TanglFiles = parseTanglDir(d.Path)
	}

	if d, ok := has[KindModel]; ok {
		fs.ModelDir = d.Path
	}

	// fullend built-in pkg/ specs.
	if pkgRoot := findFullendPkgRoot(); pkgRoot != "" {
		specs, diags := funcspec.ParseDir(pkgRoot)
		if len(diags) == 0 {
			fs.FullendPkgSpecs = specs
		}
	}

	return fs
}
