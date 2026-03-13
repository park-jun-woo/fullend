# Phase 004: fullend contract + validate 통합 ✅ 완료

## 목표

`fullend contract` 커맨드로 SSOT ↔ artifacts 계약 상태를 한눈에 확인하고, `fullend validate`에 계약 검증을 통합한다.

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/contract/scan.go` | 신규 — artifacts 디렉토리 전체 스캔, 디렉티브 수집 |
| `internal/contract/verify.go` | 신규 — SSOT 존재 확인 + contract hash 재계산 비교 |
| `internal/orchestrator/contract.go` | 신규 — contract 커맨드 오케스트레이션 |
| `internal/reporter/contract.go` | 신규 — contract 상태 출력 포매터 |
| `internal/orchestrator/validate.go` | 수정 — validate 흐름에 contract 검증 추가 |
| `internal/reporter/validate.go` | 수정 — validate 출력에 Contract 행 추가 |
| `cmd/fullend/main.go` | 수정 — `contract` 서브커맨드 추가 |

## 상세 설계

### scan — artifacts 스캔

```go
package contract

type FuncStatus struct {
    File      string
    Function  string
    Directive Directive
    Status    string    // "gen", "preserve", "broken", "orphan"
    Detail    string    // 위반 상세 (예: "arg 'deadline' added")
}

// ScanDir은 artifacts 디렉토리의 모든 Go/TSX 파일에서 //fullend: 디렉티브를 추출한다.
func ScanDir(artifactsDir string) ([]FuncStatus, error)
```

### verify — 계약 검증

```go
// Verify는 각 FuncStatus의 ssot 경로가 존재하는지, contract hash가 현재 SSOT와 일치하는지 확인한다.
func Verify(specsDir string, funcs []FuncStatus) []FuncStatus
```

검증 규칙:

| 조건 | Status |
|---|---|
| SSOT 파일 존재 + hash 일치 | `gen` 또는 `preserve` (기존 유지) |
| SSOT 파일 존재 + hash 불일치 | `broken` |
| SSOT 파일 미존재 | `orphan` |

### contract 커맨드

```bash
fullend contract <specs-dir> <artifacts-dir>
```

출력:

```
Contract Status:
  gen       service/gig/list_gigs.go      ListGigs         fullend 소유
  preserve  service/gig/create_gig.go     CreateGig        계약 유지 ✓
  broken    service/gig/update_gig.go     UpdateGig        계약 위반 ✗ (arg added)
  orphan    service/gig/old_feature.go    OldFeature       SSOT 삭제됨 ⚠

  gen:      1 functions
  preserve: 1 functions
  broken:   1 functions
  orphan:   1 functions
```

### validate 통합

`fullend validate` 마지막에 artifacts 디렉토리가 존재하면 contract 검증도 수행:

```
✓ Config       my-project
✓ OpenAPI      7 endpoints
✓ Cross        0 mismatches
✗ Contract     1 violation, 1 orphan
```

artifacts 디렉토리가 없으면 Contract 행 스킵 (gen 전에는 검증할 대상 없음).

## 의존성

- Phase 001 (`internal/contract` — Directive, Hash)
- Phase 003 (preserve 함수가 존재하는 artifacts)

## 검증

```bash
go test ./internal/contract/...
go test ./...
```

1. `fullend contract` — gen/preserve/broken/orphan 각 상태 정확히 분류
2. `fullend validate` — Contract 행 출력, violation 시 ✗ 표시
3. artifacts 없을 때 — Contract 행 스킵
4. SSOT 삭제된 함수 — orphan 감지
