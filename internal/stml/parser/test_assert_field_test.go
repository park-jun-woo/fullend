//ff:func feature=stml-parse type=test-helper control=sequence
//ff:what FieldBind 속성을 검증하는 테스트 헬퍼
package parser

import "testing"

func assertField(t *testing.T, f FieldBind, name, tag, typ string) {
	t.Helper()
	if f.Name != name { t.Errorf("Field.Name = %q, want %q", f.Name, name) }
	if f.Tag != tag { t.Errorf("Field.Tag = %q, want %q", f.Tag, tag) }
	if f.Type != typ { t.Errorf("Field.Type = %q, want %q", f.Type, typ) }
}
