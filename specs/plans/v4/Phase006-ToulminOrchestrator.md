# Phase006 — orchestrator 통합 + SSOT 카운트 조정

## 목표

orchestrator의 SSOT 탐지·파싱·검증 파이프라인을 Rego→toulmin 전환에 맞게 갱신한다. SSOT 10개 → 9개 (Rego 폐기).

## 변경 상세

### 1. SSOT 카운트 조정

| 파일 | 변경 |
|------|------|
| `internal/orchestrator/ssot_kind.go` | `KindPolicy` → `KindAuthz` 이름 변경 |
| `internal/orchestrator/all_ssot_kinds.go` | 필수 SSOT 10개 → 9개. 또는 KindAuthz를 유지하되 설명 갱신 |

### 2. SSOT 탐지

| 파일 | 변경 |
|------|------|
| `internal/orchestrator/detect_ssots.go` | `policy/` 디렉토리 + `.rego` 확장자 → `authz/` 디렉토리 + `.go` + `.yaml` 확장자 |

### 3. 파싱

| 파일 | 변경 |
|------|------|
| `internal/orchestrator/parsed.go` | `policy.ParseDir()` 호출은 동일 (Phase002에서 내부 변경됨). 경로만 `authz/` |

### 4. 개별 검증

| 파일 | 변경 |
|------|------|
| `internal/orchestrator/validate_policy.go` | → `validate_authz.go` 이름 변경. Rego 문법 검증 → Go 파싱 가능성(go/parser) + YAML 그래프 유효성(`toulmin graph --check` 상당) 검증 |

### 5. 코드 생성

| 파일 | 변경 |
|------|------|
| `internal/orchestrator/gen_authz.go` | `gogin.GenerateAuthzPackage()` 호출 (Phase005에서 내부 변경됨) |

### 6. 리포터 메시지

| 파일 | 변경 |
|------|------|
| `internal/reporter/` | "OPA Rego" → "Authz" 또는 "Toulmin" 표시 갱신 |

### 7. manual-for-ai.md

| 파일 | 변경 |
|------|------|
| `manual-for-ai.md` | OPA Rego 문법 섹션 → Go+toulmin 정책 작성법으로 교체 |

### 8. CLAUDE.md

| 파일 | 변경 |
|------|------|
| `CLAUDE.md` | "9개 SSOT" 반영. 디렉토리 구조에서 `policy/` → `authz/` |

## 의존성

- Phase001~005 전체 완료 필수

## 검증 방법

- `fullend validate specs/dummys/zenflow-try06/` 통과 (authz/ 디렉토리 탐지)
- `fullend gen specs/dummys/zenflow-try06/ artifacts/dummys/zenflow-try06/` 통과
- `fullend status specs/dummys/zenflow-try06/` 에서 "Authz" SSOT 표시
- `go test ./...` 전체 통과
