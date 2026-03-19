✅ 완료

# Phase029: STML ↔ OpenAPI 검증을 crosscheck로 분리

## 목표

STML validator에 섞여 있는 cross-SSOT 검증(STML ↔ OpenAPI)을 crosscheck로 이동하여,
"개별 SSOT validator = 내부 정합성, crosscheck = SSOT 간 정합성" 원칙을 일관되게 적용한다.

## 동기

현재 아키텍처 위반:
- 다른 모든 SSOT validator(SSaC, DDL, Policy, States 등)는 내부 정합성만 검증
- crosscheck가 SSOT 간 교차 검증 전담
- **STML validator만** 내부에서 OpenAPI를 직접 로드하여 cross-SSOT 검증 수행 중

이로 인해:
1. STML validator가 `LoadOpenAPI()`로 OpenAPI를 독자 파싱 — orchestrator의 ParseAll()과 중복
2. STML validator 자체 `SymbolTable` 타입이 crosscheck의 `SymbolTable`과 별도 존재
3. STML ↔ OpenAPI 에러가 reporter의 Cross 섹션이 아닌 STML 섹션에 출력

## 설계

### 1단계: STML validator에서 OpenAPI 의존성 제거

현재 STML validator가 수행하는 OpenAPI 관련 검증:

| 검증 | 내용 |
|---|---|
| operationId 존재 | `st.Operations[operationID]` lookup |
| HTTP method 일치 | Fetch=GET, Action≠GET |
| parameter 존재 | OpenAPI parameters에 포함 여부 |
| request field 존재 | request schema에 필드 존재 |
| response field 존재 | response schema에 필드 존재 (+ custom.ts 폴백) |
| 배열 타입 확인 | `data-each` 대상이 배열인지 |
| x-pagination 존재 | `data-paginate` 시 extension 필요 |
| x-sort/x-filter 허용 | allowed 목록 포함 여부 |

**모두** crosscheck로 이동 대상.

### 2단계: crosscheck에 CheckSTMLOpenAPI 추가

```go
// stml_openapi.go
func CheckSTMLOpenAPI(
    pages []stmlparser.PageSpec,
    doc *openapi3.T,
) []CrossError
```

기존 STML validator의 `validateFetchBlock()`, `validateActionBlock()` 로직을
crosscheck 컨텍스트로 재작성. `stml/validator/symbol.go`의 `SymbolTable` 구조를
crosscheck 내부에서 OpenAPI doc으로부터 직접 구성하거나, 공용 유틸로 추출.

### 3단계: CrossValidateInput 확장

```go
type CrossValidateInput struct {
    // ... 기존 필드 ...
    STMLPages []stmlparser.PageSpec // STML 파서 결과
}
```

### 4단계: STML validator 축소

`stml/validator/validator.go`에서 OpenAPI 관련 코드 제거:
- `LoadOpenAPI()` 호출 제거
- `validateFetchBlock()`, `validateActionBlock()`에서 OpenAPI 검증 제거
- 남는 검증: component TSX 파일 존재, custom.ts 폴백 (파일시스템 검증)

`stml/validator/symbol.go` — crosscheck로 이동 또는 공용 추출 후 삭제.

### 5단계: orchestrator 연결

`orchestrator/validate.go`에서:
- STML 파서 결과를 `CrossValidateInput.STMLPages`에 전달
- STML validator 호출 시 `projectRoot` 대신 축소된 인터페이스 사용

## 변경 파일

| 파일 | 변경 |
|---|---|
| `internal/crosscheck/stml_openapi.go` | 신규 — CheckSTMLOpenAPI |
| `internal/crosscheck/stml_openapi_test.go` | 신규 |
| `internal/crosscheck/types.go` | CrossValidateInput에 STMLPages 추가 |
| `internal/crosscheck/rules.go` | STML → OpenAPI rule 등록 |
| `internal/stml/validator/validator.go` | OpenAPI 검증 로직 제거, 파일시스템 검증만 잔류 |
| `internal/stml/validator/symbol.go` | 삭제 또는 crosscheck로 이동 |
| `internal/stml/validator/errors.go` | OpenAPI 관련 에러 함수 제거 |
| `internal/stml/validator/validator_test.go` | OpenAPI 관련 테스트 crosscheck로 이동 |
| `internal/orchestrator/validate.go` | STMLPages를 CrossValidateInput에 전달 |

## 의존성

Phase027 이후. Phase028(genapi 분리)과 독립.

## 검증

1. `go test ./internal/crosscheck/...` — 신규 STML ↔ OpenAPI 테스트 통과
2. `go test ./internal/stml/...` — 축소된 validator 테스트 통과
3. `go run ./cmd/fullend validate specs/gigbridge` — STML 에러가 Cross 섹션에 출력
4. `go vet ./...` 통과

## 리스크

- STML validator의 `custom.ts` 폴백 로직은 파일시스템 접근이 필요 — crosscheck는 현재 파일시스템 비의존이므로, 이 검증은 STML validator에 잔류시키거나 crosscheck에 projectRoot를 전달해야 함
- 기존 STML 에러 메시지 포맷이 바뀌면 사용자 혼란 가능 — 메시지 내용 유지 권장
