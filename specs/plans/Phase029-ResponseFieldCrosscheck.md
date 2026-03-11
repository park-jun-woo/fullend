# ✅ 완료 Phase 029: SSaC @response ↔ OpenAPI response 필드 교차 검증

## 목표

DDL → OpenAPI 직접 비교를 삭제하고, SSaC @response 필드 ↔ OpenAPI response schema properties 검증으로 교체한다.

데이터 흐름은 DDL → SSaC → OpenAPI이므로, DDL과 OpenAPI를 직접 비교하는 것은 `password_hash` 같은 내부 컬럼까지 OpenAPI에 선언하라는 잘못된 압박이 된다. SSaC @response가 자연스러운 필터 역할을 한다.

## 변경 내용

### 1. `internal/crosscheck/ddl_coverage.go` — Rule 2 삭제

- DDL 컬럼 → OpenAPI 스키마 직접 비교 (56-88행) 삭제
- `buildSchemaProps()` 함수 삭제 (이 함수의 유일한 소비자가 Rule 2)
- `openapi3` import 삭제
- `CheckDDLCoverage` 시그니처에서 `doc *openapi3.T` 파라미터 제거
- Rule 1 (DDL 테이블 → SSaC 참조)은 유지

### 2. `internal/crosscheck/ssac_openapi.go` — @response 필드 검증 추가

기존 operationId 매칭(Rule 3, 4) 유지. 새로운 Rule 5 추가:

**Rule 5: SSaC @response 필드 → OpenAPI response schema property (ERROR)**

- 각 SSaC 함수에서 `@response { field: var }` 형태의 시퀀스를 찾는다
- `Fields` map의 키(JSON 필드명)를 추출한다
- `st.Operations[funcName].ResponseFields`와 대조한다
- @response 필드가 OpenAPI response에 없으면 → ERROR
- `@response varName` (shorthand) 형태는 개별 필드가 아니므로 스킵

**Rule 6: OpenAPI response property → SSaC @response 필드 (WARNING)**

- OpenAPI response에 있는데 SSaC @response에 없는 필드 → WARNING
- `@response varName` (shorthand) 함수는 스킵

### 3. `internal/crosscheck/crosscheck.go` — 호출부 수정

- `CheckDDLCoverage` 호출에서 `input.OpenAPIDoc` 인자 제거
- `CheckSSaCOpenAPI` 호출에 `input.OpenAPIDoc` 인자 추가 (ResponseFields 접근용)

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/crosscheck/ddl_coverage.go` | Rule 2 삭제, `buildSchemaProps` 삭제, 시그니처 변경 |
| `internal/crosscheck/ssac_openapi.go` | Rule 5, 6 추가 |
| `internal/crosscheck/crosscheck.go` | 호출부 시그니처 맞춤 |
| `internal/crosscheck/ddl_coverage_test.go` | 테스트 수정 (있으면) |

## 의존성

- `ssacparser.Sequence.Fields` — @response 필드 키
- `ssacparser.Sequence.Target` — shorthand 감지
- `ssacvalidator.SymbolTable.Operations[].ResponseFields` — OpenAPI response 필드

## 검증 방법

1. `go test ./internal/crosscheck/...` 통과
2. `go run ./cmd/fullend validate specs/gigbridge` — DDL → OpenAPI 23건 WARNING 소멸, 대신 SSaC → OpenAPI 필드 불일치가 있으면 ERROR로 표시
3. `go vet ./...` 통과
