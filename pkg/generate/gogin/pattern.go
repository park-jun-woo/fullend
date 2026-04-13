//ff:type feature=gen-gogin type=model topic=interface-derive
//ff:what Pattern — generate_method_from_iface dispatch 결과 열거

package gogin

// Pattern identifies which method body implementation to emit.
type Pattern int

const (
	// PatternSkip: WithTx special case — emitted separately.
	PatternSkip Pattern = iota
	// PatternCursorPagination: List* + pagination.Cursor[T] return.
	PatternCursorPagination
	// PatternOffsetPagination: List* + pagination.Page[T] or []T return.
	PatternOffsetPagination
	// PatternSliceReturn: non-List slice return ([]T).
	PatternSliceReturn
	// PatternFind: single-row lookup (Find*/Get* or seqType=="get").
	PatternFind
	// PatternCreate: seqType=="post" inline QueryRowContext + scan.
	PatternCreate
	// PatternUpdateDelete: seqType=="put"/"delete" inline ExecContext.
	PatternUpdateDelete
	// PatternFallbackOne: default with query.Cardinality=="one".
	PatternFallbackOne
	// PatternFallbackExec: default otherwise.
	PatternFallbackExec
)
