//ff:type feature=stat type=model
//ff:what if/for/switch 블록 내부 라인 수 통계 레코드
package main

type blockStat struct {
	File  string
	Func  string
	Kind  string // "if", "for", "range", "switch"
	Lines int    // 블록 내부 라인 수
}
