# Phase 027: Scenario Crosscheck 강화 ✅ 완료

## 목표

기존 `CheckScenarios`에 누락된 3가지 검증을 추가한다:

1. **캡처 참조 유효성** — 캡처 변수가 사용되기 전에 선언되었는지
2. **토큰 흐름 / role 매칭** — OPA policy 기반으로 올바른 role의 토큰이 사용되는지
3. **status code 유효성** — assertion의 status code가 OpenAPI responses에 정의되어 있는지

## 검증 규칙 상세

### Rule 4: 캡처 참조 유효성 (Scenario 내부)

시나리오 내에서 `-> varName`으로 캡처한 변수가 이후 단계의 JSON `{...}` 안에서 `varName.field` 형태로 참조될 때, 해당 `varName`이 이전에 캡처된 적이 있는지 확인.

```
When POST Login {...} -> clientToken     # clientToken 캡처
When POST CreateGig {...} -> gigResult   # gigResult 캡처
When PUT PublishGig {"id": gigResult.gig.id}  # gigResult 참조 — OK
When PUT PublishGig {"id": unknown.gig.id}    # unknown 미캡처 — ERROR
```

- 레벨: ERROR
- 참조 SSOT: feature 파일 내부 (외부 SSOT 불필요)
- 구현: 시나리오별로 캡처된 변수명 set 관리, JSON 내 `varName.xxx` 패턴 추출 후 set 검사

### Rule 5: 토큰 흐름 / role 매칭 (Scenario ↔ OPA Policy)

Login 캡처 직전의 Register에서 지정한 role이, 이후 호출하는 operation의 OPA policy에서 허용하는 role과 일치하는지 확인.

```
Given POST Register {"role": "client"}
When POST Login {...} -> clientToken       # clientToken은 client role
When POST SubmitProposal {...}             # OPA: freelancer만 허용 → WARNING
```

- 레벨: WARNING (시나리오가 의도적 거부 테스트일 수 있음, `@invariant` + 4xx 기대 시 스킵)
- 참조 SSOT: OPA Rego policy
- 구현:
  1. Register 단계에서 `role` 필드 추출 → 다음 Login 캡처에 연결
  2. 캡처명 → role 매핑 테이블 구성
  3. 각 action step에서 마지막 Login 캡처의 role과 OPA 허용 role 비교

### Rule 6: status code 유효성 (Scenario ↔ OpenAPI)

`Then status == CODE` assertion의 CODE가 해당 operation의 OpenAPI responses에 정의되어 있는지 확인.

```yaml
# OpenAPI
/gigs/{id}/publish:
  put:
    operationId: PublishGig
    responses:
      "200": ...
      "404": ...
```

```
When PUT PublishGig {...}
Then status == 200    # OK — OpenAPI에 200 정의됨
Then status == 500    # WARNING — OpenAPI에 500 미정의
```

- 레벨: WARNING (OpenAPI에 모든 에러 코드를 명시하지 않을 수 있음)
- 참조 SSOT: OpenAPI responses
- 구현: opMap에서 operation 가져온 후 `op.Responses` 키에 status code 존재 여부 확인

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/crosscheck/scenario.go` | Rule 4, 5, 6 구현 추가 |
| `internal/crosscheck/crosscheck.go` | `CheckScenarios` 호출에 `Policies` 전달 (시그니처 변경) |
| `internal/crosscheck/scenario_test.go` | 신규 — 3개 규칙 테스트 |

## 의존성

- `internal/scenario` — Feature/Step/Capture 타입 (변경 없음)
- `internal/policy` — OPA Policy 파서 (변경 없음, 기존 `RolesForOperation` 활용)
- OpenAPI `openapi3.T` — responses 접근 (변경 없음)

## 검증 방법

1. `go test ./internal/crosscheck/...` — 신규 테스트 통과
2. `go test ./...` — 전체 테스트 통과
3. GigBridge specs로 `fullend validate` 실행 — 새 규칙이 잘못된 시나리오를 탐지하는지 확인
