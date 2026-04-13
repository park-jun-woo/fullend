//ff:func feature=gen-gogin type=decider control=selection topic=interface-derive
//ff:what DecideMethodPattern — MethodFacts → Pattern (depth 2 switch)

package gogin

// DecideMethodPattern selects the implementation Pattern.
// Depth 2: outer switch (depth 1) + default-branch if (depth 2).
func DecideMethodPattern(f MethodFacts) Pattern {
	switch {
	case f.MethodName == "WithTx":
		return PatternSkip
	case f.IsListPrefix && f.IsCursorReturn:
		return PatternCursorPagination
	case f.IsListPrefix:
		return PatternOffsetPagination
	case f.IsSliceReturn:
		return PatternSliceReturn
	case f.IsFindPrefix || f.SeqType == "get":
		return PatternFind
	case f.SeqType == "post":
		return PatternCreate
	case f.SeqType == "put" || f.SeqType == "delete":
		return PatternUpdateDelete
	default:
		if f.Cardinality == "one" {
			return PatternFallbackOne
		}
		return PatternFallbackExec
	}
}
