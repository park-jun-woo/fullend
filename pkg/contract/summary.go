//ff:func feature=contract type=util control=iteration dimension=1
//ff:what FuncStatus 목록에서 상태별 개수를 집계한다
package contract

// Summary returns counts by status.
func Summary(funcs []FuncStatus) (gen, preserve, broken, orphan int) {
	for _, f := range funcs {
		gen, preserve, broken, orphan = addStatusCount(f.Status, gen, preserve, broken, orphan)
	}
	return
}
