//ff:type feature=gen-gogin type=model
//ff:what parsed interface method from models_gen.go

package gogin

// ifaceMethod represents a parsed interface method from models_gen.go.
type ifaceMethod struct {
	Name       string
	ParamSig   string // e.g. "courseID int64, opts QueryOpts"
	ReturnSig  string // e.g. "(*Course, error)"
	Params     []ifaceParam
}
