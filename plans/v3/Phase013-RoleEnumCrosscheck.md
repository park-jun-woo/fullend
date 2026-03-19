# Phase013: role enum 교차 검증 ✅ 완료

## 목표

`fullend.yaml`에 선언된 role enum과 OPA Rego의 `input.role` 값을 교차 검증하여 role 오타를 검출한다.

## 배경

MUT-13에서 Rego `input.role == "client"` → `"Client"` 변경 시 crosscheck이 미검출. role 값의 원천이 어디에도 SSOT로 선언되지 않았기 때문.

## 기존 코드 분석

- `internal/policy/parser.go`에 이미 `reRoleRef` 정규식으로 `input.role == "xxx"` 패턴을 추출하고 있음
- `AllowRule.RoleValue` 필드에 role 값이 이미 저장됨
- **policy 파서 수정은 불필요** — `AllowRule.RoleValue`에서 직접 수집하면 됨

## 변경 내용

### 1. `specs/gigbridge/fullend.yaml` — roles 필드 추가

```yaml
auth:
  secret_env: JWT_SECRET
  claims:
    ID: user_id
    Email: email
    Role: role
  roles: [client, freelancer]  # NEW
```

### 2. `internal/projectconfig/projectconfig.go` — roles 파싱

`Auth` struct에 `Roles []string` 필드 추가.

```go
type Auth struct {
    SecretEnv string            `yaml:"secret_env"`
    Claims    map[string]string `yaml:"claims"`
    Roles     []string          `yaml:"roles"`
}
```

### 3. `internal/crosscheck/crosscheck.go` — `CrossValidateInput`에 roles 전달

```go
type CrossValidateInput struct {
    // ...
    Roles []string // from fullend.yaml auth.roles
}
```

`Run()`에서 `CheckRoles()` 호출 추가:
```go
// Roles ↔ Policy
if len(input.Policies) > 0 && len(input.Roles) > 0 {
    errs = append(errs, CheckRoles(input.Policies, input.Roles)...)
}
```

### 4. `internal/crosscheck/roles.go` — 신규

```go
// CheckRoles validates that OPA Rego input.role values match fullend.yaml auth.roles.
func CheckRoles(policies []*policy.Policy, roles []string) []CrossError
```

로직:
1. `roles`가 비어있으면 스킵 (선택적 기능)
2. 각 policy의 `AllowRule.RoleValue`에서 role 값 수집 (이미 파서가 추출해둠)
3. 수집된 role 값이 `roles` 목록에 있는지 확인
4. 없으면 ERROR: `Rego role "Client"가 fullend.yaml auth.roles에 없습니다`
5. 역방향: roles에 있지만 Rego에서 한 번도 안 쓰인 role → WARNING

### 5. `internal/orchestrator/validate.go` — roles 전달

`CrossValidateInput` 생성 시 `projConfig.Backend.Auth.Roles`를 전달. 기존 `claims` 추출 패턴과 동일:
```go
var roles []string
if projConfig != nil && projConfig.Backend.Auth != nil {
    roles = projConfig.Backend.Auth.Roles
}
```

### 6. 테스트

#### `internal/crosscheck/roles_test.go` — crosscheck 테스트
- roles=[client, freelancer] + rego=[client, freelancer] → 통과
- roles=[client, freelancer] + rego=[client, Client] → ERROR "Client"
- roles=[client, freelancer, admin] + rego=[client] → WARNING "freelancer", "admin" 미사용

## 영향 없는 범위

- roles 미선언 시 기존 동작 그대로 (스킵)
- policy 파서 — 변경 없음 (AllowRule.RoleValue 이미 존재)
- SSaC/OpenAPI/DDL — 변경 없음
- 기존 crosscheck 규칙 — 변경 없음

## 검증

```
go test ./internal/crosscheck/...
go test ./...
fullend validate specs/gigbridge  # roles crosscheck 통과 확인
```

MUT-13 재실행: Rego `"client"` → `"Client"` 변경 시 ERROR 검출 확인.
