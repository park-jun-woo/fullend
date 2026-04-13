//ff:func feature=funcspec type=util control=sequence
//ff:what firstResultIsPointer — FuncType Results 의 첫 반환이 *T 형태인지

package funcspec

import "go/ast"

// firstResultIsPointer reports whether a FuncType's first result is a pointer type.
// Used to detect @empty-safe funcs for nilable checks.
func firstResultIsPointer(ft *ast.FuncType) bool {
	if ft == nil || ft.Results == nil || len(ft.Results.List) == 0 {
		return false
	}
	_, ok := ft.Results.List[0].Type.(*ast.StarExpr)
	return ok
}
