//ff:func feature=gen-gogin type=test control=sequence topic=main-init
//ff:what DecideMainInit — 빈 facts 는 auth/authz 불필요 검증

package gogin

import "testing"

func TestDecideMainInit_AuthEqualsAuthz(t *testing.T) {
	needs := DecideMainInit(MainFacts{})
	if needs.Auth || needs.Authz {
		t.Errorf("empty facts should not need auth/authz, got auth=%v authz=%v", needs.Auth, needs.Authz)
	}
}
