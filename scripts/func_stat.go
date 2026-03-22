//ff:type feature=stat type=model
//ff:what 함수별 순수 body 라인 수 통계 레코드
package main

type funcStat struct {
	File      string
	Name      string
	TotalBody int // 전체 body 라인
	PureBody  int // if/for 블록 제외 순수 라인
}
