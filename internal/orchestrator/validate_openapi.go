//ff:func feature=orchestrator type=rule
//ff:what OpenAPI 스펙 검증 — 엔드포인트 수 집계 + 경로 파라미터 충돌 검사
package orchestrator

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/geul-org/fullend/internal/reporter"
)

func validateOpenAPI(path string, doc *openapi3.T) reporter.StepResult {
	step := reporter.StepResult{Name: string(KindOpenAPI)}
	if doc == nil {
		// Parse failed in ParseAll; try again for error message.
		var err error
		doc, err = openapi3.NewLoader().LoadFromFile(path)
		if err != nil {
			step.Status = reporter.Fail
			step.Errors = append(step.Errors, fmt.Sprintf("OpenAPI load error: %v", err))
			return step
		}
	}
	count := 0
	for _, pi := range doc.Paths.Map() {
		for range pi.Operations() {
			count++
		}
	}

	// Check path param name conflicts.
	if conflicts := checkPathParamConflicts(doc); len(conflicts) > 0 {
		for _, c := range conflicts {
			step.Errors = append(step.Errors, c)
		}
	}

	if len(step.Errors) > 0 {
		step.Status = reporter.Fail
	} else {
		step.Status = reporter.Pass
	}
	step.Summary = fmt.Sprintf("%d endpoints", count)
	return step
}
