//ff:func feature=orchestrator type=rule control=iteration dimension=1
//ff:what STML 페이지 스펙 검증 — pkg/validate/stml toulmin 기반
package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/reporter"
	stmlparser "github.com/park-jun-woo/fullend/internal/stml/parser"
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
	pkgstml "github.com/park-jun-woo/fullend/pkg/validate/stml"
)

func validateSTML(root string, pages []stmlparser.PageSpec) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindSTML)}
	if pages == nil {
		step.Status = reporter.Fail
		step.Errors = append(step.Errors, "STML parse failed")
		return step
	}

	bindings := 0
	for _, p := range pages {
		bindings += len(p.Fetches) + len(p.Actions)
	}

	detected, _ := fullend.DetectSSOTs(root)
	fs := fullend.ParseAll(root, detected, nil)
	ground := buildSTMLGround(fs)

	verrs := pkgstml.Validate(fs.STMLPages, ground)
	if len(verrs) > 0 {
		step.Status = reporter.Fail
		for _, ve := range verrs {
			step.Errors = append(step.Errors, fmt.Sprintf("%s [%s] — %s",
				ve.File, ve.Rule, ve.Message))
		}
	} else {
		step.Status = reporter.Pass
	}
	step.Summary = fmt.Sprintf("%d pages, %d bindings", len(pages), bindings)
	return step
}
