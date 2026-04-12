# pkg/generate

Toulmin Trace 기반 코드 생성 엔진.
**외부 도구가 할 수 있는 것은 외부 도구에 위임**하고, fullend 고유 로직(SSaC 시퀀스, STML 바인딩, Hurl 시나리오 순서)만 toulmin으로 구현한다.

## 재발명 금지 원칙

| 대상 | 도구 | 역할 |
|------|------|------|
| OpenAPI → Go types | **oapi-codegen** (외부) | `types.gen.go` — 요청/응답 구조체 |
| OpenAPI → HTTP server skeleton | **oapi-codegen** (외부) | `server.gen.go` — 라우팅, 파라미터 추출 |
| DDL + queries → Go DB 접근 코드 | **sqlc** (외부) | `queries.sql.go` — 타입 안전 쿼리 함수 |
| SSaC 시퀀스 → 핸들러 body | **pkg/generate/backend** | 시퀀스 로직 조합 (@get/@post/@auth/@state/@call/...) |
| STML → React 페이지 | **pkg/generate/frontend** | data-* 바인딩 → useQuery/useMutation |
| OpenAPI + 상태 → Hurl 시나리오 | **pkg/generate/hurl** | scenario ordering + request 생성 |

## 산출물 체인

```
fullend gen <specs-dir> <artifacts-dir>
  ↓
1. oapi-codegen 실행 (외부)
   → artifacts/backend/internal/api/types.gen.go
   → artifacts/backend/internal/api/server.gen.go

2. sqlc 실행 (외부, DDL + queries)
   → artifacts/backend/internal/db/*.sql.go

3. pkg/generate/backend (toulmin)
   → artifacts/backend/internal/service/**/*.go        — 핸들러 body
   → artifacts/backend/internal/model/**/*.go          — 모델 인터페이스
   → artifacts/backend/cmd/main.go                     — 초기화 + 라우터 연결
   → artifacts/backend/internal/middleware/*.go        — auth, claims
   → artifacts/backend/internal/authz/*.go             — Rego 연동

4. pkg/generate/frontend (정적 + toulmin)
   → artifacts/frontend/src/{App,main}.tsx, api/client.ts
   → artifacts/frontend/src/pages/*.tsx                — STML 바인딩

5. pkg/generate/hurl (toulmin + topological sort)
   → artifacts/tests/smoke-*.hurl                      — 자동 생성
   → artifacts/tests/scenario-*.hurl                   — 사용자 작성(생략 가능)

6. Contract 디렉티브 삽입
   → 생성된 Go/TSX 파일에 //fullend:gen ssot=... contract=... 주입
```

## 아키텍처

```
Fullstack (파싱 결과)
  → ground.Build(fs) → *rule.Ground
  → SSOT별 AST 노드 순회:
      seq[0] → Evaluate(ctx, Trace: true) → Trace 패턴 → 코드 생성기 실행
      seq[1] → Evaluate(ctx, Trace: true) → Trace 패턴 → 코드 생성기 실행
      ...
  → 생성된 코드 조립 → oapi-codegen/sqlc 산출물과 결합 → 파일 출력
```

## Trace 패턴 매칭

각 AST 노드(시퀀스, 페이지 블록 등)에 대해 Toulmin Graph를 평가하면,
모든 warrant의 activated/deactivated 상태가 Trace로 반환된다.

```go
results, _ := graph.Evaluate(ctx, toulmin.EvalOption{Trace: true})
// Trace = [
//   {Name: "IsGet",       Activated: true},
//   {Name: "HasFK",       Activated: true},
//   {Name: "HasPaginate", Activated: false},
//   {Name: "IsSubscribe", Activated: false},
// ]
// → 패턴: Get + FK → "FK 참조 Get + @empty 가드" 코드 생성
```

이 Trace 패턴은 **심볼릭 뉴럴 네트워크**처럼 동작한다:
- 각 warrant = 뉴런 (activated/deactivated)
- defeat edge = 억제 연결
- Trace 패턴 = 활성화 벡터
- 코드 생성기 = 출력 레이어

## if-else 대비 장점

| | if-else 코드젠 | Toulmin 코드젠 |
|---|---|---|
| 조건 표현 | 중첩 분기, 순서 의존 | 선언적 규칙, 관계 정의 |
| 규칙 추가 | 기존 코드 수정 필요 | 새 warrant 등록, defeat edge 추가 |
| 가독성 | 분기 깊이에 비례하여 하락 | 규칙 수에 무관하게 유지 |
| 추적성 | 디버그 로그 수동 삽입 | Trace 자동 제공 |
| 테스트 | 모든 분기 조합 수동 작성 | RunCases + Trace 패턴 검증 |

## Backend 생성 상세

