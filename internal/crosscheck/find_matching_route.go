//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=scenario-check
//ff:what Hurl 세그먼트와 메서드에 매칭되는 라우트 검색
package crosscheck

// findMatchingRoute finds a matching route for hurl segments and method.
func findMatchingRoute(hurlSegs []string, method string, routes []apiRoute) (matched bool, matchedRoute *apiRoute) {
	for i := range routes {
		if !segmentsMatch(hurlSegs, routes[i].Segments) {
			continue
		}
		matched = true
		if routes[i].Method == method {
			matchedRoute = &routes[i]
			return
		}
	}
	return
}
