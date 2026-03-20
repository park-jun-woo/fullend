//ff:func feature=funcspec type=parser control=iteration dimension=1
//ff:what 단일 func spec .go 파일을 파싱하여 FuncSpec을 반환한다
package funcspec

import (
	"go/parser"
	"go/token"
	"strings"

	"github.com/park-jun-woo/fullend/pkg/diagnostic"
)

// ParseFile parses a single func spec .go file.
func ParseFile(path string) (*FuncSpec, []diagnostic.Diagnostic) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, []diagnostic.Diagnostic{{
			File:    path,
			Line:    0,
			Phase:   diagnostic.PhaseParse,
			Level:   diagnostic.LevelError,
			Message: "Go parse error: " + err.Error(),
		}}
	}

	spec := &FuncSpec{
		Package: f.Name.Name,
	}

	// Extract @func and @description from file-level comments.
	for _, cg := range f.Comments {
		parseCommentGroup(cg, spec)
	}

	if spec.Name == "" {
		return nil, nil // Not a func spec file.
	}

	// Extract imports.
	for _, imp := range f.Imports {
		p := strings.Trim(imp.Path.Value, `"`)
		spec.Imports = append(spec.Imports, p)
	}

	// Extract Request/Response structs and function body.
	expectedRequest := ucFirst(spec.Name) + "Request"
	expectedResponse := ucFirst(spec.Name) + "Response"

	for _, decl := range f.Decls {
		processDecl(decl, fset, spec, expectedRequest, expectedResponse)
	}

	return spec, nil
}
