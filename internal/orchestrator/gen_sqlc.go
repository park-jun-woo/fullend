//ff:func feature=orchestrator type=command
//ff:what genSqlc runs sqlc code generation for DDL.

package orchestrator

import (
	"fmt"

	"github.com/geul-org/fullend/internal/reporter"
)

func genSqlc(specsDir, artifactsDir string) reporter.StepResult {
	step := reporter.StepResult{Name: "sqlc"}

	// Auto-generate sqlc.yaml if not present.
	configPath, err := generateSqlcConfig(specsDir, artifactsDir)
	if err != nil {
		step.Status = reporter.Skip
		step.Summary = err.Error()
		return step
	}

	res := RunExec("sqlc", "generate", "-f", configPath)
	if res.Skipped {
		step.Status = reporter.Skip
		step.Summary = "sqlc 미설치, 스킵"
		step.Errors = append(step.Errors, "[WARN] sqlc가 설치되어 있지 않습니다 — go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest")
		return step
	}
	if res.Err != nil {
		step.Status = reporter.Skip
		step.Summary = "sqlc generate 실패, 스킵"
		step.Errors = append(step.Errors, fmt.Sprintf("[WARN] %v", res.Err))
		if res.Stderr != "" {
			step.Errors = append(step.Errors, res.Stderr)
		}
		return step
	}
	step.Status = reporter.Pass
	step.Summary = "DB models generated"
	return step
}
