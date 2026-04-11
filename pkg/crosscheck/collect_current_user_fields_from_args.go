//ff:func feature=crosscheck type=util control=iteration dimension=1
//ff:what collectCurrentUserFieldsFromArgs — Args에서 currentUser 필드 수집
package crosscheck

import "github.com/park-jun-woo/fullend/pkg/parser/ssac"

func collectCurrentUserFieldsFromArgs(args []ssac.Arg, seen map[string]bool) {
	for _, arg := range args {
		if arg.Source == "currentUser" && arg.Field != "" {
			seen[arg.Field] = true
		}
	}
}
