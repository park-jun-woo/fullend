//ff:func feature=crosscheck type=util control=iteration dimension=1 topic=config-check
//ff:what SSaC 시퀀스에서 currentUser.X 참조 필드를 수집
package crosscheck

import ssacparser "github.com/park-jun-woo/fullend/internal/ssac/parser"

// collectCurrentUserFields scans all SSaC sequences for currentUser.X references.
// Returns map[fieldName][]location.
func collectCurrentUserFields(funcs []ssacparser.ServiceFunc) map[string][]string {
	result := make(map[string][]string)

	for _, sf := range funcs {
		loc := sf.FileName + ":" + sf.Name
		for _, seq := range sf.Sequences {
			collectCurrentUserFromInputs(seq, loc, result)
		}
	}

	return result
}
