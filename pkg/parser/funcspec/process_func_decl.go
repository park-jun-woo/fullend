//ff:func feature=funcspec type=parser control=sequence
//ff:what processFuncDecl — ast.FuncDecl 에서 HasBody + ResponsePointer 추출

package funcspec

import (
	"go/ast"
	"go/token"
)

func processFuncDecl(d *ast.FuncDecl, fset *token.FileSet, spec *FuncSpec) {
	if d.Name.Name != ucFirst(spec.Name) {
		return
	}
	if d.Body != nil {
		spec.HasBody = !isStubBody(fset, d.Body)
	}
	spec.ResponsePointer = firstResultIsPointer(d.Type)
}
