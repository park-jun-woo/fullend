//ff:func feature=rule type=generator control=sequence
//ff:what buildSeqGraph — 시퀀스 타입 분류 + defeat 관계 Graph 구성
package backend

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// buildSeqGraph constructs the Toulmin graph for sequence code generation.
// Each warrant identifies a sequence type or pattern, and defeat edges
// encode priority: cursor > offset > fk > simple for @get variants.
func buildSeqGraph() *toulmin.Graph {
	g := toulmin.NewGraph("seq-codegen")

	// Sequence type warrants
	g.Rule(IsGet)
	g.Rule(IsPost)
	g.Rule(IsPut)
	g.Rule(IsDelete)
	g.Rule(IsEmpty)
	g.Rule(IsExists)
	g.Rule(IsState)
	g.Rule(IsAuth)
	g.Rule(IsCall)
	g.Rule(IsPublish)
	g.Rule(IsResponse)

	// Get sub-pattern warrants (with defeat priority)
	simple := g.Rule(IsGet)
	fk := g.Rule(HasFKRef)
	offset := g.Rule(HasPaginateOffset)
	cursor := g.Rule(HasPaginateCursor)
	slice := g.Rule(HasSliceResult)

	fk.Attacks(simple)
	offset.Attacks(simple)
	offset.Attacks(fk)
	cursor.Attacks(offset)
	slice.Attacks(simple)

	return g
}
