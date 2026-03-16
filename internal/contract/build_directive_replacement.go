//ff:func feature=contract type=util control=iteration dimension=1
//ff:what Doc 코멘트에서 디렉티브를 찾아 preserve 교체 항목을 생성한다
package contract

import (
	"go/ast"
	"go/token"
)

// buildDirectiveReplacement finds the directive comment in doc and returns a replacement
// that changes it to preserve with the old contract hash.
func buildDirectiveReplacement(doc *ast.CommentGroup, fset *token.FileSet, newD *Directive, oldContract string) *spliceReplacement {
	for _, c := range doc.List {
		if _, err := Parse(c.Text); err == nil {
			lineStart := fset.Position(c.Pos()).Offset
			lineEnd := lineStart + len(c.Text)
			preserveD := &Directive{
				Ownership: "preserve",
				SSOT:      newD.SSOT,
				Contract:  oldContract,
			}
			return &spliceReplacement{
				start: lineStart,
				end:   lineEnd,
				text:  preserveD.String(),
			}
		}
	}
	return nil
}
