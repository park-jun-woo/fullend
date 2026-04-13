//ff:func feature=gen-gogin type=test-helper control=sequence topic=main-init
//ff:what DecideMainInit 백엔드 활성화·ContextImport 기대값 검증 헬퍼

package gogin

import "testing"

func assertInitNeedsBackends(t *testing.T, facts MainFacts, wantCtx, wantSes, wantCac, wantFil bool) {
	t.Helper()
	needs := DecideMainInit(facts)
	if needs.NeedsContextImport != wantCtx {
		t.Errorf("NeedsContextImport: want %v got %v", wantCtx, needs.NeedsContextImport)
	}
	if needs.Session.Enabled != wantSes {
		t.Errorf("Session.Enabled: want %v got %v", wantSes, needs.Session.Enabled)
	}
	if needs.Cache.Enabled != wantCac {
		t.Errorf("Cache.Enabled: want %v got %v", wantCac, needs.Cache.Enabled)
	}
	if needs.File.Enabled != wantFil {
		t.Errorf("File.Enabled: want %v got %v", wantFil, needs.File.Enabled)
	}
}
