//ff:func feature=contract type=walker control=sequence
//ff:what 단일 FuncDecl에서 보존 함수를 추출한다

package contract

import (
	"go/ast"
	"go/token"
)

// extractPreservedFunc extracts a preserved function from a single FuncDecl, if applicable.
func extractPreservedFunc(fd *ast.FuncDecl, src string, fset *token.FileSet) (string, *PreservedFunc, bool) {
	if fd.Body == nil {
		return "", nil, false
	}

	d := extractDirectiveFromDoc(fd.Doc)
	if d == nil || d.Ownership != "preserve" {
		return "", nil, false
	}

	bodyStart := fset.Position(fd.Body.Lbrace).Offset
	bodyEnd := fset.Position(fd.Body.Rbrace).Offset
	bodyText := src[bodyStart+1 : bodyEnd]

	return fd.Name.Name, &PreservedFunc{
		Directive: *d,
		BodyText:  bodyText,
	}, true
}
