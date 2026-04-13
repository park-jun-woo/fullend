//ff:func feature=gen-gogin type=test control=iteration dimension=1 topic=main-init
//ff:what DecideMainInit 6축 독립 + ContextImport 파생 조합 테이블 테스트

package gogin

import (
	"testing"

	"github.com/park-jun-woo/fullend/pkg/parser/manifest"
)

func TestDecideMainInit_Backends(t *testing.T) {
	localFile := &manifest.FileBackend{Backend: "local"}
	s3File := &manifest.FileBackend{Backend: "s3"}
	cases := []struct {
		name    string
		facts   MainFacts
		wantCtx bool
		wantSes bool
		wantCac bool
		wantFil bool
	}{
		{"all empty", MainFacts{}, false, false, false, false},
		{"session postgres → context", MainFacts{SessionBackend: "postgres"}, true, true, false, false},
		{"session memory → no context", MainFacts{SessionBackend: "memory"}, false, true, false, false},
		{"cache postgres → context", MainFacts{CacheBackend: "postgres"}, true, false, true, false},
		{"cache memory → no context", MainFacts{CacheBackend: "memory"}, false, false, true, false},
		{"file local → no context", MainFacts{FileConfig: localFile}, false, false, false, true},
		{"file s3 → context", MainFacts{FileConfig: s3File}, true, false, false, true},
		{"all three + postgres → context", MainFacts{SessionBackend: "postgres", CacheBackend: "memory", FileConfig: localFile}, true, true, true, true},
		{"unknown session backend ignored", MainFacts{SessionBackend: "mongo"}, false, false, false, false},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assertInitNeedsBackends(t, tc.facts, tc.wantCtx, tc.wantSes, tc.wantCac, tc.wantFil)
		})
	}
}
