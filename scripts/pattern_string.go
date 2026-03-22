//ff:func feature=stat type=util control=sequence
//ff:what pattern을 "if=N loop=N switch=N" 형식 문자열로 변환
package main

import "fmt"

func (p pattern) String() string {
	return fmt.Sprintf("if=%d loop=%d switch=%d", p.Ifs, p.Loops, p.Switches)
}
