//ff:type feature=contract type=model
//ff:what 스플라이스 결과와 경고를 담는 구조체
package contract

// SpliceResult holds the merged content and any warnings.
type SpliceResult struct {
	Content  string
	Warnings []Warning
}
