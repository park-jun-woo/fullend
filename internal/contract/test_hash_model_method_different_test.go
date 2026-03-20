//ff:func feature=contract type=rule control=sequence topic=go-interface
//ff:what HashModelMethodDifferent: 다른 모델 메서드에 대해 다른 해시를 생성하는지 테스트
package contract

import (
	"testing"
)

func TestHashModelMethod_Different(t *testing.T) {
	h1 := HashModelMethod("Create", []string{"*Gig"}, []string{"*Gig", "error"})
	h2 := HashModelMethod("FindByID", []string{"int64"}, []string{"*Gig", "error"})
	if h1 == h2 {
		t.Errorf("different methods produced same hash: %s", h1)
	}
}
