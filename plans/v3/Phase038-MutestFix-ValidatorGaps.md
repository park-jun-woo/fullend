# Phase038: Mutest FAIL 수정 — Validator 검출 갭 4건 ✅ 완료

## 목표

Mutation test 전수 실행 결과 **FAIL 4건** 처리. 실제 갭 3건 수정 + mutest 기대값 오류 1건 정정.

## 현황 (Mutest 전수 실행 결과)

| 구분 | 건수 |
|---|---|
| PASS | 65 (non-SSaC 52 + SSaC 샘플 13) |
| FAIL | 4 |
| SKIP | 11 (dummy fixture 부재) |

### FAIL 4건 근본 원인 (코드 재현 확인 완료)

| ID | 유형 | 근본 원인 | 조치 |
|---|---|---|---|
| MUT-SCENARIO-001 | 실제 갭 | `validate_scenario.go` line 16: `.feature` 검사가 `scenario/` 디렉토리만 검색. `tests/` 미검사 | 코드 수정 |
| MUT-STML-002 | 실제 갭 | fetch 블록 내부의 `data-action`이 `page.Actions`에 수집되지 않음 → validator 미검증. 검증 로직(line 150) 자체는 정상이나 fetch 내 action이 전달 안 됨 | 코드 수정 |
| MUT-SSAC-OPENAPI-002 | mutest 오류 | Login `@response token` — shorthand 경로. 변수 단위 반환이므로 개별 property 이름 비교 불가. 설계상 정상 동작 | mutest 결과 정정 |
| MUT-SSAC-014 | 실제 갭 | `validate_model.go`는 `seq.Model`만 검증. Result 타입 대소문자(`gig` vs `Gig`) 미검증. DDL crosscheck도 소문자 `gig` → `gigs` 매칭 성공 | 코드 수정 |

## 설계

### 1. MUT-SCENARIO-001: `.feature` 검사 경로 확대

**수정**: `validate_scenario.go` line 16 — `tests/` 디렉토리도 `.feature` 검사 대상에 추가.

```go
// 현재: scenarioDir := filepath.Join(specsRoot, "scenario")
// 수정: testsDir도 검사
```

**변경 파일**: `internal/orchestrator/validate_scenario.go`

### 2. MUT-STML-002: fetch 내부 action 수집 누락

**수정**: STML validator가 fetch 블록의 children을 순회하여 내부 action도 검증하도록 수정. `validateFetchBlock` 내에서 `fb.Children` 중 kind="action"인 것을 `validateActionBlock`으로 위임.

**변경 파일**: `internal/stml/validator/validator.go`

### 3. MUT-SSAC-OPENAPI-002: mutest 결과 정정

**수정**: mutest 기대값이 틀림 — shorthand `@response token`은 변수 단위 반환이므로 property-level 비교 불가. mutest 결과를 SKIP("shorthand 경로에서 property 이름 비교 불가")으로 정정.

**변경 파일**: `files/mutests/ssac-openapi.md`

### 4. MUT-SSAC-014: Result 타입 PascalCase 강제

**수정**: SSaC validator에서 @get/@post의 Result 타입이 대문자로 시작하는지 검증. 소문자 시작이면 ERROR.

**변경 파일**: `internal/ssac/validator/validate_model.go` 또는 새 파일

## 변경 파일

- `internal/orchestrator/validate_scenario.go` — `.feature` 검사에 `tests/` 추가
- `internal/stml/validator/validator.go` — fetch 내부 action 검증 추가
- `internal/ssac/validator/validate_model.go` — Result 타입 PascalCase 검증
- `files/mutests/ssac-openapi.md` — MUT-SSAC-OPENAPI-002 결과 SKIP 정정

## 검증

1. `go test ./...`
2. Mutest 4건 재실행 → 3건 PASS + 1건 SKIP
3. `fullend validate specs/dummys/gigbridge-try02/` — 정상 통과
4. `fullend validate specs/dummys/zenflow-try05/` — 정상 통과
