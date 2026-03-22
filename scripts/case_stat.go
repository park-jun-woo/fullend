//ff:type feature=stat type=model
//ff:what switch/select의 case 절 내부 라인 수 통계 레코드
package main

type caseStat struct {
	File  string
	Func  string
	Lines int
}
