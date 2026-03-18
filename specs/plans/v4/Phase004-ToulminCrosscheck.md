# Phase004 — crosscheck 규칙 적응: Rego 기반 → toulmin 기반

## 목표

crosscheck의 Policy 관련 규칙 4개를 toulmin 기반 데이터 구조에 적응시킨다. D1(claims 검증 누락 탐지)이 이 Phase에서 자연스럽게 해소된다.

## 배경

Phase002에서 AllowRule에 `UsesClaims`, `Role`, `Defeats`, `Backing` 필드가 추가된다. YAML에서 role/defeats가 명시적이므로 기존 crosscheck 규칙이 단순화된다.

## 현재 Policy 관련 crosscheck 규칙 (4개)

| rules.go 위치 | 규칙 이름 | 호출 함수 |
|---------------|----------|----------|
| L51-57 | `Policy ↔ SSaC/DDL/States` | CheckPolicy() |
| L105-111 | `Policy → Config (claims)` | CheckClaimsRego() |
| L159-165 | `Policy → Config (roles)` | CheckRoles() |
| L168-174 | `Policy → DDL (roles)` | CheckRegoRoleDDL() |

## 규칙별 변경

### 1. CheckPolicy (유지, 내부 단순화)

`checkRegoPairsCoverage` → `checkAuthzPairsCoverage` 이름 변경.
`checkSSaCPairsCoverage` → 유지.
`checkOwnershipAnnotations` → Ownership이 Go 구조체이므로 AST 파싱 시 자동 검증. 단순화.
`checkOwnershipDDL` → 유지.

### 2. CheckClaimsRego → CheckClaimsAuthz (단순화)

기존: Rego 텍스트에서 `input.claims.xxx` 문자열을 regex로 수집 → fullend.yaml claims 값과 비교
신규: AllowRule.ClaimsRefs (Go AST에서 정확히 추출) → fullend.yaml claims 값과 비교

파일 이름 변경: `check_claims_rego.go` → `check_claims_authz.go`
유틸 삭제: `collect_rego_claims_refs.go` (AllowRule에 이미 포함)

### 3. CheckRoles (유지, 입력 변경)

기존: `collectPolicyRoles()` → regex로 수집
신규: AllowRule.RoleValue (Go AST에서 정확히 추출)

유틸 삭제: `collect_policy_roles.go`, `collect_rego_roles.go`

### 4. CheckRegoRoleDDL → CheckAuthzRoleDDL (이름 변경)

로직 동일. 파일 이름만 변경.

### 5. 신규: CheckClaimsPresence (D1 해소)

```go
func CheckClaimsPresence(policies []*policy.Policy, openapi *openapi3.T, serviceFuncs []ServiceFunc) []CrossError {
    for _, p := range policies {
        for _, rule := range p.Rules {
            if !rule.UsesClaims && guardsBearerAuthEndpoint(rule, openapi, serviceFuncs) {
                // bearerAuth endpoint인데 claims 미참조 → ERROR
            }
        }
    }
}
```

AllowRule.UsesClaims가 Phase002에서 Go AST 분석으로 추가되므로 D1이 자연스럽게 해소된다.

### 6. 신규: defeats 그래프 정합성 검증

YAML의 defeats 참조가 실제 존재하는 함수를 가리키는지 검증. `toulmin graph --check`가 이미 수행하지만, fullend crosscheck에서도 이중 검증.

## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| `internal/crosscheck/check_claims_rego.go` | → `check_claims_authz.go` 이름 변경 + 단순화 |
| `internal/crosscheck/check_rego_pairs_coverage.go` | → `check_authz_pairs_coverage.go` 이름 변경 |
| `internal/crosscheck/check_rego_role_ddl.go` | → `check_authz_role_ddl.go` 이름 변경 |
| `internal/crosscheck/collect_rego_claims_refs.go` | 삭제 |
| `internal/crosscheck/collect_policy_claims_refs.go` | 삭제 |
| `internal/crosscheck/collect_policy_roles.go` | 삭제 |
| `internal/crosscheck/collect_rego_roles.go` | 삭제 |
| `internal/crosscheck/check_claims_presence.go` (신규) | D1 규칙 |
| `internal/crosscheck/rules.go` | 규칙 이름/함수 참조 갱신 + 신규 규칙 등록 |
| `internal/crosscheck/claims_test.go` | 테스트 갱신 |
| `internal/crosscheck/rego_role_ddl_test.go` | → `authz_role_ddl_test.go` 갱신 |

## 의존성

- Phase002 완료 필수 (AllowRule 확장)

## 검증 방법

- `go test ./internal/crosscheck/...` 통과
- zenflow-try06에 대해 `fullend validate` 실행 → 기존 검증 결과와 동일 + D1 ERROR 추가 탐지
