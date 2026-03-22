//ff:type feature=stat type=model
//ff:what 함수의 패턴 분석 결과 (위치, 이름, 패턴, 라인 수)
package main

type funcInfo struct {
	File    string
	Name    string
	Pattern pattern
	Lines   int
}
