//ff:func feature=contract type=util control=sequence
//ff:what 단일 함수 선언에 대해 보존 본문 교체 항목과 경고를 생성한다
package contract

import (
	"go/ast"
	"go/token"
)

// spliceFunc processes a single FuncDecl for preserved body restoration.
// Returns replacement entries and an optional warning.
func spliceFunc(fd *ast.FuncDecl, fset *token.FileSet, pf *PreservedFunc, filePath string) ([]spliceReplacement, *Warning) {
	var replacements []spliceReplacement
	var warn *Warning

	// Check contract change.
	newD := extractDirectiveFromDoc(fd.Doc)
	if newD != nil && newD.Contract != pf.Directive.Contract {
		warn = &Warning{
			File:        filePath,
			Function:    fd.Name.Name,
			OldContract: pf.Directive.Contract,
			NewContract: newD.Contract,
		}
	}

	// Replace body.
	bodyStart := fset.Position(fd.Body.Lbrace).Offset
	bodyEnd := fset.Position(fd.Body.Rbrace).Offset
	replacements = append(replacements, spliceReplacement{
		start: bodyStart + 1,
		end:   bodyEnd,
		text:  pf.BodyText,
	})

	// Update directive from gen to preserve.
	if newD != nil {
		r := buildDirectiveReplacement(fd.Doc, fset, newD, pf.Directive.Contract)
		if r != nil {
			replacements = append(replacements, *r)
		}
	}

	return replacements, warn
}
