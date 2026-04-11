//ff:func feature=crosscheck type=command control=sequence
//ff:what Run — Fullstack에서 모든 교차 검증을 실행하여 CrossError 목록 반환
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/fullend"

// Run executes all cross-validation rules against parsed SSOTs.
func Run(fs *fullend.Fullstack) []CrossError {
	g := BuildGround(fs)
	var errs []CrossError
	// SSaC ↔ OpenAPI
	errs = append(errs, checkSSaCOpenAPI(g, fs)...)
	errs = append(errs, checkOpenAPISSaC(g, fs)...)
	errs = append(errs, checkResponseSchema(g, fs)...)
	errs = append(errs, checkErrStatus(fs)...)
	errs = append(errs, checkShorthandResponse(fs)...)
	// OpenAPI ↔ DDL
	errs = append(errs, checkOpenAPIDDL(g, fs)...)
	errs = append(errs, checkSortIndex(g, fs)...)
	errs = append(errs, checkXInclude(g, fs)...)
	errs = append(errs, checkCursor(fs)...)
	errs = append(errs, checkGhostProperties(g, fs)...)
	// SSaC ↔ DDL
	errs = append(errs, checkSSaCDDL(g, fs)...)
	// States
	errs = append(errs, checkStates(g, fs)...)
	errs = append(errs, checkSSaCStates(g, fs)...)
	errs = append(errs, checkStatesOpenAPI(g, fs)...)
	errs = append(errs, checkStatesGuard(g, fs)...)
	errs = append(errs, checkStatesDDL(g, fs)...)
	// Policy
	errs = append(errs, checkPolicy(g, fs)...)
	errs = append(errs, checkPolicyReverse(g, fs)...)
	errs = append(errs, checkOwnership(g, fs)...)
	errs = append(errs, checkOwnershipAnnotation(g, fs)...)
	errs = append(errs, checkOwnershipVia(g, fs)...)
	// Hurl
	errs = append(errs, checkHurl(g, fs)...)
	errs = append(errs, checkHurlMethod(g, fs)...)
	errs = append(errs, checkHurlStatus(fs)...)
	// Func
	errs = append(errs, checkFuncs(g, fs)...)
	errs = append(errs, checkFuncCoverage(g, fs)...)
	errs = append(errs, checkFuncPurity(fs)...)
	errs = append(errs, checkFuncDetails(g, fs)...)
	errs = append(errs, checkCallFuncName(fs)...)
	errs = append(errs, checkCallTypeMatch(g, fs)...)
	errs = append(errs, checkCallSourceVar(fs)...)
	// Config
	errs = append(errs, checkConfig(g, fs)...)
	errs = append(errs, checkClaims(g, fs)...)
	errs = append(errs, checkClaimsRego(g, fs)...)
	errs = append(errs, checkConfigReverse(g, fs)...)
	errs = append(errs, checkEndpointSecurity(fs)...)
	errs = append(errs, checkJWTClaims(g, fs)...)
	// Roles
	errs = append(errs, checkRoles(g, fs)...)
	errs = append(errs, checkDDLCheckRoles(g, fs)...)
	// Queue
	errs = append(errs, checkQueue(g)...)
	errs = append(errs, checkQueueSchema(g, fs)...)
	// Coverage
	errs = append(errs, checkMiddleware(g, fs)...)
	errs = append(errs, checkDDLCoverage(g)...)
	// Sensitive + Constraints
	errs = append(errs, checkSensitive(fs)...)
	errs = append(errs, checkConstraints(fs)...)
	errs = append(errs, checkDDLOpenAPIConstraints(g, fs)...)
	// Authz
	errs = append(errs, checkAuthzInputs(fs)...)
	return errs
}
