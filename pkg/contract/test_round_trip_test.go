//ff:func feature=contract type=rule control=sequence
//ff:what RoundTrip: 디렉티브 String→Parse 라운드트립이 동일 값을 유지하는지 테스트
package contract

import (
	"testing"
)

func TestRoundTrip(t *testing.T) {
	original := &Directive{Ownership: "preserve", SSOT: "db/gigs.sql", Contract: "e1d9f20"}
	s := original.String()
	parsed, err := Parse(s)
	if err != nil {
		t.Fatalf("Parse(String()) error: %v", err)
	}
	if parsed.Ownership != original.Ownership || parsed.SSOT != original.SSOT || parsed.Contract != original.Contract {
		t.Errorf("roundtrip mismatch: %+v != %+v", parsed, original)
	}
}
