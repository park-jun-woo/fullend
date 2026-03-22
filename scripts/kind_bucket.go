//ff:type feature=stat type=model
//ff:what 제어문 종류별 라인 수 분포 버킷
package main

type kindBucket struct {
	b1_5   int
	b6_10  int
	b11_20 int
	b21_50 int
	b51    int
	total  int
	max    int
}
