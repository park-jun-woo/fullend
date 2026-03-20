//ff:func feature=contract type=rule control=sequence
//ff:what DirectiveString: String() 메서드가 올바른 디렉티브 문자열을 생성하는지 테스트
package contract

import (
	"testing"
)

func TestDirectiveString(t *testing.T) {
	d := &Directive{Ownership: "gen", SSOT: "service/gig/create_gig.ssac", Contract: "a3f8c10"}
	got := d.String()
	want := "//fullend:gen ssot=service/gig/create_gig.ssac contract=a3f8c10"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
