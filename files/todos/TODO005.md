# TODO005: fullend validate 개선 목록

## Part A: 수동 분석으로 발견한 미탐지 이슈

gigbridge specs 전체 분석 시 수동으로 발견했으나 `fullend validate`가 탐지하지 못한 항목.

## 1. OPA RaiseDispute 코멘트 ↔ 실제 동작 불일치

- **파일**: `specs/gigbridge/policy/authz.rego`
- **내용**: 코멘트는 "Either client (owner) or freelancer (assignee)"이지만 실제로는 `resource == "gig"` (= `gigs.client_id`)만 체크하므로 client only
- **SSaC**: `raise_dispute.ssac`에서 `ResourceID: gig.ClientID`만 전달
- **수정 방안**: (A) 코멘트를 "Owner only"로 수정, 또는 (B) freelancer도 분쟁 가능하게 `gig_assignee` resource 두 번째 rule 추가
- **검증 자동화**: OPA 코멘트 파싱 → 실제 rule 조건과 비교하는 crosscheck 규칙 추가 가능

## 2. OpenAPI User 스키마에 password_hash 노출

- **파일**: `specs/gigbridge/api/openapi.yaml` — `components/schemas/User`
- **내용**: API 스키마 정의에 `password_hash` 필드 포함. 현재 응답에서 직접 `$ref: User`를 쓰지 않지만, API 문서/클라이언트 코드젠에서 노출 위험
- **수정 방안**: `password_hash` 제거 또는 `writeOnly: true` 표시
- **검증 자동화**: OpenAPI 스키마에 `password`, `secret`, `hash` 등 민감 필드명 패턴 감지 규칙 추가 가능

## 3. Scenario 테스트 커버리지 갭

- **파일**: `specs/gigbridge/scenario/gig_lifecycle.feature`
- **내용**: 12개 operation 중 RejectProposal, RaiseDispute 2개가 시나리오에서 미테스트
- **수정 방안**: RejectProposal (거절 후 status=="rejected"), RaiseDispute (under_review → disputed) 시나리오 추가
- **검증 자동화**: Scenario에서 사용된 operationId 집합 vs OpenAPI operationId 집합 diff → 미커버 operation WARN

## 4. model/dto.go Token DTO 미사용

- **파일**: `specs/gigbridge/model/dto.go`
- **내용**: `Token` struct 정의되어 있으나 어떤 SSaC에서도 참조하지 않음. Login은 `auth.IssueTokenResponse`를 직접 사용
- **수정 방안**: 삭제 또는 실제 사용처 연결
- **검증 자동화**: DTO 정의 vs SSaC/Func 참조 diff → 미사용 DTO WARN

## 5. Frontend STML 불완전 — 누락된 액션

- **파일**: `specs/gigbridge/frontend/gig-detail.html`
- **내용**: Publish, SubmitWork, ApproveWork 버튼만 있고 AcceptProposal, RejectProposal, RaiseDispute 버튼 없음. Proposal 목록 표시 UI도 없음
- **수정 방안**: proposal 목록 섹션 + accept/reject 버튼 + dispute 버튼 추가
- **검증 자동화**: STML `data-action` 집합 vs OpenAPI security 필요 operationId 집합 diff → 미구현 액션 WARN

---

## Part B: fullend validate 신뢰도 개선

validate의 약속 = "specs만 잘 쓰면 나머지는 fullend가 보장한다". 현재 약 80%. 아래 해결 시 100%에 근접.

### B-1. False positive 정리 — 노이즈가 많으면 WARN을 무시하게 됨

현재 validate 실행 시 27건 WARN이 출력되지만 전부 false positive.
사용자가 WARN을 무시하는 습관이 들면 진짜 문제도 놓치게 된다.
false positive를 0으로 만들어야 WARN의 신뢰도가 올라간다.

### B-2. 의미적 검증 부재 — 구조는 맞지만 의미가 틀린 케이스

fullend는 "이 필드가 존재하는가", "이 operation이 정의되었는가" 같은 구조적 검증은 잘 하지만:
- OPA 코멘트와 실제 rule 조건의 불일치 (Part A #1)
- 민감 필드(password_hash) 노출 위험 (Part A #2)
- 이런 "구조는 맞지만 의미가 틀린" 케이스는 탐지하지 못함

구현 방안:
- OPA 코멘트 내 role/ownership 키워드 파싱 → rule 조건과 비교
- OpenAPI 스키마 필드명에 `password`, `secret`, `hash`, `token`, `key` 패턴 감지
- SSaC `@response`에서 민감 필드가 포함된 모델을 직접 반환하는지 검사

### B-3. 시나리오 커버리지 측정 — 어떤 operation이 테스트 안 되는지 알려줘야 함

현재 validate는 시나리오의 문법과 참조 유효성만 검증.
OpenAPI에 정의된 12개 operation 중 몇 개가 시나리오에서 커버되는지 리포트하지 않음.
`fullend status`에 시나리오 커버리지 퍼센트를 표시하면 갭이 즉시 보인다.

구현 방안:
- Scenario operationId 집합 vs OpenAPI operationId 집합 diff
- `fullend status` 출력에 `Scenario coverage: 10/12 (83%)` 추가
- 미커버 operation 목록 표시

---

## Part C: fullend validate false positive 상세

### SSaC ListGigs offset/limit WARN (false positive)

- `list_gigs.ssac`에서 `@get Page[Gig] gigPage = Gig.List({Query: query})` 패턴 사용
- `{Query: query}`의 `query`는 SSaC 내장 pagination 변수로, offset/limit/sort/filter를 통째로 전달
- 생성된 코드에서 `QueryOpts`로 변환되어 offset/limit 정상 동작 확인됨
- SSaC 검증기가 `request.offset`, `request.limit` 명시적 참조가 없다고 WARN 출력
- **수정 방안**: `Page[T]` + `{Query: query}` 패턴일 때 offset/limit WARN 억제

### DDL → OpenAPI 25건 WARN (false positive)

- `components/schemas`에 모든 컬럼이 정의되어 있는데 "OpenAPI 스키마에 없습니다" 출력
- 원인 추정: crosscheck가 component schemas를 인식하지 못하고 inline response schema만 탐색하며, 중첩 응답 (`{ gig: { ... } }`) 내부를 탐색하지 못함
- fullend crosscheck 로직 개선 필요
