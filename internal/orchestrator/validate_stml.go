//ff:func feature=orchestrator type=rule control=iteration dimension=1
//ff:what STML 페이지 스펙 검증 — 바인딩 수 집계 + stml validator 실행
package orchestrator

import (
	"fmt"

	"github.com/park-jun-woo/fullend/internal/reporter"
	stmlparser "github.com/park-jun-woo/fullend/internal/stml/parser"
	stmlvalidator "github.com/park-jun-woo/fullend/internal/stml/validator"
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

	verrs := stmlvalidator.Validate(pages, root)
	if len(verrs) > 0 {
		step.Status = reporter.Fail
		for _, ve := range verrs {
			step.Errors = append(step.Errors, fmt.Sprintf("%s [%s] — %s",
				ve.File, ve.Attr, ve.Message))
		}
	} else {
		step.Status = reporter.Pass
	}
	step.Summary = fmt.Sprintf("%d pages, %d bindings", len(pages), bindings)
	return step
}
