//ff:func feature=crosscheck type=rule control=iteration dimension=1 topic=authz-check
//ff:what SSaC @auth 입력 필드가 authz CheckRequest에 존재하는지 검증
package crosscheck

import (
	ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"
)

// defaultCheckRequestFields are the fields of the default pkg/authz CheckRequest struct.
var defaultCheckRequestFields = map[string]bool{
	"Action":     true,
	"Resource":   true,
	"UserID":     true,
	"Role":       true,
	"ResourceID": true,
}

// CheckAuthz validates @auth inputs against the authz CheckRequest fields.
func CheckAuthz(funcs []ssacparser.ServiceFunc, authzPackage string) []CrossError {
	var errs []CrossError

	// If custom authz package is set, skip validation (we don't have its source to check).
	if authzPackage != "" {
		return errs
	}

	for _, fn := range funcs {
		errs = append(errs, checkAuthzFunc(fn, defaultCheckRequestFields)...)
	}

	return errs
}
