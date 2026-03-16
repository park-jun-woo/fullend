//ff:func feature=crosscheck type=util control=sequence topic=config-check
//ff:what JWT 빌트인 함수 목록 정의
package crosscheck

// jwtBuiltinFuncs are claims-dependent functions that are generated into internal/auth/
// (not in pkg/auth) when auth.type is jwt. Skip funcspec lookup for these.
var jwtBuiltinFuncs = map[string]bool{
	"auth.issueToken":   true,
	"auth.verifyToken":  true,
	"auth.refreshToken": true,
}
