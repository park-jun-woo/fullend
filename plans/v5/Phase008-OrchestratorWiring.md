# Phase008 — orchestrator 배선 교체 (internal/gen → pkg/generate)

## 목표

Phase006/007 에서 pkg/generate 가 단위 호출 가능해졌다. 이제 **orchestrator 가 pkg/generate 를 호출** 하도록 배선 교체.

1. **ParsedSSOTs ↔ Fullstack 어댑터** — orchestrator 가 parse 시 둘 다 구축 또는 어댑터 함수.
2. **gen_glue.go 교체** — `gen.Generate(parsed, cfg, stmlOut)` → `generate.Generate(fs, ground, cfg, stmlOut)`.
3. **gen_ssac.go / gen_stml.go** — 필요 시 pkg/generate/gogin/ssac, pkg/generate/react/stml 쪽으로 전환 또는 무변경 (orchestrator 가 이걸 직접 호출하지 않을 수도).
4. **grep 검증** — orchestrator 에서 `internal/gen`, `internal/genapi`, `internal/ssac/generator`, `internal/stml/generator` import 0건.

성공 기준:
- `go build ./pkg/... ./internal/... ./cmd/...` 통과.
- `go run ./cmd/fullend gen dummys/gigbridge/specs /tmp/out` 실행 가능 (품질 Tier 3 허용).
- `fullend validate` 는 종전 동작 유지.
- orchestrator import grep 결과 internal/gen 계열 0건.

---

## 전제

- **Phase006 완료** — gogin.Generate 활성.
- **Phase007 완료** — react/hurl.Generate 활성.
- pkg/generate.Generate 호출 가능 상태.

---

## 범위

### 포함

- `internal/orchestrator/gen_glue.go` — pkg/generate 호출로 교체
- `internal/orchestrator/parsed.go` — Fullstack + Ground 생성 로직 추가 (또는 어댑터 함수)
- `internal/orchestrator/gen_ssac.go` — pkg/generate/gogin/ssac 호출로 전환 (현재 internal/ssac/generator 호출 중일 가능성)
- `internal/orchestrator/gen_stml.go` — pkg/generate/react/stml 쪽 호출로 전환

### 포함하지 않음

- validate/crosscheck 는 이미 Ground 사용 중이라 무변경
- internal/gen 삭제 — 별도 Phase
- 구조 정리·Toulmin — Phase009
- Dummy 검증 — Phase010

---

## 작업 순서

### Step 1. 어댑터 함수 작성

Parse 결과에서 Fullstack + Ground 생성:

```go
// internal/orchestrator/build_pkg_context.go (신설)
import (
    "github.com/park-jun-woo/fullend/pkg/fullend"
    "github.com/park-jun-woo/fullend/pkg/ground"
    "github.com/park-jun-woo/fullend/pkg/rule"
)

func buildPkgContext(specsDir string) (*fullend.Fullstack, *rule.Ground, error) {
    detected, _ := fullend.DetectSSOTs(specsDir)
    fs := fullend.ParseAll(specsDir, detected, nil)
    g := ground.Build(fs)
    return fs, g, nil
}
```

기존 `ParseAll()` (ParsedSSOTs 반환) 는 그대로 유지 — validate/status/chain 이 아직 쓰므로.

### Step 2. gen_glue.go 교체

```go
// internal/orchestrator/gen_glue.go
import (
    ...
    pkggen "github.com/park-jun-woo/fullend/pkg/generate"
    reactgen "github.com/park-jun-woo/fullend/pkg/generate/react"
)

func genGlue(specsDir, artifactsDir string, ... stmlDeps, stmlPages, stmlPageOps) (reporter.StepResult, bool) {
    fs, ground, err := buildPkgContext(specsDir)
    if err != nil { ... }
    cfg := &pkggen.Config{ArtifactsDir: artifactsDir, SpecsDir: specsDir, ModulePath: ...}
    stmlOut := &reactgen.STMLGenOutput{Deps: stmlDeps, Pages: stmlPages, PageOps: stmlPageOps}
    if err := pkggen.Generate(fs, ground, cfg, stmlOut); err != nil { ... }
    return success
}
```

### Step 3. gen_ssac.go / gen_stml.go 전환

orchestrator 가 SSaC/STML generator 를 직접 호출하는지 확인:
- 만약 호출한다면 pkg 버전으로 import 교체
- 만약 gen_glue 가 통합 호출한다면 별도 수정 불필요

### Step 4. grep 검증

```bash
grep -rn "internal/gen\"\|internal/genapi\|internal/ssac/generator\|internal/stml/generator" internal/orchestrator/
```

결과 0건이어야 함.

### Step 5. 동작 검증

```bash
go build ./pkg/... ./internal/... ./cmd/...
go vet
go test ./pkg/...

# 실사용 테스트
go run ./cmd/fullend validate dummys/gigbridge/specs
go run ./cmd/fullend gen dummys/gigbridge/specs /tmp/gigbridge-pkg
ls /tmp/gigbridge-pkg/backend /tmp/gigbridge-pkg/tests
```

---

## 주의사항

### R1. Fullstack + ParsedSSOTs 병행

orchestrator 의 validate/status/chain 은 ParsedSSOTs 계속 사용. gen 만 Fullstack 으로 전환. 이중 파싱은 아까우나 Phase010 이후 정리.

### R2. internal/gen 은 삭제 금지

dead code 로 남김. 실제 삭제는 별도 Phase (Phase00N — 로드맵 끝).

### R3. 부분 실패 허용

`fullend gen` 실행 시 일부 생성 에러가 나도 OK (품질은 Tier 2/3 허용, Phase010 에서 검증).

---

## 완료 조건 (Definition of Done)

- [ ] `buildPkgContext` 어댑터 함수 작성
- [ ] gen_glue.go 가 pkg/generate 호출
- [ ] orchestrator 의 internal/gen 계열 import 0건
- [ ] `go build ./pkg/... ./internal/... ./cmd/...` 통과
- [ ] `go vet` 통과
- [ ] `fullend validate` 동작 유지
- [ ] `fullend gen dummys/gigbridge/specs /tmp/x` 실행 성공 (품질 불문)
- [ ] 커밋: `feat(orchestrator): gen 배선을 pkg/generate 로 교체`

---

## 다음 Phase

- **Phase009** — 구조 정리 + Toulmin 포인트 3군데.
- **Phase010** — Dummy 실용 검증 + 지표.
