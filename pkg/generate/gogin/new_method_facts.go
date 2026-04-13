//ff:func feature=gen-gogin type=util control=sequence topic=interface-derive
//ff:what NewMethodFacts — ifaceMethod + sqlcQuery + seqType → MethodFacts 프로젝션

package gogin

import "strings"

// NewMethodFacts computes the axis values for a given method signature.
func NewMethodFacts(m ifaceMethod, query *sqlcQuery, seqType string) MethodFacts {
	hasOpts := hasQueryOptsParam(m)
	cardinality := ""
	if query != nil {
		cardinality = query.Cardinality
	}
	return MethodFacts{
		MethodName:     m.Name,
		SeqType:        seqType,
		Cardinality:    cardinality,
		HasQueryOpts:   hasOpts,
		IsListPrefix:   isListMethod(m.Name) && hasOpts,
		IsFindPrefix:   strings.HasPrefix(m.Name, "Find"),
		IsCursorReturn: strings.Contains(m.ReturnSig, "pagination.Cursor["),
		IsSliceReturn:  strings.Contains(m.ReturnSig, "[]"),
	}
}