oapi-codegen과 sqlc가 생성한 산출물을 **전제**로 SSaC 시퀀스 → Go 핸들러 body를 생성한다.

### 핸들러 모드

기존 internal/gen의 이중 모드를 유지한다:

| 모드 | 조건 | 구조 |
|------|------|------|
| **Flat** | 모든 ServiceFunc의 Domain이 비어있음 | 단일 `service/Server` 구조체 |
| **Domain** | ServiceFunc 중 하나라도 Domain 필드 존재 | `service/{domain}/Handler` + 중앙 `service/Server` |

모드 선택은 warrant로 표현:
```go
hasDomains := g.Rule(HasDomainFuncs)
// Trace에 HasDomainFuncs: true → domain mode dispatcher 실행
```

### SSaC 시퀀스 → 코드 생성 Graph

| Graph | warrant 예시 | 생성 결과 |
|-------|-------------|----------|
| `get-codegen` | IsSimpleGet, IsFKGet, IsPaginatedOffset, IsPaginatedCursor, HasSort, HasFilter | 모델 조회 호출 + 스캔 |
| `mutation-codegen` | IsPost, IsPut, IsDelete, HasResult | sqlc exec/scan 호출 |
| `guard-codegen` | HasEmpty, HasExists, HasState, HasAuth, HasErrStatus | 가드 체크 + 조기 반환 |
| `call-codegen` | IsCall, HasResult, IsBuiltinAuth, IsCustomFunc | 함수 호출 + result 바인딩 |
| `response-codegen` | IsExplicit, IsShorthand, IsPaginatedResult | JSON 응답 구성 |
| `publish-codegen` | HasPublish, HasDelay, HasPriority | queue.Publish 호출 |
| `subscribe-codegen` | IsSubscribe, HasMessageStruct | queue.Subscribe 핸들러 |
| `transaction-codegen` | NeedsTransaction (DB 쓰기 포함) | sql.Tx 감싸기 |

### Defeat 활용 예시

```go
g := toulmin.NewGraph("get-codegen")

simple    := g.Rule(IsSimpleGet)
fk        := g.Rule(IsFKGet)
paginated := g.Rule(IsPaginatedOffset)
cursor    := g.Rule(IsPaginatedCursor)

fk.Attacks(simple)
paginated.Attacks(simple)
paginated.Attacks(fk)
cursor.Attacks(paginated)  // cursor가 offset pagination을 대체
```

### 모델 인터페이스 생성

DDL + sqlc 산출물로부터 모델 인터페이스를 파생한다:

| Graph | warrant | 생성 결과 |
|-------|---------|----------|
| `model-iface` | HasCreate, HasFindByID, HasList, HasUpdate, HasDelete | 인터페이스 메서드 시그니처 |
| `model-impl` | UsesDB, UsesCache | sqlc 쿼리 호출 래퍼 |

### 초기화 생성

| Graph | warrant | 생성 결과 |
|-------|---------|----------|
| `main-init` | HasSession, HasCache, HasFile, HasQueue, HasAuth | `cmd/main.go`의 Init 블록 |
| `middleware-init` | HasBearerAuth, HasClaims | `middleware/auth.go` 주입 |

## Frontend 생성 상세

기존 internal/gen/react는 **대부분 정적 파일 생성** (package.json, vite.config.ts 등)이다.
STML 바인딩 → React 컴포넌트 부분만 toulmin으로 구현.

### 정적 파일 (toulmin 불필요)

- package.json, vite.config.ts, tsconfig.json, index.html, main.tsx, App.tsx
- STML 의존성 및 pageOps 매핑에서 라우트 자동 구성

### 페이지 컴포넌트 (toulmin 적용)

| Graph | warrant | 생성 결과 |
|-------|---------|----------|
| `fetch-codegen` | HasFetch, HasPaginate, HasSort, HasFilter | `useQuery` 훅 + 상태 관리 |
| `action-codegen` | HasAction, HasFields, HasNavigate | `useMutation` 훅 |
| `bind-codegen` | IsBind, IsEach, IsState | JSX 바인딩 표현식 |
| `nested-codegen` | HasNestedFetch, HasChildAction | 중첩 블록 재귀 생성 |

## Hurl 생성 상세

기존 internal/gen/hurl의 가장 큰 가치는 **시나리오 순서 결정**이다.
toulmin이 가장 유리한 영역 — 우선순위 그래프가 자연스러움.

### 시나리오 순서 (5-phase)

기존 로직을 유지하되 toulmin으로 선언:

```
1. Auth          — Register → Login (고정 순서)
2. Creates + Transitions (interleaved)
                 — top-level creates (path depth ≤ 2) 먼저
                 — state transitions BFS 순서
                 — nested creates (depth > 2) 부모 transition 후
3. Updates       — PUT without @state
4. Reads         — GET (가장 파괴적이지 않은 순서)
5. Deletes       — FK 의존성 topological sort
```

