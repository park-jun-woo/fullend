//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validatePaginationRefs — Page/Cursor 타입 ↔ x-pagination 존재 검증 (S-54~S-56)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validatePaginationRefs(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	if fn.Subscribe != nil {
		return nil
	}
	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		if seq.Type != "get" || seq.Result == nil || seq.Result.Wrapper == "" {
			continue
		}
		if !ground.Config["pagination."+fn.Name] {
			errs = append(errs, validate.ValidationError{
				Rule: "S-54", File: fn.FileName, Func: fn.Name, SeqIdx: i, Level: "ERROR",
				Message: seq.Result.Wrapper + " type used but no x-pagination in OpenAPI",
			})
		}
	}
	return errs
}
