//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateStaleResponse — @put/@delete 후 갱신 없이 @response 사용 WARNING (S-36)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateStaleResponse(fn parsessac.ServiceFunc) []validate.ValidationError {
	mutated := make(map[string]bool)
	refreshed := make(map[string]bool)
	for _, seq := range fn.Sequences {
		if seq.Type == "put" || seq.Type == "delete" {
			mutated[extractModelFromSeq(seq)] = true
		}
		if seq.Type == "get" && seq.Result != nil {
			refreshed[seq.Result.Type] = true
		}
	}
	var errs []validate.ValidationError
	for _, seq := range fn.Sequences {
		if seq.Type != "response" {
			continue
		}
		errs = append(errs, checkStaleFields(fn, seq, mutated, refreshed)...)
	}
	return errs
}
