//ff:func feature=gen-gogin type=util
//ff:what converts OpenAPI path params to Go 1.22 mux style (pass-through)

package gogin

// convertPathParams converts OpenAPI path params like {CourseID} to Go 1.22 mux style {CourseID}.
// Go 1.22 mux uses the same brace syntax, so this is mostly a pass-through.
func convertPathParams(path string) string {
	return path
}
