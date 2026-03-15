# Feature Chain

## 개념

Feature Chain은 하나의 API 기능(operationId)에 연결된 모든 SSOT 노드를 추출한 것이다.

풀스택 애플리케이션에서 하나의 기능은 여러 레이어에 걸쳐 있다. "제안 수락"이라는 기능 하나를 수정하려면 OpenAPI 스펙, 서비스 로직, DB 스키마, 인가 정책, 상태 전이, 외부 함수, 테스트 시나리오, 프론트엔드까지 모두 파악해야 한다. 기존에는 Grep을 수십 번 돌리거나 코드를 수동으로 추적해야 했다.

Feature Chain은 이 문제를 해결한다. operationId 하나를 입력하면, SSOT 간 심볼 참조를 따라가며 관련된 모든 파일과 라인 번호를 한 번에 보여준다.

## 왜 가능한가

fullend의 SSOT는 이미 서로를 심볼릭하게 참조하고 있다:

- SSaC의 `@get Model.Method` → DDL 테이블
- SSaC의 `@auth action resource` → Rego 정책
- SSaC의 `@state diagramID` → Mermaid stateDiagram
- SSaC의 `@call pkg.Func` → Func Spec
- OpenAPI의 operationId → SSaC 파일명
- Gherkin의 action step → operationId
- STML의 endpoint → OpenAPI path

crosscheck가 이미 이 참조들을 검증하고 있으므로, 같은 파싱 인프라를 재활용하면 그래프 탐색만으로 feature chain이 추출된다.

## 탐색 경로

```
operationId (시작점)
├── OpenAPI → path + method
├── SSaC → 서비스 함수 파일
│   ├── @get → DDL 테이블들
│   ├── @auth → Rego 정책 규칙
│   ├── @state → Mermaid stateDiagram 전이
│   ├── @call → Func Spec 구현체
│   └── @publish → 큐 구독자
├── Gherkin → operationId를 참조하는 시나리오들
└── STML → endpoint를 참조하는 프론트엔드 파일
```

## CLI

```bash
fullend chain <operationId> <specs-dir>
```

## 출력 예시

```
── Feature Chain: AcceptProposal ──

  OpenAPI    api/openapi.yaml:296                          POST /proposals/{id}/accept
  SSaC       service/proposal/accept_proposal.ssac:19      @get @empty @auth @state @put @call @post @response
  DDL        db/gigs.sql:1                                 CREATE TABLE gigs
  DDL        db/proposals.sql:1                            CREATE TABLE proposals
  DDL        db/transactions.sql:1                         CREATE TABLE transactions
  Rego       policy/authz.rego:3                           resource: gig
  StateDiag  states/gig.md:7                               diagram: gig → AcceptProposal
  StateDiag  states/proposal.md:6                          diagram: proposal → AcceptProposal
  FuncSpec   func/billing/hold_escrow.go:8                 @func billing.HoldEscrow
  Gherkin    scenario/gig_lifecycle.feature:4              Scenario: Happy Path - Full Gig Lifecycle
  Gherkin    scenario/gig_lifecycle.feature:42             Scenario: Unauthorized Access
```

연결되지 않는 SSOT 레이어는 출력하지 않는다.

## 활용

1. **수정 범위 파악** — 기능 하나를 변경할 때 어떤 파일들을 건드려야 하는지 즉시 파악
2. **AI 코드 수정** — AI에게 operationId 하나만 알려주면 전체 수정 범위를 자동으로 식별
3. **코드 리뷰** — PR에서 누락된 레이어가 있는지 chain과 대조
4. **온보딩** — 새 개발자가 기능 하나의 전체 구조를 한눈에 파악

## 미래: GEUL + SILK

모든 SSOT가 GEUL 그래프로 변환되면, feature chain은 SILK의 SIDX 비트와이즈 AND 쿼리가 된다. 현재는 파서별 탐색이지만, GEUL 통합 후에는 단일 인덱스 쿼리로 동일한 결과를 얻을 수 있다.
