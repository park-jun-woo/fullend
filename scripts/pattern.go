//ff:type feature=stat type=model
//ff:what 함수 depth-1 제어문(if/loop/switch) 개수 패턴
package main

type pattern struct {
	Ifs      int
	Loops    int // for + range
	Switches int
}
