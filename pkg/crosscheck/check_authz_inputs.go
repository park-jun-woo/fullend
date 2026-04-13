//ff:func feature=crosscheck type=rule control=iteration dimension=1
//ff:what checkAuthzInputs — @auth inputs → authz CheckRequest 필드 검증 (X-60)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/fullend"
	"github.com/park-jun-woo/fullend/pkg/rule"
)

var authzCheckRequestFields = rule.StringSet{
	"Action": true, "Resource": true, "UserID": true, "Role": true, "ResourceID": true,
}

func checkAuthzInputs(fs *fullend.Fullstack) []CrossError {
	var errs []CrossError
	for _, fn := range fs.ServiceFuncs {
		errs = append(errs, checkAuthzInputSeqs(fn.Name, fn.Sequences)...)
	}
	return errs
}
