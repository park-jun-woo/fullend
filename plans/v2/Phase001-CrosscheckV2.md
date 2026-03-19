✅ 완료

# Phase 001: crosscheck SSaC v2 마이그레이션

## 목표

SSaC v2 파서 타입 변경에 맞춰 `internal/crosscheck/` 전체 모듈을 업데이트한다.

## 배경

SSaC가 v2로 대격변했다. 멀티라인 `@sequence` 블록이 원라인 DSL(`@get`, `@post`, `@empty`, `@call` 등)로 바뀌면서 파서 IR 타입이 변경되었다.

### 주요 타입 변경

| 항목 | v1 | v2 |
|------|----|----|
| Sequence.Type | `"get"`, `"authorize"`, `"call"` 등 | `"get"`, `"auth"`, `"call"` 등 (동일 상수) |
| 인자 | `Sequence.Params []Param` (Name, Source, Column) | `Sequence.Args []Arg` (Source, Field, Literal) |
| 가드 | `"guard nil"`, `"guard exists"`, `"guard state"` | `"empty"`, `"exists"`, `"state"` |
| @auth | `"authorize"` + @action/@resource/@id | `"auth"` + Action, Resource, Inputs |
| @response | `"response json"` + @var 목록 | `"response"` + Fields map |
| @call | @func pkg.funcName + @param | Model = "pkg.Func", Args |

### 영향받는 파일

| 파일 | LOC | 변경 수준 |
|------|-----|----------|
| `crosscheck/ssac_ddl.go` | 138 | 중 — Param → Arg 구조 |
| `crosscheck/ssac_openapi.go` | 45 | 소 — fn.Name 접근만 (변경 없을 수 있음) |
| `crosscheck/ddl_coverage.go` | 177 | 소 — seq.Model, seq.Result 접근 |
| `crosscheck/openapi_ddl.go` | 332 | 소 — seq.Model 접근 |
| `crosscheck/states.go` | 226 | 중 — guard state → @state, Params → Inputs |
| `crosscheck/policy.go` | 171 | 중 — authorize → auth, @id → Inputs |
| `crosscheck/func.go` | 389 | 대 — @func → @call, Params → Args 전면 |
| `crosscheck/func_test.go` | 336 | 대 — 테스트 데이터 전면 재작성 |
| `crosscheck/crosscheck.go` | 79 | 소 — 타입 참조만 |
| `crosscheck/archived.go` | 112 | 없음 — DDL 전용 |

## 변경 항목

### A. ssac_ddl.go — SSaC ↔ DDL 검증

- `seq.Params` → `seq.Args` 변환
- `p.Source`, `p.Name`, `p.Column` → `arg.Source`, `arg.Field`, `arg.Literal`

### B. states.go — States ↔ SSaC 검증

- `seq.Type == "guard state"` → `seq.Type == "state"`
- `seq.Target` (stateDiagram ID) → `seq.DiagramID`
- `seq.Params` 기반 필드 추출 → `seq.Inputs` 맵 사용

### C. policy.go — Policy ↔ SSaC 검증

- `seq.Type == "authorize"` → `seq.Type == "auth"`
- `seq.Action`, `seq.Resource` 필드는 유지
- `@id` 처리 → `seq.Inputs` 맵에서 추출

### D. func.go — Func ↔ SSaC 검증 (최대 규모)

- `seq.Type == "call"` 유지
- `seq.Func` → `seq.Model` (v2에서 @call도 Model 필드 사용)
- `seq.Package` → Model에서 `.` 앞부분 추출
- `seq.Params` → `seq.Args` 전면 교체
- `p.Source`, `p.Name` → `arg.Source`, `arg.Field`

### E. func_test.go — 테스트 데이터 재작성

- 모든 `ssacparser.Param{}` → `ssacparser.Arg{}`
- 시퀀스 타입 상수 변경 반영

### F. ddl_coverage.go, openapi_ddl.go — 경미한 수정

- `seq.Result` 구조 동일 (Type, Var) — 변경 최소

## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| `internal/crosscheck/ssac_ddl.go` | Param → Arg 구조 전환 |
| `internal/crosscheck/states.go` | guard state → state, Params → Inputs |
| `internal/crosscheck/policy.go` | authorize → auth, @id → Inputs |
| `internal/crosscheck/func.go` | @func → @call Model, Params → Args |
| `internal/crosscheck/func_test.go` | 테스트 데이터 전면 재작성 |
| `internal/crosscheck/ddl_coverage.go` | Result 접근 확인 |
| `internal/crosscheck/openapi_ddl.go` | Model 접근 확인 |
| `internal/crosscheck/ssac_openapi.go` | 변경 여부 확인 |
| `internal/crosscheck/crosscheck.go` | 타입 임포트 확인 |

## 의존성

- SSaC v2 파서 (`go.mod` replace 디렉티브로 로컬 참조)
- `go mod tidy` 후 빌드 확인

## 검증 방법

```bash
# 1. 빌드
go build ./internal/crosscheck/...

# 2. 단위 테스트
go test ./internal/crosscheck/... -count=1

# 3. fullend 전체 빌드
go build ./cmd/fullend/
```
