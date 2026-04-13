//ff:type feature=contract type=model
//ff:what 바이트 오프셋 기반 텍스트 교체 항목을 나타내는 구조체
package contract

// spliceReplacement describes a byte-range text replacement.
type spliceReplacement struct {
	start int
	end   int
	text  string
}
