//ff:type feature=gen-gogin type=model
//ff:what holds OpenAPI path parameter metadata for route generation

package gogin

type pathParamInfo struct {
	Name   string // original param name e.g. "CourseID"
	GoName string // PascalCase e.g. "CourseID"
	IsInt  bool
}
