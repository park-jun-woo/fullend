//ff:func feature=crosscheck type=command control=sequence
//ff:what Run — Fullstack에서 모든 교차 검증을 실행하여 CrossError 목록 반환
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/fullend"

// Run executes all cross-validation rules against parsed SSOTs.
func Run(fs *fullend.Fullstack) []CrossError {
	g := BuildGround(fs)
	var errs []CrossError
	errs = append(errs, checkSSaCOpenAPI(g, fs)...)
	errs = append(errs, checkOpenAPISSaC(g, fs)...)
	errs = append(errs, checkOpenAPIDDL(g, fs)...)
	errs = append(errs, checkStates(g, fs)...)
	errs = append(errs, checkSSaCStates(g, fs)...)
	errs = append(errs, checkPolicy(g, fs)...)
	errs = append(errs, checkOwnership(g, fs)...)
	errs = append(errs, checkHurl(g, fs)...)
	errs = append(errs, checkFuncs(g, fs)...)
	errs = append(errs, checkFuncCoverage(g, fs)...)
	errs = append(errs, checkFuncPurity(fs)...)
	errs = append(errs, checkConfig(g, fs)...)
	errs = append(errs, checkClaims(g, fs)...)
	errs = append(errs, checkClaimsRego(g, fs)...)
	errs = append(errs, checkRoles(g, fs)...)
	errs = append(errs, checkQueue(g)...)
	errs = append(errs, checkMiddleware(g, fs)...)
	errs = append(errs, checkDDLCoverage(g)...)
	errs = append(errs, checkSensitive(fs)...)
	return errs
}
