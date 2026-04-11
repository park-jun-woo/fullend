//ff:func feature=rule type=rule control=iteration dimension=1
//ff:what validateFetch — fetch ブロックの param/bind/each/paginate 検証 (TM-4~TM-12)
package stml

import (
	parsestml "github.com/park-jun-woo/fullend/pkg/parser/stml"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateFetch(pages []parsestml.PageSpec, ground *rule.Ground) []validate.ValidationError {
	var errs []validate.ValidationError
	for _, page := range pages {
		for _, fb := range page.Fetches {
			errs = append(errs, validateFetchBlock(fb, page.FileName, ground)...)
		}
	}
	return errs
}
