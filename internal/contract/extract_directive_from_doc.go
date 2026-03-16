//ff:func feature=contract type=util control=iteration dimension=1
//ff:what AST Doc 코멘트 그룹에서 fullend 디렉티브를 찾는다
package contract

import "go/ast"

// extractDirectiveFromDoc finds a //fullend: directive in a doc comment group.
func extractDirectiveFromDoc(doc *ast.CommentGroup) *Directive {
	if doc == nil {
		return nil
	}
	for _, c := range doc.List {
		if d, err := Parse(c.Text); err == nil {
			return d
		}
	}
	return nil
}
