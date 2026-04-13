//ff:func feature=orchestrator type=command control=sequence
//ff:what genGlue generates glue code via pkg/generate (gogin + react + hurl).

package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/reporter"
	"github.com/park-jun-woo/fullend/pkg/fullend"
	pkggen "github.com/park-jun-woo/fullend/pkg/generate"
	"github.com/park-jun-woo/fullend/pkg/generate/react"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

func genGlue(specsDir, artifactsDir string, fs *fullend.Fullstack, g *rule.Ground, stmlDeps map[string]string, stmlPages []string, stmlPageOps map[string]string) reporter.StepResult {
	step := reporter.StepResult{Name: "glue-gen"}

	modulePath := determinePkgModulePath(specsDir, artifactsDir, fs.Manifest)

	cfg := &pkggen.Config{
		ArtifactsDir: artifactsDir,
		SpecsDir:     specsDir,
		ModulePath:   modulePath,
	}

	var stmlOut *react.STMLGenOutput
	if stmlDeps != nil || stmlPages != nil || stmlPageOps != nil {
		stmlOut = &react.STMLGenOutput{
			Deps:    stmlDeps,
			Pages:   stmlPages,
			PageOps: stmlPageOps,
		}
	}

	if err := pkggen.Generate(fs, g, cfg, stmlOut); err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("glue-gen error: %v", err))
		return step
	}

	step.Status = reporter.Pass
	step.Summary = "server + main.go + frontend setup generated"
	return step
}
