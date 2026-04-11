//ff:func feature=crosscheck type=rule control=sequence
//ff:what checkEndpointSecurity — endpoint security 사용 → middleware 설정 필수 (X-52)
package crosscheck

import (
	"github.com/park-jun-woo/fullend/pkg/parser/fullend"
)

func checkEndpointSecurity(fs *fullend.Fullstack) []CrossError {
	if fs.OpenAPIDoc == nil || fs.Manifest == nil {
		return nil
	}
	if !hasOpenAPISecurity(fs) {
		return nil
	}
	if len(fs.Manifest.Backend.Middleware) == 0 {
		return []CrossError{{Rule: "X-52", Context: "fullend.yaml", Level: "ERROR",
			Message: "OpenAPI endpoints use security but backend.middleware not configured"}}
	}
	return nil
}
