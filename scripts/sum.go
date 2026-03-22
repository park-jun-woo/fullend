//ff:func feature=stat type=util control=iteration dimension=1
//ff:what 정수 슬라이스의 합계 계산
package main

func sum(a []int) int {
	s := 0
	for _, v := range a {
		s += v
	}
	return s
}
