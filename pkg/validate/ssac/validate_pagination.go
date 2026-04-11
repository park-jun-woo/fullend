//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validatePagination — x-pagination ↔ SSaC query/Page/Cursor 정합성 (S-52~S-56)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validatePagination(fn parsessac.ServiceFunc) []validate.ValidationError {
	if fn.Subscribe != nil {
		return nil
	}
	var errs []validate.ValidationError
	for i, seq := range fn.Sequences {
		if seq.Type != "get" || seq.Result == nil {
			continue
		}
		if seq.Result.Wrapper == "Page" || seq.Result.Wrapper == "Cursor" {
			errs = append(errs, checkPaginationQuery(fn.FileName, fn.Name, i, seq)...)
		}
	}
	return errs
}
