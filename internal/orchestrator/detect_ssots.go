//ff:func feature=orchestrator type=command
//ff:what DetectSSOTs scans root for known SSOT directories and returns what exists.

package orchestrator

import (
	"os"
	"path/filepath"
)

// DetectSSOTs scans root for known SSOT directories and returns what exists.
func DetectSSOTs(root string) ([]DetectedSSOT, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, &NotDirError{Path: abs}
	}

	var found []DetectedSSOT

	checks := []struct {
		kind    SSOTKind
		pattern string
	}{
		{KindOpenAPI, "api/openapi.yaml"},
		{KindDDL, "db/*.sql"},
		{KindSSaC, "service/*.ssac"},
		{KindModel, "model/*.go"},
		{KindSTML, "frontend/*.html"},
	}

	for _, c := range checks {
		matches, _ := filepath.Glob(filepath.Join(abs, c.pattern))
		if len(matches) > 0 {
			dir := filepath.Dir(matches[0])
			if c.kind == KindOpenAPI {
				dir = matches[0] // file path, not dir
			}
			found = append(found, DetectedSSOT{Kind: c.kind, Path: dir})
		} else if c.kind == KindSSaC {
			// Also check for domain folder structure: service/{domain}/*.ssac
			subMatches, _ := filepath.Glob(filepath.Join(abs, "service", "*", "*.ssac"))
			if len(subMatches) > 0 {
				found = append(found, DetectedSSOT{Kind: c.kind, Path: filepath.Join(abs, "service")})
			}
		}
	}

	// Check for states/ directory (Mermaid stateDiagram files).
	statesDir := filepath.Join(abs, "states")
	if statesMatches, _ := filepath.Glob(filepath.Join(statesDir, "*.md")); len(statesMatches) > 0 {
		found = append(found, DetectedSSOT{Kind: KindStates, Path: statesDir})
	}

	// Check for policy/ directory (OPA Rego files).
	policyDir := filepath.Join(abs, "policy")
	if policyMatches, _ := filepath.Glob(filepath.Join(policyDir, "*.rego")); len(policyMatches) > 0 {
		found = append(found, DetectedSSOT{Kind: KindPolicy, Path: policyDir})
	}

	// Check for tests/ directory (scenario .hurl files).
	testsDir := filepath.Join(abs, "tests")
	scenarioHurls, _ := filepath.Glob(filepath.Join(testsDir, "scenario-*.hurl"))
	invariantHurls, _ := filepath.Glob(filepath.Join(testsDir, "invariant-*.hurl"))
	if len(scenarioHurls)+len(invariantHurls) > 0 {
		found = append(found, DetectedSSOT{Kind: KindScenario, Path: testsDir})
	} else {
		// Also detect deprecated .feature files so validate can emit ERROR.
		scenarioDir := filepath.Join(abs, "scenario")
		if featureFiles, _ := filepath.Glob(filepath.Join(scenarioDir, "*.feature")); len(featureFiles) > 0 {
			found = append(found, DetectedSSOT{Kind: KindScenario, Path: testsDir})
		}
	}

	// Check for fullend.yaml (project config).
	configPath := filepath.Join(abs, "fullend.yaml")
	if _, err := os.Stat(configPath); err == nil {
		found = append(found, DetectedSSOT{Kind: KindConfig, Path: configPath})
	}

	// Check for func/ directory (custom func spec files).
	funcDir := filepath.Join(abs, "func")
	if fi, err := os.Stat(funcDir); err == nil && fi.IsDir() {
		// Check for any .go files in subdirectories.
		funcMatches, _ := filepath.Glob(filepath.Join(funcDir, "*", "*.go"))
		if len(funcMatches) > 0 {
			found = append(found, DetectedSSOT{Kind: KindFunc, Path: funcDir})
		}
	}

	return found, nil
}
