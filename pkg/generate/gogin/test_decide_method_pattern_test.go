//ff:func feature=gen-gogin type=test control=iteration dimension=1 topic=interface-derive
//ff:what DecideMethodPattern 의 9 Pattern 전수 테이블 테스트

package gogin

import "testing"

func TestDecideMethodPattern(t *testing.T) {
	cases := []struct {
		name  string
		facts MethodFacts
		want  Pattern
	}{
		{"WithTx short-circuits everything", MethodFacts{MethodName: "WithTx", SeqType: "post", IsListPrefix: true, IsCursorReturn: true}, PatternSkip},
		{"List + Cursor return", MethodFacts{MethodName: "ListUsers", SeqType: "get", IsListPrefix: true, IsCursorReturn: true}, PatternCursorPagination},
		{"List without Cursor", MethodFacts{MethodName: "ListUsers", SeqType: "get", IsListPrefix: true}, PatternOffsetPagination},
		{"Slice return (non-List)", MethodFacts{MethodName: "FindByTag", SeqType: "get", IsFindPrefix: true, IsSliceReturn: true}, PatternSliceReturn},
		{"Find prefix", MethodFacts{MethodName: "FindByID", SeqType: "", IsFindPrefix: true}, PatternFind},
		{"seqType get without Find prefix", MethodFacts{MethodName: "GetUser", SeqType: "get"}, PatternFind},
		{"POST inline", MethodFacts{MethodName: "CreateUser", SeqType: "post"}, PatternCreate},
		{"PUT inline", MethodFacts{MethodName: "UpdateUser", SeqType: "put"}, PatternUpdateDelete},
		{"DELETE inline", MethodFacts{MethodName: "DeleteUser", SeqType: "delete"}, PatternUpdateDelete},
		{"default + cardinality=one", MethodFacts{MethodName: "ExistsByEmail", SeqType: "", Cardinality: "one"}, PatternFallbackOne},
		{"default + cardinality=many", MethodFacts{MethodName: "Bulk", SeqType: "", Cardinality: "many"}, PatternFallbackExec},
		{"default + cardinality empty", MethodFacts{MethodName: "Sync", SeqType: "", Cardinality: ""}, PatternFallbackExec},
		{"WithTx beats List+Cursor", MethodFacts{MethodName: "WithTx", IsListPrefix: true, IsCursorReturn: true}, PatternSkip},
		{"List+Cursor beats Slice", MethodFacts{MethodName: "ListX", IsListPrefix: true, IsCursorReturn: true, IsSliceReturn: true}, PatternCursorPagination},
		{"List beats Slice", MethodFacts{MethodName: "ListX", IsListPrefix: true, IsSliceReturn: true}, PatternOffsetPagination},
		{"Slice beats Find when no Find prefix", MethodFacts{MethodName: "AllBy", IsSliceReturn: true}, PatternSliceReturn},
		{"Find beats POST", MethodFacts{MethodName: "FindX", IsFindPrefix: true, SeqType: "post"}, PatternFind},
		{"POST beats PUT fallthrough", MethodFacts{MethodName: "CreateX", SeqType: "post"}, PatternCreate},
	}
	for _, tc := range cases {
		if got := DecideMethodPattern(tc.facts); got != tc.want {
			t.Errorf("%s: want %d got %d (facts=%+v)", tc.name, tc.want, got, tc.facts)
		}
	}
}
