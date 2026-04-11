//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectCurrentUserFields — SSaC에서 currentUser.Field 참조 수집
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/fullend"

func collectCurrentUserFields(fs *fullend.Fullstack) []string {
	seen := make(map[string]bool)
	for _, fn := range fs.ServiceFuncs {
		collectCurrentUserFieldsFromSeqs(fn.Sequences, seen)
	}
	var fields []string
	for f := range seen {
		fields = append(fields, f)
	}
	return fields
}
