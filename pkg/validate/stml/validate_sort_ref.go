//ff:func feature=rule type=rule control=sequence
//ff:what validateSortRef — data-sort 컬럼이 OpenAPI x-sort allowed에 있는지 검증 (TM-10)
package stml

import (
	"strings"

	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

func validateSortRef(fb parsestml.FetchBlock, file string, ground *rule.Ground) []validate.ValidationError {
	graph := toulmin.NewGraph("stml-sort-" + fb.OperationID)
	graph.Rule(rule.RefExists).With(&rule.RefExistsSpec{
		BaseSpec:  rule.BaseSpec{Rule: "TM-10", Level: "ERROR", Message: "data-sort column not in OpenAPI x-sort allowed"},
		LookupKey: "OpenAPI.sort." + fb.OperationID,
	})
	col := fb.Sort.Column
	if idx := strings.IndexByte(col, ':'); idx >= 0 {
		col = col[:idx]
	}
	ctx := toulmin.NewContext()
	ctx.Set("ground", ground)
	ctx.Set("claim", col)
	results, _ := graph.Evaluate(ctx)
	return toSTMLErrors(results, file, fb.OperationID)
}
