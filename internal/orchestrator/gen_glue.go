//ff:func feature=orchestrator type=command control=sequence
//ff:what genGlue generates glue code (Server struct, main.go, frontend setup).

package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/gen"
	"github.com/park-jun-woo/fullend/internal/genapi"
	"github.com/park-jun-woo/fullend/internal/reporter"
)

func genGlue(specsDir, artifactsDir string, has map[SSOTKind]DetectedSSOT, parsed *genapi.ParsedSSOTs, stmlDeps map[string]string, stmlPages []string, stmlPageOps map[string]string) reporter.StepResult {
	step := reporter.StepResult{Name: "glue-gen"}

	modulePath := determineModulePath(specsDir, artifactsDir, parsed.Config)

	cfg := &genapi.GenConfig{
		ArtifactsDir: artifactsDir,
		SpecsDir:     specsDir,
		ModulePath:   modulePath,
	}

	var stmlOut *genapi.STMLGenOutput
	if stmlDeps != nil || stmlPages != nil || stmlPageOps != nil {
		stmlOut = &genapi.STMLGenOutput{
			Deps:    stmlDeps,
			Pages:   stmlPages,
			PageOps: stmlPageOps,
		}
	}

	if err := gen.Generate(parsed, cfg, stmlOut); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("glue-gen error: %v", err))
		return step
	}

	step.Status = reporter.Pass
	step.Summary = "server + main.go + frontend setup generated"
	return step
}
