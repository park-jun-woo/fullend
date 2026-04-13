//ff:func feature=contract type=rule control=sequence topic=go-interface
//ff:what HashModelMethodDeterministic: 동일 모델 메서드에 대해 동일 해시를 생성하는지 테스트
package contract

import (
	"testing"
)

func TestHashModelMethod_Deterministic(t *testing.T) {
	h1 := HashModelMethod("Create", []string{"*Gig"}, []string{"*Gig", "error"})
	h2 := HashModelMethod("Create", []string{"*Gig"}, []string{"*Gig", "error"})
	if h1 != h2 {
		t.Errorf("same input produced different hashes: %s vs %s", h1, h2)
	}
}
