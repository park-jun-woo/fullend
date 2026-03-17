# Phase048: 잔여 FAIL 2건 수정 — IANA HTTP status + input key case ✅ 완료

## 목표

mutest-report03 불확실 FAIL 2건을 해소하여 SKIP 제외 FAIL 0건을 달성한다.

## 변경 파일 목록

### 1. MUT-SSAC-020 — HTTP status 999 범위 검증 없음

- **chain**: `validateFunc` → `validateErrStatus`
- **현상**: SSaC 파서(`parse_guard.go:21`)가 `code >= 100 && code <= 599` 범위 체크로 999를 거부. ErrStatus=0으로 파싱 → validator에 도달하지 않음
- **수정**:
  - SSaC 파서 4개 파일(`parse_guard.go`, `parse_state.go`, `parse_auth.go`, `parse_call.go`)의 범위 제한 `code >= 100 && code <= 599` → `code > 0`으로 완화. 양의 정수면 파싱 허용
  - `valid_http_status.go` — IANA 등록 HTTP status code 테이블 (1xx~5xx)
  - `validate_err_status.go` — ErrStatus가 IANA 미등록이면 WARNING
  - `err_ctx.go` — `warn()` 헬퍼 추가
  - `validate_func_internal.go` — `validateErrStatus` 호출 추가

### 2. MUT-SSAC-005 — input key "Id" vs sqlc 파라미터 "ID" 대소문자 미검출

- **chain**: `CheckInputKeyCase` (신규) ← `rules.go`
- **현상**: `checkParamColumn`이 input key를 `pascalToSnake`로 변환하여 DDL 컬럼과 비교. `Id` → `id`, `ID` → `id` 모두 동일 컬럼에 매핑. 하지만 코드젠 시 Go struct 필드명이 되므로 `Id` ≠ `ID` (Go initialism)
- **수정**: `check_input_key_case.go` — SSaC input key를 sqlc `MethodInfo.Params`와 exact match 비교. case-insensitive 매칭은 되지만 exact match 실패 시 WARNING
- **파일**: `internal/crosscheck/check_input_key_case.go` (신규), `rules.go`에 Rule 추가

## 검증 결과

- `go test ./...` 전체 통과
- MUT-SSAC-020: `@empty wf "Workflow not found" 999` → "[WARN] HTTP status 999는 IANA 등록 코드가 아닙니다" 검출 ✓
- MUT-SSAC-005: `{Id: request.id}` → "[WARN] input key "Id"와 sqlc 파라미터 "ID" — 대소문자 불일치" 검출 ✓
- mutest-report03: FAIL 0건, 통과율 97.8%