### toulmin warrant 설계

| warrant | 역할 |
|---------|------|
| IsAuthStep | Auth phase 우선 |
| IsTopLevelCreate | phase 2 앞부분 |
| IsStateTransition | phase 2 중간 |
| IsNestedCreate | 부모 transition 이후 |
| IsUpdate | phase 3 |
| IsRead | phase 4 |
| IsDelete | phase 5 |
| HasFKDependency | 다른 delete 앞에 와야 함 (defeat으로 순서 조정) |

### Step 생성

| Graph | warrant | 생성 결과 |
|-------|---------|----------|
| `path-params` | HasPathParam, HasCapturedVar | `{{var}}` 치환 |
| `request-body` | HasRequestBody, HasEnumField, HasRequiredField | JSON body 생성 |
| `auth-token` | RequiresRole, HasCapturedToken | Bearer 토큰 삽입 |
| `captures` | ReturnsID, ReturnsObject | `[Captures]` 블록 |
| `asserts` | HasResponseSchema | `[Asserts]` 블록 |

### 상태 관리 (captures)

기존 구조 유지 — `captures map[string]bool`로 변수 추적, path param 해소 가능한 step만 실행.
toulmin warrant는 `HasCapturedVar(varName)` 형태로 동적 체크.

## 패키지 구조 (계획)

```
pkg/generate/
├── generate.go                  — Generate(fs, specsDir, artifactsDir) 엔트리포인트
├── run_external.go              — oapi-codegen, sqlc 호출
├── backend/
│   ├── generate.go              — 모드 선택 + phase 실행
│   ├── graph_get.go             — Get 코드젠 Graph 구성
│   ├── graph_mutation.go        — Post/Put/Delete
│   ├── graph_guard.go           — Empty/Exists/State/Auth
│   ├── graph_call.go            — @call
│   ├── graph_response.go        — @response
│   ├── graph_transaction.go     — Tx 감싸기
│   ├── graph_publish.go         — @publish
│   ├── graph_subscribe.go       — @subscribe
│   ├── graph_model_iface.go     — 모델 인터페이스
│   ├── graph_main_init.go       — main.go 초기화
│   ├── emit_*.go                — warrant별 코드 생성기
│   ├── assemble_handler.go      — 시퀀스 조립
│   ├── assemble_server.go       — Server 구조체 + 라우터
│   └── domain/
│       └── (domain mode 전용)
├── frontend/
│   ├── generate.go              — 정적 파일 + 페이지 생성
│   ├── static_files.go          — package.json, vite.config 등
│   ├── graph_fetch.go
│   ├── graph_action.go
│   ├── graph_bind.go
│   └── emit_*.go
├── hurl/
│   ├── generate.go              — 5-phase 오케스트레이션
│   ├── scenario_order.go        — topological sort
│   ├── graph_phase.go           — step phase 결정
│   ├── graph_step.go            — 단일 step 생성
│   ├── graph_captures.go        — 변수 추적
│   └── emit_*.go
├── contract/
│   └── inject_directive.go      — //fullend:gen 주입
└── trace/
    ├── pattern.go               — Trace 패턴 매칭 유틸
    └── dispatch.go              — 패턴 → 생성기 디스패치
```

## Contract 연동

생성된 모든 함수에 `//fullend:gen` 디렉티브를 삽입한다.
개발자가 `gen` → `preserve`로 변경하면 해당 함수 body는 재생성 시 보존된다.

```go
//fullend:gen ssot=service/gig/create_gig.ssac contract=a3f8c1
func (h *Handler) CreateGig(c *gin.Context) {
    // generated code...
}
```

## Trace 디버그

`fullend gen --trace` 플래그로 Trace 출력을 활성화하면,
각 노드에서 어떤 warrant가 발동했고 어떤 코드 생성기가 선택됐는지 확인할 수 있다.

```
[CreateGig] seq[0] @post
  ✓ IsPost (activated)
  ✗ IsGet (not activated)
  ✓ HasModel (activated)
  ✓ HasResult (activated)
  → emit: PostWithResult

[AcceptProposal] seq[11] @response
  ✓ IsExplicit (activated)
  ✗ IsShorthand (not activated)
  ✓ HasPagination (activated)
  → emit: PaginatedExplicitResponse
```

## 외부 도구 의존성

| 도구 | 설치 | 호출 시점 |
|------|------|----------|
| oapi-codegen | `go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest` | backend 생성 전 |
| sqlc | `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` | backend 생성 전 |

외부 도구 누락 시 `fullend gen`이 명확한 에러 메시지로 중단된다.
