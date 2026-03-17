# Phase045: 뮤테이션 테스트 FAIL 수정 ✅ 완료

## 목표

`files/mutest-report01.md` 확정 FAIL을 수정한다.

## 사전 확인: report ↔ mutest source 불일치

| ID | report01 | mutest source | 판단 |
|---|---|---|---|
| MUT-CONFIG-OPENAPI-001 | FAIL | PASS (Phase014 수정 완료) | **제외** — 이미 수정됨. report 오류 |
| MUT-FUNC-002 | FAIL | PASS (TODO/panic 감지) | **확인 필요** — `isStubBody` 코드에 panic 분기 없음. source 기록이 선행 기록일 가능성 |
| MUT-SSAC-002 | FAIL | 결과 미기재 | 수정 대상 |
| MUT-SSAC-025 | FAIL | 결과 미기재 | 수정 대상 |
| MUT-SSAC-035 | FAIL | 결과 미기재 | 수정 대상 |

→ 확정 수정 대상: **4건** (MUT-FUNC-002, MUT-SSAC-002, MUT-SSAC-025, MUT-SSAC-035)
→ MUT-CONFIG-OPENAPI-001은 `rules.go:70` Requires가 `in.OpenAPIDoc != nil`로 이미 middleware 무관하게 호출됨. `CheckMiddleware` 코드도 빈 middleware에서 schemeNames 에러를 정상 보고. 재실행으로 PASS 확인 후 제외.

## 변경 파일 목록

### 1. MUT-FUNC-002 — `panic("TODO")` stub body 미감지

- **파일**: `internal/funcspec/is_stub_body.go`
- **chain**: `isStubBody` ← `processDecl` ← `ParseFile`
- **현상**: `isStubBody`가 빈 body(len=0)와 단일 return만 stub 판정. 단일 `panic(...)` 호출을 놓침
- **수정**: 단일 statement가 `*ast.ExprStmt` → `*ast.CallExpr` → ident `panic`인 경우도 stub 판정
- **테스트**: `internal/funcspec/parser_test.go`에 `panic("TODO")` body 케이스 추가

### 2. MUT-SSAC-002 — @get result type "Gigs" 변이 미검출

- **파일**: `internal/crosscheck/check_result_type.go`
- **chain**: `checkResultType` ← `checkSSaCDDLFunc` ← `CheckSSaCDDL`
- **현상**: `checkSSaCDDLFunc`의 skip 조건은 정상. 실제 원인은 `modelToTable` 변환:
  - `modelToTable("Gig")` → `pascalToSnake("Gig")` → `"gig"` → `inflection.Plural("gig")` → `"gigs"` ✓
  - `modelToTable("Gigs")` → `pascalToSnake("Gigs")` → `"gigs"` → `inflection.Plural("gigs")` → `"gigs"` ✓
  - **"Gig"과 "Gigs" 둘 다 동일 테이블 "gigs"에 매핑** → 변이 통과
- **수정**: `checkResultType`에 모델명 단수형 검증 추가. `inflection.Singular(typeName) != typeName`이면 WARNING ("result type은 단수형 모델명이어야 함"). `"github.com/jinzhu/inflection"` import 추가 필요 (현재 이 파일에 없음, `model_to_table.go`에서만 사용 중)
- **테스트**: result type "Gigs"(복수형) 입력 시 WARNING 검출 케이스 추가

### 3. MUT-SSAC-025 — @state transition 이벤트명 대소문자 미검출

- **파일**: `internal/crosscheck/check_func_guard_states.go`
- **chain**: `checkFuncGuardStates` ← `checkGuardStates` ← `CheckStates`
- **현상**: `checkFuncGuardStates`가 `d.ValidFromStates(fn.Name)`으로 검증. 검증 대상이 `seq.Transition`이어야 하는데 `fn.Name`(함수명)을 사용. gigbridge에서는 transition="PublishGig"과 fn.Name="PublishGig"이 동일하여 우연히 통과했으나, `parser_test.go:158`의 transition="cancel" vs fn.Name="CancelReservation" 케이스에서는 기존 코드가 오탐 ERROR를 낸다
- **수정**: `d.ValidFromStates(fn.Name)`을 `d.ValidFromStates(seq.Transition)`으로 **교체**. 에러 메시지도 fn.Name → seq.Transition으로 변경
- **테스트**: transition="publishGig"(소문자 변이), diagram event="PublishGig" 시 ERROR 검출 + transition="cancel", diagram event="cancel" 시 정상 통과 케이스 추가

### 4. MUT-SSAC-035 — @call pkg 함수명 소문자 미검출

- **파일**: `internal/crosscheck/check_service_func_calls.go`
- **chain**: `checkServiceFuncCalls` → `parseCallKey` (co-called) → `checkSingleCall` → `validateCallSpec`
- **현상**: `parseCallKey`가 `strcase.ToGoCamel("issueToken")` → `"IssueToken"`으로 정규화하여 spec 매칭 성공. SSaC 원본의 소문자 함수명이 묵인됨
- **수정**: `checkServiceFuncCalls`에서 `parseCallKey` 호출 전에 `strings.SplitN(seq.Model, ".", 2)`로 원본 함수명을 추출, 첫 글자가 소문자이면 ERROR 보고. `parseCallKey` 내부의 `callParts[1]`은 외부에서 접근 불가하므로 인라인 split 사용. Go 소문자 시작 함수는 unexported이므로 패키지 외부 호출 불가 — 컨벤션이 아니라 언어 규칙 위반이므로 ERROR. 패키지 없는 `@call SomeFunc(...)` 케이스도 처리: `len(parts) == 1`이면 `parts[0]`, `len(parts) == 2`이면 `parts[1]`의 첫 글자 검사
- **테스트**: `@call auth.issueToken` 시 ERROR 검출 + `@call someFunc(...)` (패키지 없음) 시 ERROR 검출 케이스 추가

## 의존성

- 추가 외부 패키지 없음
- 기존 `inflection`, `strcase`, `go/ast`, `go/token` 그대로 사용

## 검증 방법

1. `go test ./internal/funcspec/... ./internal/crosscheck/...` — 단위 테스트 전체 통과
2. `go test ./...` — 전체 테스트 통과
3. MUT-CONFIG-OPENAPI-001 재실행하여 PASS 확인 (report 오류 정정)
4. 나머지 4건 뮤테이션 재실행하여 PASS 전환 확인
