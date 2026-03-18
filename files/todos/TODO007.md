# TODO007: codegen 설계 미비 (2층 이슈)

zenflow-try06 전수 리뷰(codegen-review2.md)에서 발견된 **codegen 아키텍처 수준 결함**.
1층(기계적 버그)과 3층(스펙 설계)은 제외. fullend 엔진 구조 변경이 필요한 것만 수록.

---

## D1. Rego 규칙 claims 검증 누락 탐지

### 현상

bearerAuth로 보호된 endpoint의 Rego allow 규칙이 `input.claims`를 한 번도 참조하지 않으면, 인증 필요로 선언해놓고 실제로는 비인증 접근을 허용하는 구멍이 생긴다.

zenflow-try06에서 8개 엔드포인트가 이 상태:
- ExecuteWorkflow, ListWorkflowVersions, ExecuteWithReport, GetExecutionReport
- CloneTemplate, ListWebhooks, GetSchedule, ListExecutionLogs

### 현재 crosscheck 상태

| 규칙 | 파일 | 하는 일 |
|------|------|---------|
| `Policy → Config (claims)` | check_claims_rego.go | Rego `input.claims.xxx` 참조가 fullend.yaml claims 값과 일치하는지 검증 |
| `Policy ↔ SSaC` | check_rego_pairs_coverage.go | Rego allow rule ↔ SSaC authorize 시퀀스 매칭 |
| `Config → OpenAPI` | check_middleware.go | fullend.yaml middleware ↔ OpenAPI securitySchemes 일치 |

**빠진 검증**: "bearerAuth endpoint의 Rego rule이 `input.claims`를 참조하는가?" — 이 교차가 없음.

### 검증 로직

```
for each Rego allow rule (action, resource):
    해당 (action, resource)에 대응하는 OpenAPI endpoint를 찾는다
    그 endpoint에 security: [bearerAuth] 가 있는가?
    있다면 → 이 allow rule이 input.claims를 한 번이라도 참조하는가?
    참조 0회 → ERROR
```

### 필요한 데이터

- **AllowRule에 claims 참조 여부 추가**: 현재 `AllowRule`에 `UsesOwner`, `UsesRole`은 있지만 `UsesClaims`가 없음. rule-level claims 참조 추적 필요
- **OpenAPI endpoint → (action, resource) 매핑**: SSaC `@authorize` 에서 추출하는 (action, resource) 쌍과 OpenAPI endpoint의 security 설정을 연결해야 함. 현재 `buildSSaCAuthPairs`가 이 매핑을 만들지만 endpoint 정보는 포함하지 않음

### 변경 대상 파일

| 파일 | 변경 |
|------|------|
| `internal/policy/allow_rule.go` | `UsesClaims bool` 필드 추가 |
| `internal/policy/process_allow_block.go` | allow 블록 파싱 시 `input.claims` 참조 탐지 |
| `internal/crosscheck/` (신규) | `CheckRegoClaimsPresence` — bearerAuth endpoint의 Rego rule claims 참조 검증 |
| `internal/crosscheck/rules.go` | 새 규칙 등록 |
| `internal/crosscheck/cross_validate_input.go` | OpenAPIDoc 이미 포함되어 있으므로 추가 불필요 |

### 에러 메시지 예시

```
[Policy ↔ OpenAPI] ERROR: Rego allow rule "ExecuteWorkflow" (resource: workflow)
  guards bearerAuth endpoint POST /workflows/{id}/execute
  but never references input.claims — effectively unauthenticated
```

### 판단 기준

| OpenAPI security | Rego input.claims 참조 | 판정 |
|-----------------|----------------------|------|
| bearerAuth 있음 | 1회 이상 | OK |
| bearerAuth 있음 | 0회 | **ERROR** |
| security 없음 (public) | 0회 | OK |
| security 없음 (public) | 1회 이상 | WARNING (dead check) |

---

## D2. 트랜잭션 경계와 외부 부수효과 위치

### 현상

SSaC에서 `@publish`, `@call Func` (외부 HTTP 호출 등)의 실행 시점이 트랜잭션 경계 기준으로 구분되지 않음. codegen이 모든 시퀀스를 트랜잭션 내부에 배치하여:

1. `queue.Publish()` → `tx.Commit()` 순서 — Commit 실패 시 발행된 메시지 회수 불가
2. `webhook.Deliver()` (외부 HTTP) → `tx.Commit()` — 외부 지연이 DB 트랜잭션 점유

### 수정 방향

**방안 A**: SSaC 문법에 `@after-commit` 블록 도입

```ssac
@transaction {
  @put Workflow.UpdateStatus(...)
  @put ExecutionLog.Create(...)
}
@after-commit {
  @publish "workflow.executed" { ... }
  @call webhook.Deliver(...)
}
```

**방안 B**: codegen이 `@publish`와 외부 `@call`을 자동으로 Commit 이후로 배치 (암묵적 규칙)

방안 B가 SSaC 문법 변경 없이 가능하지만, 명시성이 떨어짐.

---

## D3. ResourceID 바인딩 규칙 부재

### 현상

`@authorize` 시퀀스에서 `ResourceID`에 무엇을 넣을지 규칙이 없음. codegen이 SSaC를 그대로 옮기면:

- create: 아직 리소스가 없으므로 `currentUser.ID` 전달 (잘못됨)
- list: 특정 리소스가 없으므로 `currentUser.ID` 전달 (잘못됨)
- get/update/delete: 조회한 리소스의 ID 전달 (정상)

### 수정 방향

1. create/list에서는 `ResourceID`를 생략하거나 0으로 전달하는 관례 정립
2. crosscheck에서 `@authorize ResourceID: currentUser.ID` + action이 create/list일 때 WARNING
3. 또는 SSaC에서 `ResourceID: _` (생략 표기) 문법 도입

---

## D4. 프론트엔드 인증 흐름 미생성

### 현상

STML/codegen이 로그인 API 호출 코드는 생성하지만:
- 토큰 저장소 (localStorage/메모리) 미생성
- fetch wrapper에 Authorization 헤더 자동 부착 미생성
- 토큰 만료 시 refresh 흐름 미생성
- 401 응답 시 로그인 페이지 리다이렉트 미생성

### 영향

인증 필요 API 전부 401. 프론트엔드 사실상 미작동.

### 수정 방향

fullend.yaml에 `frontend.auth` 설정이 있거나, OpenAPI에 bearerAuth가 정의되어 있으면:
- `api.ts`에 토큰 저장/주입 fetch wrapper 자동 생성
- login 핸들러에 토큰 저장 코드 생성
- 401 인터셉터 생성

---

## 우선순위

| 이슈 | 긴급도 | 난이도 | 비고 |
|------|--------|--------|------|
| D1. Rego claims 검증 | 높음 | 낮음 | crosscheck 규칙 1개 추가 |
| D2. 트랜잭션 경계 | 높음 | 중간 | SSaC 문법 확장 또는 codegen 암묵 규칙 |
| D3. ResourceID 바인딩 | 중간 | 낮음 | crosscheck WARNING + 관례 정립 |
| D4. 프론트엔드 인증 | 높음 | 중간 | STML codegen 확장 |
