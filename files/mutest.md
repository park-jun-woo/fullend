# Mutation Test — gigbridge crosscheck 검출력 검증

## 목적

gigbridge SSOT에 미묘한 오류를 하나씩 주입하고 `fullend validate`가 검출하는지 확인한다.
변경은 최대한 교묘하게 — 대소문자 한 글자, 언더스코어 추가/제거, 복수형 등.

## 방법

1. 변경 적용
2. `go run ./cmd/fullend validate specs/gigbridge` 실행
3. 검출 여부 기록 (PASS = 잡음, FAIL = 못 잡음)
4. `git checkout -- specs/gigbridge/` 로 되돌림

---

## 시나리오

### MUT-01: OpenAPI operationId 대소문자 (SSaC ↔ OpenAPI)
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: `operationId: Login` → `operationId: login`
- 기대: ERROR — SSaC function "Login"과 불일치

### MUT-02: OpenAPI 응답 property 대소문자 (shorthand @response ↔ OpenAPI)
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: Login 응답 `AccessToken` → `accessToken`
- 기대: ERROR — json tag `access_token`과 불일치 (둘 다 틀림)

### MUT-03: DDL 컬럼명 언더스코어 (DDL ↔ OpenAPI)
- 대상: `specs/gigbridge/db/gigs.sql`
- 변경: `client_id` → `clientid` (언더스코어 제거)
- 기대: ERROR — OpenAPI x-model Gig의 `client_id`와 불일치

### MUT-04: 상태명 대소문자 (States ↔ SSaC)
- 대상: `specs/gigbridge/states/gig.md`
- 변경: `draft --> open : PublishGig` → `Draft --> open : PublishGig`
- 기대: ERROR — SSaC @state draft와 불일치

### MUT-05: 상태 전이 함수명 오타 (States ↔ SSaC)
- 대상: `specs/gigbridge/states/gig.md`
- 변경: `draft --> open : PublishGig` → `draft --> open : Publishgig`
- 기대: ERROR — SSaC function "PublishGig"과 불일치

### MUT-06: Policy action명 대소문자 (Policy ↔ SSaC)
- 대상: `specs/gigbridge/policy/authz.rego`
- 변경: `input.action == "PublishGig"` → `input.action == "publishGig"`
- 기대: ERROR — SSaC @auth action "PublishGig"과 불일치

### MUT-07: Policy resource명 오타 (Policy ↔ SSaC)
- 대상: `specs/gigbridge/policy/authz.rego`
- 변경: `input.resource == "gig"` → `input.resource == "gigs"`
- 기대: ERROR — SSaC @auth resource "gig"과 불일치

### MUT-08: SSaC 함수명 변경 (SSaC ↔ OpenAPI)
- 대상: `specs/gigbridge/service/gig/create_gig.ssac`
- 변경: `func CreateGig()` → `func Creategig()`
- 기대: ERROR — OpenAPI operationId "CreateGig"과 불일치

### MUT-09: @response 필드명 변경 (SSaC @response ↔ OpenAPI)
- 대상: `specs/gigbridge/service/gig/create_gig.ssac`
- 변경: `@response { gig: gig }` → `@response { Gig: gig }`
- 기대: ERROR — OpenAPI 응답 property "gig"과 불일치

### MUT-10: DDL 테이블명 변경 (DDL ↔ SSaC)
- 대상: `specs/gigbridge/db/gigs.sql`
- 변경: `CREATE TABLE gigs` → `CREATE TABLE gig`
- 기대: ERROR — SSaC @get Gig.List의 테이블 "gigs"와 불일치

### MUT-11: DDL 컬럼 추가 누락 (DDL ↔ OpenAPI)
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: Gig schema에 `budget` property 제거
- 기대: WARNING — DDL gigs.budget이 OpenAPI에 없음

### MUT-12: 상태 전이 누락 (States ↔ SSaC)
- 대상: `specs/gigbridge/states/gig.md`
- 변경: `under_review --> disputed : RaiseDispute` 행 삭제
- 기대: ERROR — SSaC RaiseDispute의 @state 전이가 States에 없음

### MUT-13: Policy role명 변경 (Policy ↔ SSaC)
- 대상: `specs/gigbridge/policy/authz.rego`
- 변경: `input.role == "client"` → `input.role == "Client"` (PublishGig 규칙)
- 기대: 현재 crosscheck 범위 밖 (role 값은 검증 안 함) — FAIL 예상

### MUT-14: OpenAPI path 변경 (Hurl ↔ OpenAPI)
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: `/auth/login` → `/auth/Login`
- 기대: ERROR — Hurl scenario의 path와 불일치

### MUT-15: SSaC @call 패키지명 오타 (Func ↔ SSaC)
- 대상: `specs/gigbridge/service/auth/login.ssac`
- 변경: `auth.IssueToken` → `Auth.IssueToken`
- 기대: ERROR — funcspec package "auth"와 불일치

---

## 결과 기록

| ID | 변경 | 기대 | 실제 | 비고 |
|---|---|---|---|---|
| MUT-01 | operationId Login→login | PASS | PASS | SSaC→OpenAPI + OpenAPI→SSaC 양방향 검출 |
| MUT-02 | AccessToken→accessToken | PASS | PASS | Phase012 shorthand 검증이 access_token vs accessToken 잡음 |
| MUT-03 | client_id→clientid | PASS | PASS | x-include FK + Policy ownership 컬럼 부재 검출 |
| MUT-04 | draft→Draft | PASS | PASS | States 상태명 대소문자 불일치 검출 |
| MUT-05 | PublishGig→Publishgig | PASS | PASS | States↔SSaC + States↔OpenAPI 3건 검출 |
| MUT-06 | action "PublishGig"→"publishGig" | PASS | PASS | Policy↔SSaC 양방향 검출 |
| MUT-07 | resource "gig"→"gigs" | PASS | PASS | Policy↔SSaC resource 불일치 검출 |
| MUT-08 | func CreateGig→Creategig | PASS | PASS | SSaC→OpenAPI + OpenAPI→SSaC 검출 |
| MUT-09 | @response { gig→Gig } | PASS | PASS | @response→OpenAPI + OpenAPI→@response 양방향 |
| MUT-10 | TABLE gigs→gig | PASS | PASS | SSaC @result↔DDL 9건 + index 누락 등 대량 검출 |
| MUT-11 | OpenAPI Gig budget 제거 | PASS | PASS | SSaC @post budget 필드 OpenAPI 부재 검출 |
| MUT-12 | 전이 RaiseDispute 삭제 | PASS | PASS | States↔SSaC 전이 누락 검출 |
| MUT-13 | role "client"→"Client" | FAIL | PASS | Phase013 구현 후 검출 성공 |
| MUT-14 | path /auth/login→/auth/Login | PASS | PASS | Hurl→OpenAPI 6건 path 불일치 검출 |
| MUT-15 | auth.IssueToken→Auth.IssueToken | PASS | PASS | Func↔SSaC 패키지 불일치 검출 |

**결과: 15/15 검출 (100%)** — Phase013 구현 후 MUT-13도 검출 성공
