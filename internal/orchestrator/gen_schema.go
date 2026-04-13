//ff:func feature=orchestrator type=command control=sequence
//ff:what genSchema — DDL 위상정렬 + schema.sql 통합 산출 (Phase018)

package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/reporter"
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/generate/db"
)

func genSchema(specsDDLDir, artifactsDir string, fs *fullend.Fullstack) reporter.StepResult {
	step := reporter.StepResult{Name: "schema-gen"}
	if len(fs.DDLTables) == 0 {
		step.Status = reporter.Skip
		step.Summary = "no DDL tables"
		return step
	}
	autoSeed := false
	if fs.Manifest != nil && fs.Manifest.Backend.DB != nil {
		autoSeed = fs.Manifest.Backend.DB.AutoNobodySeed
	}
	cfg := db.Config{
		SpecsDDLDir:    specsDDLDir,
		ArtifactsDir:   artifactsDir,
		AutoNobodySeed: autoSeed,
	}
	tables, seeds, err := db.GenerateSchema(fs.DDLTables, cfg)
	if err != nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, fmt.Sprintf("schema-gen error: %v", err))
		return step
	}
	step.Status = reporter.Pass
	step.Summary = fmt.Sprintf("schema.sql generated (%d tables, %d seeds)", tables, seeds)
	return step
}
