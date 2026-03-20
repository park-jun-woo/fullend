//ff:func feature=contract type=rule control=sequence topic=go-interface
//ff:what HashClaims: claims 맵에 대해 결정적 해시를 생성하는지 테스트
package contract

import (
	"testing"

	"github.com/park-jun-woo/fullend/internal/projectconfig"
)

func TestHashClaims(t *testing.T) {
	claims := map[string]projectconfig.ClaimDef{
		"ID":    {Key: "user_id", GoType: "int64"},
		"Email": {Key: "email", GoType: "string"},
		"Role":  {Key: "role", GoType: "string"},
	}
	h1 := HashClaims(claims)
	h2 := HashClaims(claims)
	if h1 != h2 {
		t.Errorf("same input produced different hashes: %s vs %s", h1, h2)
	}
	if len(h1) != 7 {
		t.Errorf("hash length = %d, want 7", len(h1))
	}
}
