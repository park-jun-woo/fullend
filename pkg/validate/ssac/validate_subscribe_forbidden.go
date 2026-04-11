//ff:func feature=rule type=rule control=sequence
//ff:what validateSubscribeForbidden — @subscribe에서 request/query 금지 (S-42~S-43), HTTP에서 message 금지 (S-44)
package ssac

import (
	parsessac "github.com/park-jun-woo/fullend/pkg/parser/ssac"
	"github.com/park-jun-woo/fullend/pkg/rule"
	"github.com/park-jun-woo/fullend/pkg/validate"
)

func validateSubscribeForbidden(fn parsessac.ServiceFunc, ground *rule.Ground) []validate.ValidationError {
	if fn.Subscribe != nil {
		return validateSubForbiddenRefs(fn, ground)
	}
	return validateHTTPForbiddenRefs(fn, ground)
}
