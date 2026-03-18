# Phase001 — 정책 SSOT 포맷 전환: Rego → Go+toulmin

## 목표

OPA Rego 파일(`.rego`)을 폐기하고, Go 규칙 함수 + YAML 그래프 정의를 정책 SSOT로 도입한다. SSOT 10개 → 9개.

## 배경

기존: `specs/<project>/policy/authz.rego` — Rego DSL로 인가 규칙 선언
전환: `specs/<project>/authz/` — Go 함수(`//tm:backing`) + YAML 그래프 정의(defeats 관계)

toulmin 현재 문법에서 **규칙 함수는 판정 로직만 담당**하고, **역할(warrant/rebuttal/defeater)과 defeats 관계는 YAML 또는 Graph Builder에서 선언**한다. 동일한 함수가 그래프에 따라 다른 역할을 가질 수 있다.

## SSOT 포맷 정의

### 디렉토리 구조

```
specs/<project>/
├── authz/                       # 기존 policy/ 대체
│   ├── rules.go                 # 인가 규칙 함수 (//tm:backing 어노테이션)
│   ├── authz.yaml               # defeats 그래프 정의 (역할, qualifier, defeats)
│   └── ownership.go             # ownership 매핑 (기존 Rego 주석 → Go 구조체)
```

### 규칙 함수 (`rules.go`)

함수 시그니처: `func(claim any, ground any) (bool, any)` — 판정(bool) + 증거(evidence)

```go
package authz

//tm:backing "기본 거부 정책"
func DenyByDefault(claim, ground any) (bool, any) {
    return true, nil
}

//tm:backing "인증된 사용자만 워크플로우 실행 가능"
func ExecuteWorkflow(claim, ground any) (bool, any) {
    c := claim.(CheckRequest)
    if c.Action == "ExecuteWorkflow" && c.Resource == "workflow" && c.Claims.UserID > 0 {
        return true, nil
    }
    return false, nil
}

//tm:backing "admin은 모든 워크플로우 관리 가능"
func AdminAllowAll(claim, ground any) (bool, any) {
    c := claim.(CheckRequest)
    return c.Claims.Role == "admin", nil
}

//tm:backing "정지 계정은 모든 접근 차단"
func SuspendedBlock(claim, ground any) (bool, any) {
    c := claim.(CheckRequest)
    return c.Claims.Status == "suspended", nil
}

//tm:backing "소유자는 자기 리소스 접근 가능"
func OwnerAccess(claim, ground any) (bool, any) {
    c := claim.(CheckRequest)
    return c.ResourceOwnerID == c.Claims.UserID, nil
}
```

### 그래프 정의 (`authz.yaml`)

역할, qualifier, defeats 관계를 YAML로 선언. `toulmin graph authz.yaml`로 검증 + 코드 생성.

```yaml
graph: authz
rules:
  - name: DenyByDefault
    role: warrant
  - name: ExecuteWorkflow
    role: rebuttal
  - name: AdminAllowAll
    role: rebuttal
  - name: OwnerAccess
    role: rebuttal
  - name: SuspendedBlock
    role: rebuttal
defeats:
  - from: ExecuteWorkflow
    to: DenyByDefault
  - from: AdminAllowAll
    to: DenyByDefault
  - from: OwnerAccess
    to: DenyByDefault
  - from: SuspendedBlock
    to: AdminAllowAll
  - from: SuspendedBlock
    to: OwnerAccess
  - from: SuspendedBlock
    to: ExecuteWorkflow
```

### Ownership 매핑 (`ownership.go`)

```go
package authz

var Ownerships = []OwnershipMapping{
    {Resource: "workflow", Table: "workflows", Column: "org_id"},
    {Resource: "template", Table: "templates", Column: "org_id"},
    {Resource: "webhook", Table: "webhooks", Column: "org_id"},
}
```

### 기존 Rego와 1:1 매핑

| Rego | Go+toulmin |
|------|-----------|
| `allow if { input.action == "X" }` | `func X(claim, ground any) (bool, any)` + YAML `role: rebuttal` |
| `input.claims.role == "admin"` | `c.Claims.Role == "admin"` |
| `input.resource_owner` | `c.ResourceOwnerID == c.Claims.UserID` |
| `# @ownership res: table.col` | `var Ownerships = []OwnershipMapping{...}` |
| else 체이닝 (암묵적 우선순위) | YAML `defeats:` (명시적 관계) |
| Rego 단일 파일 내 규칙 순서 | defeats 그래프 (순서 무관) |

## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| `specs/dummys/zenflow-try06/policy/authz.rego` | 삭제 |
| `specs/dummys/zenflow-try06/authz/rules.go` (신규) | 규칙 함수 |
| `specs/dummys/zenflow-try06/authz/authz.yaml` (신규) | 그래프 정의 |
| `specs/dummys/zenflow-try06/authz/ownership.go` (신규) | ownership 매핑 |

## 의존성

- `github.com/park-jun-woo/toulmin/pkg/toulmin` — Graph Builder, EvalResult 타입 참조

## 검증 방법

- zenflow-try06 기존 Rego 규칙을 Go+toulmin으로 1:1 변환
- `toulmin graph authz.yaml --check` → 순환 없음, 참조 유효 확인
- 변환 후 기존 (action, resource) 쌍, role 값, claims 참조, ownership 매핑이 동일한지 수동 확인
- Phase002 이후 파서가 동일한 Policy 구조체를 생성하는지 자동 검증
