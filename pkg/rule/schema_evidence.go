//ff:type feature=rule type=model
//ff:what SchemaEvidence — 스키마 불일치 시 누락 필드 포함 결과
package rule

// SchemaEvidence extends Evidence with missing field details.
type SchemaEvidence struct {
	Rule    string
	Level   string
	Missing []string
	Message string
}
