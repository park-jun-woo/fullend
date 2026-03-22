//ff:func feature=stat type=util control=selection
//ff:what kindBucket에 라인 수를 분류하여 추가
package main

func (b *kindBucket) add(lines int) {
	b.total++
	if lines > b.max {
		b.max = lines
	}
	switch {
	case lines <= 5:
		b.b1_5++
	case lines <= 10:
		b.b6_10++
	case lines <= 20:
		b.b11_20++
	case lines <= 50:
		b.b21_50++
	default:
		b.b51++
	}
}
