//ff:func feature=orchestrator type=command
//ff:what Chain traces all SSOT nodes connected to the given operationId.

package orchestrator

import (
	"fmt"
	"path/filepath"

	ssacparser "github.com/geul-org/fullend/internal/ssac/parser"
)

// Chain traces all SSOT nodes connected to the given operationId.
func Chain(specsDir string, operationID string) ([]ChainLink, error) {
	abs, err := filepath.Abs(specsDir)
	if err != nil {
		return nil, err
	}

	detected, err := DetectSSOTs(abs)
	if err != nil {
		return nil, err
	}

	has := make(map[SSOTKind]DetectedSSOT)
	for _, d := range detected {
		has[d.Kind] = d
	}

	// Parse all available SSOTs once.
	parsed := ParseAll(abs, detected, nil)

	allFuncSpecs := append(parsed.ProjectFuncSpecs, parsed.FullendPkgSpecs...)

	// Trace the chain.
	var links []ChainLink

	// 1. OpenAPI
	if parsed.OpenAPIDoc != nil {
		link := traceOpenAPI(parsed.OpenAPIDoc, operationID, abs)
		if link != nil {
			links = append(links, *link)
		} else {
			return nil, fmt.Errorf("operationId %q not found in OpenAPI", operationID)
		}
	}

	// Find the matching SSaC function.
	var matchedFunc *ssacparser.ServiceFunc
	for i := range parsed.ServiceFuncs {
		if parsed.ServiceFuncs[i].Name == operationID {
			matchedFunc = &parsed.ServiceFuncs[i]
			break
		}
	}

	// 2. SSaC
	if matchedFunc != nil {
		links = append(links, traceSSaC(matchedFunc, abs))
	}

	// 3. DDL — trace tables referenced by SSaC sequences
	if matchedFunc != nil && parsed.SymbolTable != nil {
		ddlLinks := traceDDL(matchedFunc, parsed.SymbolTable, abs)
		links = append(links, ddlLinks...)
	}

	// 4. Rego — trace policies referenced by @auth sequences
	if matchedFunc != nil && parsed.Policies != nil {
		regoLinks := tracePolicy(matchedFunc, parsed.Policies, abs)
		links = append(links, regoLinks...)
	}

	// 5. StateDiagram — trace diagrams referenced by @state sequences
	if matchedFunc != nil && parsed.StateDiagrams != nil {
		stateLinks := traceStates(matchedFunc, parsed.StateDiagrams, abs)
		links = append(links, stateLinks...)
	}

	// 6. FuncSpec — trace funcs referenced by @call sequences
	if matchedFunc != nil && len(allFuncSpecs) > 0 {
		funcLinks := traceFuncSpecs(matchedFunc, allFuncSpecs, abs)
		links = append(links, funcLinks...)
	}

	// 7. Hurl scenario — trace .hurl files referencing this endpoint
	if d, ok := has[KindScenario]; ok {
		hurlLinks := traceHurlScenarios(operationID, parsed.OpenAPIDoc, d.Path, abs)
		links = append(links, hurlLinks...)
	}

	// 8. STML — trace frontend files referencing this endpoint
	if parsed.OpenAPIDoc != nil {
		if d, ok := has[KindSTML]; ok {
			stmlLinks := traceSTML(parsed.OpenAPIDoc, operationID, d.Path, abs)
			links = append(links, stmlLinks...)
		}
	}

	// 9. Artifacts — trace generated code referencing this operationId
	if matchedFunc != nil {
		artifactsDir := inferArtifactsDir(abs)
		if artifactsDir != "" {
			artifactLinks := traceArtifacts(artifactsDir, operationID, matchedFunc)
			links = append(links, artifactLinks...)
		}
	}

	return links, nil
}
