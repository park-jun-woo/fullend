//ff:type feature=gen-gogin type=model topic=interface-derive
//ff:what MethodFacts — DecideMethodPattern 입력 축 값

package gogin

// MethodFacts carries the axis values used by DecideMethodPattern.
// Fields are precomputed by NewMethodFacts to keep the decider side-effect-free.
type MethodFacts struct {
	MethodName     string
	SeqType        string // "get"/"post"/"put"/"delete"/""
	Cardinality    string // "one"/"many"/""
	HasQueryOpts   bool
	IsListPrefix   bool // isListMethod(Name) && HasQueryOpts
	IsFindPrefix   bool // Name startsWith "Find"
	IsCursorReturn bool // ReturnSig contains "pagination.Cursor["
	IsSliceReturn  bool // ReturnSig contains "[]"
}
