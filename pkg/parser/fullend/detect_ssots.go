//ff:func feature=orchestrator type=command control=sequence
//ff:what specs 루트에서 SSOT 파일/디렉토리를 탐지하여 목록 반환
package fullend

import (
	"fmt"
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
		return nil, fmt.Errorf("not a directory: %s", abs)
	}

	var found []DetectedSSOT

	// fullend.yaml
	configPath := filepath.Join(abs, "fullend.yaml")
	if _, err := os.Stat(configPath); err == nil {
		found = append(found, DetectedSSOT{Kind: KindConfig, Path: configPath})
	}

	// api/openapi.yaml
	openapiPath := filepath.Join(abs, "api", "openapi.yaml")
	if _, err := os.Stat(openapiPath); err == nil {
		found = append(found, DetectedSSOT{Kind: KindOpenAPI, Path: openapiPath})
	}

	// db/*.sql
	if matches, _ := filepath.Glob(filepath.Join(abs, "db", "*.sql")); len(matches) > 0 {
		found = append(found, DetectedSSOT{Kind: KindDDL, Path: filepath.Join(abs, "db")})
	}

	// service/**/*.ssac
	serviceDir := filepath.Join(abs, "service")
	if ssacFlat, _ := filepath.Glob(filepath.Join(serviceDir, "*.ssac")); len(ssacFlat) > 0 {
		found = append(found, DetectedSSOT{Kind: KindSSaC, Path: serviceDir})
	} else if ssacNested, _ := filepath.Glob(filepath.Join(serviceDir, "*", "*.ssac")); len(ssacNested) > 0 {
		found = append(found, DetectedSSOT{Kind: KindSSaC, Path: serviceDir})
	}

	// model/*.go
	modelDir := filepath.Join(abs, "model")
	if matches, _ := filepath.Glob(filepath.Join(modelDir, "*.go")); len(matches) > 0 {
		found = append(found, DetectedSSOT{Kind: KindModel, Path: modelDir})
	}

	// frontend/*.html
	frontendDir := filepath.Join(abs, "frontend")
	if matches, _ := filepath.Glob(filepath.Join(frontendDir, "*.html")); len(matches) > 0 {
		found = append(found, DetectedSSOT{Kind: KindSTML, Path: frontendDir})
	}

	// states/*.md
	statesDir := filepath.Join(abs, "states")
	if matches, _ := filepath.Glob(filepath.Join(statesDir, "*.md")); len(matches) > 0 {
		found = append(found, DetectedSSOT{Kind: KindStates, Path: statesDir})
	}

	// policy/*.rego
	policyDir := filepath.Join(abs, "policy")
	if matches, _ := filepath.Glob(filepath.Join(policyDir, "*.rego")); len(matches) > 0 {
		found = append(found, DetectedSSOT{Kind: KindPolicy, Path: policyDir})
	}

	// tests/scenario-*.hurl + invariant-*.hurl
	testsDir := filepath.Join(abs, "tests")
	scenarioHurls, _ := filepath.Glob(filepath.Join(testsDir, "scenario-*.hurl"))
	invariantHurls, _ := filepath.Glob(filepath.Join(testsDir, "invariant-*.hurl"))
	if len(scenarioHurls)+len(invariantHurls) > 0 {
		found = append(found, DetectedSSOT{Kind: KindScenario, Path: testsDir})
	}

	// func/*/*.go
	funcDir := filepath.Join(abs, "func")
	if matches, _ := filepath.Glob(filepath.Join(funcDir, "*", "*.go")); len(matches) > 0 {
		found = append(found, DetectedSSOT{Kind: KindFunc, Path: funcDir})
	}

	return found, nil
}
