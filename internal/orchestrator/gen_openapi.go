//ff:func feature=orchestrator type=command control=sequence
//ff:what genOpenAPI runs oapi-codegen to generate types and server from OpenAPI spec.

package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/fullend/internal/reporter"
)

func genOpenAPI(specsDir, artifactsDir string) reporter.StepResult {
	step := reporter.StepResult{Name: "oapi-gen"}
	apiPath := filepath.Join(specsDir, "api", "openapi.yaml")

	outDir := filepath.Join(artifactsDir, "backend", "internal", "api")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("cannot create dir: %v", err))
		return step
	}

	// Generate types.
	typesOut := filepath.Join(outDir, "types.gen.go")
	res := RunExec("oapi-codegen", "-package", "api", "-generate", "types", "-o", typesOut, apiPath)
	if res.Skipped {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, "oapi-codegen이 설치되어 있지 않습니다 — go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest")
		return step
	}
	if res.Err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("oapi-codegen types error: %v", res.Err))
		if res.Stderr != "" {
			step.Errors = append(step.Errors, res.Stderr)
		}
		return step
	}

	// Generate server (net/http std).
	serverOut := filepath.Join(outDir, "server.gen.go")
	res = RunExec("oapi-codegen", "-package", "api", "-generate", "std-http-server", "-o", serverOut, apiPath)
	if res.Err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("oapi-codegen server error: %v", res.Err))
		if res.Stderr != "" {
			step.Errors = append(step.Errors, res.Stderr)
		}
		return step
	}

	step.Status = reporter.Pass
	step.Summary = "types + server generated"
	return step
}
