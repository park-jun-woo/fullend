# Phase014: 뮤테이션 테스트 Phase 2 미검출 수정 ✅ 완료

## 목표

MUT-16(claims key 미검증), MUT-17(middleware nil 스킵), MUT-21(OpenAPI 유령 property), MUT-22(sensitive 패턴 부족)를 수정한다.
MUT-19는 sed 오류였고 실제로는 검출됨 (재테스트 PASS).

## MUT-16: claims key 교차 검증

### 문제

`fullend.yaml`의 `claims: {ID: user_id}`에서 `user_id` → `userId`로 변경해도 crosscheck이 못 잡음. 코드젠이 `c.Get("userId")`를 생성하지만 JWT에는 `user_id`로 들어있어 항상 빈 값 — 모든 `currentUser.ID` 참조가 조용히 실패.

### 검증 경로

```
fullend.yaml claims: {ID: user_id, Role: role}
    ↕ Rego input.claims.user_id 참조와 일치?
    ↕ pkg/auth IssueToken에서 JWT claim key와 일치?
```

### 수정

`internal/crosscheck/claims.go`의 `CheckClaims`에 역방향 검증 추가:

1. **Claims → Rego**: fullend.yaml claims의 value (e.g., `user_id`, `role`)가 Rego에서 `input.claims.xxx`로 참조될 때, `xxx`가 claims value 목록에 있는지 확인.

Rego에서 `input.claims.role`을 파싱하는 로직이 필요. 기존 `reRoleRef` 정규식은 role만 추출하므로, 범용 `input.claims.xxx` 참조 추출이 필요.

`internal/policy/parser.go`에 claims 참조 추출 추가:
```go
var reClaimsRef = regexp.MustCompile(`input\.claims\.(\w+)`)
```

`AllowRule` 또는 `Policy`에 `ClaimsRefs []string` 필드 추가.

`CheckClaims`에서:
- Rego claims 참조 → fullend.yaml claims values에 있는지 (ERROR)
- fullend.yaml claims values → Rego에서 한 번도 안 쓰인 건 (WARNING)

## MUT-17: middleware nil vs 빈 배열 구분

### 문제

YAML에서 `middleware:` (값 없음) → Go에서 `nil` → `crosscheck.go:75` 조건 `input.Middleware != nil`이 false → 검증 통째로 스킵.

`middleware: []` (명시적 빈 배열)이면 검출 10건. `middleware:` (값 없음)이면 0건.

### 수정

`crosscheck.go:75`의 조건을 변경:

```go
// 기존
if input.OpenAPIDoc != nil && input.Middleware != nil {

// 수정: OpenAPI 있으면 항상 검증
if input.OpenAPIDoc != nil {
```

middleware가 `nil`이면 mwSet가 비어서 securityScheme 불일치가 검출됨. middleware 키 자체가 없는 프로젝트도 OpenAPI에 securitySchemes가 있으면 경고해야 맞음.

### 영향

middleware를 안 쓰고 securitySchemes도 없는 프로젝트 → 변경 없음 (둘 다 비어서 에러 0건).

## MUT-21: OpenAPI 유령 property 검출

### 문제

OpenAPI Gig 스키마에 `rating: integer`를 추가해도 DDL `gigs`에 해당 컬럼이 없는 걸 못 잡음. 프론트엔드 개발자가 OpenAPI를 계약서로 보고 `rating`을 소비하는 UI를 만들면, 실제 API 응답에 안 와서 "왜 null이지?" 디버깅 지옥.

### 현재 상태

DDL → OpenAPI 방향만 검증 (DDL 컬럼이 OpenAPI에 있는지). 역방향 OpenAPI → DDL은 없음.

### 정당한 예외

OpenAPI에 DDL 컬럼이 아닌 property가 있을 수 있는 경우:
- `x-include` FK로 선언된 관계 필드 (e.g., `client` = User JOIN 결과)
- `@dto` 모델의 필드 (DDL 테이블이 아닌 커스텀 타입)

### 수정

`internal/crosscheck/openapi_ddl.go`의 `CheckOpenAPIDDL`에 역방향 검증 추가:

```go
// OpenAPI schema property → DDL column 존재 확인
// 예외: x-include로 선언된 필드, @dto 모델
for propName := range schema.Properties {
    if xIncludeFields[propName] {
        continue // FK join 필드 — 정당한 확장
    }
    if _, colExists := table.Columns[propName]; !colExists {
        // ERROR: 유령 property
    }
}
```

`x-include`에서 선언된 FK 필드명을 수집하여 예외 처리. 그 외 DDL에도 없고 x-include에도 없는 property는 **ERROR**.

### 레벨: ERROR

유령 property는 프론트엔드가 정당한 계약이라 믿고 소비하다 실패하는 케이스. WARNING이 아니라 ERROR.

## MUT-22: sensitive 패턴 보강

### 문제

현재 패턴: `["password", "secret", "hash", "token"]`

다음 변형은 못 잡음:
- `pw` (password 축약)
- `credential`, `cred`
- `ssn` (주민등록번호)
- `pin`
- `key` (API key 등)
- `salt`
- `otp`
- `private`
- `auth_code`
- `refresh` (refresh token)
- `access` (access token)
- `bearer`
- `cert` (certificate)
- `passphrase`
- `seed` (mnemonic seed)
- `nonce`
- `iv` (initialization vector) — 너무 짧아 false positive 위험
- `encrypted`
- `cipher`
- `digest`
- `signature`, `sig`
- `credit_card`, `card_number`, `cvv`, `expiry`
- `bank_account`, `routing_number`
- `phone` (PII)
- `address` (PII) — false positive 높음
- `birth`, `dob` (PII)
- `passport`
- `license_number`
- `biometric`

### 수정 기준

false positive를 줄이면서 실질적 위험이 높은 것만 추가:

```go
var sensitivePatterns = []string{
    // 인증 정보
    "password", "passwd", "passphrase",
    "secret", "token", "hash", "salt",
    "credential", "otp", "pin",
    // 암호화
    "private_key", "cipher", "encrypted",
    // 금융
    "credit_card", "card_number", "cvv",
    "bank_account", "routing_number",
    // 개인식별
    "ssn", "passport", "license_number",
    "biometric",
}
```

`key`, `phone`, `address`, `birth` 등은 false positive가 높아 제외. 필요 시 `@sensitive` 수동 태깅.

### 테스트

기존 테스트에 새 패턴 케이스 추가:
- `user_credential` → 검출
- `otp_code` → 검출
- `api_key` → 미검출 (의도적 제외, false positive)

## 변경 파일

| 파일 | 변경 |
|---|---|
| `internal/policy/parser.go` | `reClaimsRef` 정규식 + `Policy.ClaimsRefs` 추가 |
| `internal/crosscheck/claims.go` | Rego claims 참조 ↔ fullend.yaml claims 대조 |
| `internal/crosscheck/crosscheck.go` | middleware 조건 변경 |
| `internal/crosscheck/openapi_ddl.go` | OpenAPI→DDL 역방향 property 검증 추가 |
| `internal/crosscheck/sensitive.go` | sensitivePatterns 목록 보강 |
| `files/mutest-phase2.md` | 결과 업데이트 |

## 검증

```
go test ./internal/policy/...
go test ./internal/crosscheck/...
go test ./...
```

MUT-16, MUT-17, MUT-21, MUT-22 재실행하여 검출 확인.
