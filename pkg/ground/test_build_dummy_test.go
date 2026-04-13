//ff:func feature=rule type=loader control=sequence
//ff:what Build 통합 검증 — dummy gigbridge 로드 후 신 필드가 채워지는지 확인
package ground

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/park-jun-woo/fullend/pkg/fullend"
)

// TestBuildDummyGigbridge loads the gigbridge dummy project via ParseAll
// and verifies that the new structural fields (Models/Tables/Ops/ReqSchemas)
// are populated.
func TestBuildDummyGigbridge(t *testing.T) {
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(thisFile), "..", "..")
	specsDir := filepath.Join(repoRoot, "dummys", "gigbridge", "specs")

	detected, _ := fullend.DetectSSOTs(specsDir)
	fs := fullend.ParseAll(specsDir, detected, nil)
	g := Build(fs)

	if len(g.Tables) == 0 {
		t.Errorf("Tables empty — DDL not loaded from %s", specsDir)
	}
	if len(g.Ops) == 0 {
		t.Errorf("Ops empty — OpenAPI not loaded")
	}
	// Models 는 model/*.go 또는 db/queries/*.sql 둘 중 하나라도 있으면 채워짐.
	// gigbridge 의 model.go 는 빈 파일이고 queries/ 가 있으면 sqlc 에서 채워짐.
	// 둘 다 없으면 0 일 수 있음 — 경고만.
	if len(g.Models) == 0 {
		t.Logf("Models empty — dummy 에 model interfaces / sqlc queries 없음 (경고)")
	}
	// ReqSchemas 는 requestBody 있는 operation 이 있으면 채워짐.
	if len(g.ReqSchemas) == 0 {
		t.Logf("ReqSchemas empty — 가능 (요청 body 없는 operation 만 있으면)")
	}

	// legacy bridge 확인
	if _, ok := g.Lookup["SymbolTable.model"]; !ok {
		t.Error("legacy Lookup['SymbolTable.model'] missing — populateModelLookup 미동작")
	}
}
