//ff:func feature=orchestrator type=command control=iteration dimension=3
//ff:what Status collects SSOT stats and returns lines to display.

package orchestrator

import (
	"fmt"
	"path/filepath"
)

// Status collects SSOT stats and returns lines to display.
func Status(root string, detected []DetectedSSOT) []StatusLine {
	// Parse all SSOTs once.
	parsed := ParseAll(root, detected, nil)

	var lines []StatusLine

	absRoot, _ := filepath.Abs(root)

	for _, d := range detected {
		relPath, err := filepath.Rel(absRoot, d.Path)
		if err != nil || relPath == "" {
			relPath = d.Path
		}

		switch d.Kind {
		case KindOpenAPI:
			summary := "?"
			if parsed.OpenAPIDoc != nil {
				summary = fmt.Sprintf("%d endpoints", countEndpoints(parsed.OpenAPIDoc))
			}
			lines = append(lines, StatusLine{Kind: KindOpenAPI, Path: relPath, Summary: summary})
		case KindDDL:
			summary := "?"
			if parsed.SymbolTable != nil {
				summary = fmt.Sprintf("%d tables, %d columns", len(parsed.SymbolTable.DDLTables), countDDLColumns(parsed.SymbolTable.DDLTables))
			}
			lines = append(lines, StatusLine{Kind: KindDDL, Path: relPath, Summary: summary})
		case KindSSaC:
			summary := "?"
			if parsed.ServiceFuncs != nil {
				summary = fmt.Sprintf("%d functions", len(parsed.ServiceFuncs))
			}
			lines = append(lines, StatusLine{Kind: KindSSaC, Path: relPath, Summary: summary})
		case KindSTML:
			summary := "?"
			if parsed.STMLPages != nil {
				summary = fmt.Sprintf("%d pages", len(parsed.STMLPages))
			}
			lines = append(lines, StatusLine{Kind: KindSTML, Path: relPath, Summary: summary})
		case KindStates:
			summary := "?"
			if parsed.StateDiagrams != nil {
				summary = fmt.Sprintf("%d diagrams, %d transitions", len(parsed.StateDiagrams), countTransitions(parsed.StateDiagrams))
			}
			lines = append(lines, StatusLine{Kind: KindStates, Path: relPath, Summary: summary})
		case KindPolicy:
			summary := "?"
			if parsed.Policies != nil {
				summary = fmt.Sprintf("%d files, %d rules", len(parsed.Policies), countPolicyRules(parsed.Policies))
			}
			lines = append(lines, StatusLine{Kind: KindPolicy, Path: relPath, Summary: summary})
		case KindScenario:
			scenarioHurls, _ := filepath.Glob(filepath.Join(d.Path, "scenario-*.hurl"))
			invariantHurls, _ := filepath.Glob(filepath.Join(d.Path, "invariant-*.hurl"))
			total := len(scenarioHurls) + len(invariantHurls)
			lines = append(lines, StatusLine{Kind: KindScenario, Path: relPath, Summary: fmt.Sprintf("%d hurl files", total)})
		case KindFunc:
			summary := "?"
			if parsed.ProjectFuncSpecs != nil {
				stubs := countFuncStubs(parsed.ProjectFuncSpecs)
				if stubs > 0 {
					summary = fmt.Sprintf("%d funcs (%d TODO)", len(parsed.ProjectFuncSpecs), stubs)
				} else {
					summary = fmt.Sprintf("%d funcs", len(parsed.ProjectFuncSpecs))
				}
			}
			lines = append(lines, StatusLine{Kind: KindFunc, Path: relPath, Summary: summary})
		case KindModel:
			// Model is auxiliary; skip in status display.
		}
	}

	return lines
}
