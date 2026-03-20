# fullend

## 벤치마크: ZenFlow (멀티테넌트 워크플로우 자동화 SaaS)

| 단계 | 내용 | 시간 | 누적 |
|------|------|-----:|-----:|
| 초기 구축 | 멀티테넌트, 크레딧 시스템, 상태머신, 5 tables, 12 endpoints | 20분 | 20분 |
| +버전 관리 | 워크플로우 복제, 버전 목록, INSERT...SELECT 액션 복사 | 12분 | 32분 |
| +웹훅 알림 | 이벤트 발행, 웹훅 CRUD, queue backend | 6분 | 38분 |
| +템플릿 마켓 | cursor pagination, cross-org 복제, public endpoints | 5분 | 43분 |
| +파일 첨부 | 실행 리포트 생성, file backend | 4분 | 47분 |

**최종: 23 endpoints, 7 tables, 23 services, 18 auth rules, 65 test requests. 전부 통과.**

기능을 추가해도 느려지지 않는다. 기존 테스트는 한 번도 깨지지 않았다.

---

풀스택 SSOT 오케스트레이터. 9개 선언 소스의 정합성을 검증하고, 코드를 생성한다.

### Quick Start

```bash
# 설치
go install github.com/park-jun-woo/fullend/cmd/fullend@latest

# 예제 프로젝트로 즉시 체험
git clone https://github.com/park-jun-woo/fullend && cd fullend
fullend validate examples/zenflow
```

```
✓ Config       zenflow, go/gin, typescript/react
✓ OpenAPI      23 endpoints
✓ DDL          7 tables, 43 columns
✓ SSaC         23 service functions
✓ Model        1 files
✓ STML         2 pages, 4 bindings
✓ States       1 diagrams, 7 transitions
✓ Policy       1 files, 18 rules
✓ Scenario     7 scenario hurl files
✓ Func         7 funcs
✓ Cross        0 mismatches

All SSOT sources are consistent.
```

```bash
fullend chain ExecuteWorkflow examples/zenflow
```

```
── Feature Chain: ExecuteWorkflow ──

  OpenAPI    api/openapi.yaml                POST /workflows/{id}/execute
  SSaC       service/workflow/execute_workflow.ssac   @get @empty @auth @state @call @publish @response
  DDL        db/workflows.sql                CREATE TABLE workflows
  DDL        db/execution_logs.sql           CREATE TABLE execution_logs
  Rego       policy/authz.rego               resource: workflow
  StateDiag  states/workflow.md              diagram: workflow → ExecuteWorkflow
  FuncSpec   func/billing/check_credits.go   @func billing.CheckCredits
  FuncSpec   func/billing/deduct_credit.go   @func billing.DeductCredit
  FuncSpec   func/worker/process_actions.go  @func worker.ProcessActions
  FuncSpec   func/webhook/deliver.go         @func webhook.Deliver
  Hurl       tests/scenario-happy-path.hurl  scenario: scenario-happy-path.hurl
```

### AI와 함께 쓰기

위 벤치마크는 AI 에이전트가 SSOT를 작성하고 fullend가 검증하는 방식으로 측정되었다. Claude Code, Codex, Copilot, Cursor — 어떤 AI를 쓰든 상관없다.

AI 에이전트를 켜고 다음 프롬프트를 입력한다:

```
fullend/manual-for-ai.md 확인하고 fullend/examples/zenflow/zenflow.md대로 개발하라.
```

AI가 SSOT 명세를 쓰면, `fullend validate`가 레이어 간 불일치를 즉시 잡아낸다. AI는 자유롭게 설계하되, 레일 밖으로 나가면 검증이 실패한다.

## 9 SSOT Sources

```
specs/
├── fullend.yaml             → 프로젝트 설정 (필수)
├── api/openapi.yaml         → OpenAPI 3.x
├── db/*.sql                 → SQL DDL + sqlc 쿼리
├── service/**/*.ssac        → SSaC (서비스 시퀀스 DSL)
├── model/*.go               → Go 구조체 (// @dto)
├── func/<pkg>/*.go          → 커스텀 함수 구현 (선택)
├── states/*.md              → Mermaid stateDiagram (상태 전이)
├── policy/*.rego            → OPA Rego (인가 정책)
├── tests/scenario-*.hurl    → Hurl 시나리오 테스트
├── tests/invariant-*.hurl   → Hurl 불변성 테스트
├── frontend/*.html          → STML (HTML5 + data-*)
```

## 왜 AI가 헤매지 않는가

AI 에이전트에게 "기능 추가해"라고 하면 보통 프로젝트가 커질수록 맥락을 잃는다. fullend는 9개 SSOT가 서로를 참조하고, `validate`가 불일치를 즉시 잡아낸다. AI는 자유롭게 코드를 쓰되, SSOT 바깥으로 벗어나면 검증이 실패한다. 레일 위의 자유.

## Commands

### validate

각 SSOT를 개별 검증한 뒤, 레이어 간 교차 검증을 수행한다.

```bash
fullend validate <specs-dir>
fullend validate --skip states <specs-dir>
```

```
✓ Config       my-project, go/gin, typescript/react
✓ OpenAPI      12 endpoints
✓ DDL          4 tables, 23 columns
✓ SSaC         12 service functions
✓ Model        1 files
✓ STML         2 pages, 2 bindings
✓ States       2 diagrams, 7 transitions
✓ Policy       1 files, 7 rules, 3 ownership mappings
✓ Scenario     3 scenario hurl files
✓ Func         2 funcs
✓ Cross        0 mismatches
— Contract     no artifacts

All SSOT sources are consistent.
```

### gen

검증 후 모든 SSOT에서 코드를 생성한다.

```bash
fullend gen <specs-dir> <artifacts-dir>
```

### chain

하나의 API 오퍼레이션에 연결된 모든 SSOT 노드를 추적한다.

```bash
fullend chain <operationId> <specs-dir>
```

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
  Hurl       tests/scenario-gig-lifecycle.hurl:4           scenario: scenario-gig-lifecycle.hurl
```

### gen-model

외부 OpenAPI 문서에서 Go 모델(인터페이스 + 타입 + HTTP 클라이언트)을 생성한다.

```bash
fullend gen-model <openapi-source> <output-dir>
fullend gen-model https://api.stripe.com/openapi.yaml ./external/
```

### status

감지된 SSOT 현황을 요약한다.

```bash
fullend status <specs-dir>
```

## Cross-Validation

개별 도구(SSaC, STML)는 자기 레이어만 검증한다. fullend는 레이어 **사이**의 불일치를 잡는다:

- **fullend.yaml ↔ OpenAPI** — 미들웨어 이름이 securitySchemes 키와 일치
- **OpenAPI x-sort/x-filter ↔ DDL** — 참조 컬럼이 테이블에 존재
- **OpenAPI x-include ↔ DDL** — 참조 리소스가 테이블에 매핑
- **SSaC @result ↔ DDL** — 결과 타입이 DDL 파생 모델과 일치
- **SSaC arg ↔ DDL** — 인자 필드명이 테이블 컬럼과 일치
- **States ↔ SSaC** — 전이 이벤트가 SSaC 함수와 매칭
- **States ↔ DDL** — 상태 필드가 DDL 컬럼에 매핑
- **States ↔ OpenAPI** — 전이 이벤트가 operationId와 매칭
- **Policy ↔ SSaC** — @auth (action, resource) 쌍이 Rego allow 규칙과 매칭
- **Policy ↔ DDL** — @ownership 테이블/컬럼 참조가 DDL에 존재
- **Policy ↔ States** — @auth가 있는 상태 전이에 Rego 규칙 존재
- **Hurl ↔ OpenAPI** — 테스트가 유효한 엔드포인트를 참조
- **Queue** — @publish 토픽과 @subscribe 함수 일치, 페이로드 필드 정합성
- **Func ↔ SSaC** — @call 참조에 구현이 존재, 인자 수/타입 일치
- **STML ↔ SSaC** (간접) — 동일한 OpenAPI operationId 참조

## Default Functions (pkg/)

SSaC `@call`에서 사용 가능한 내장 함수:

| 패키지 | 함수 | 설명 |
|---|---|---|
| `auth` | `hashPassword` | bcrypt 해싱 |
| `auth` | `verifyPassword` | bcrypt 검증 |
| `auth` | `issueToken` | JWT 액세스 토큰 (24h) |
| `auth` | `verifyToken` | JWT 검증 + 클레임 추출 |
| `auth` | `refreshToken` | 리프레시 토큰 (7일) |
| `auth` | `generateResetToken` | 비밀번호 재설정 토큰 |
| `crypto` | `encrypt` | AES-256-GCM 암호화 |
| `crypto` | `decrypt` | AES-256-GCM 복호화 |
| `crypto` | `generateOTP` | TOTP 비밀 + QR URL |
| `crypto` | `verifyOTP` | TOTP 코드 검증 |
| `storage` | `uploadFile` | S3 파일 업로드 |
| `storage` | `deleteFile` | S3 파일 삭제 |
| `storage` | `presignURL` | S3 presigned URL |
| `mail` | `sendEmail` | SMTP 이메일 |
| `mail` | `sendTemplateEmail` | 템플릿 HTML 이메일 |
| `text` | `generateSlug` | 유니코드 → URL 슬러그 |
| `text` | `sanitizeHTML` | XSS 방지 HTML 새니타이즈 |
| `text` | `truncateText` | 유니코드 텍스트 자르기 |
| `image` | `ogImage` | OG 이미지 생성 (1200×630) |
| `image` | `thumbnail` | 썸네일 생성 (200×200) |

`specs/<project>/func/<pkg>/`에 커스텀 구현을 두면 내장 함수를 오버라이드할 수 있다.

## Built-in Models (pkg/)

DDL 외 I/O를 위한 패키지 레벨 @model 인터페이스. `fullend.yaml`에서 설정.

| 패키지 | 인터페이스 | 백엔드 | SSaC 사용 |
|---|---|---|---|
| `session` | `SessionModel` (Set/Get/Delete + TTL) | PostgreSQL, Memory | `session.Session.Get({key: ...})` |
| `cache` | `CacheModel` (Set/Get/Delete + TTL) | PostgreSQL, Memory | `cache.Cache.Set({key: ..., value: ..., ttl: ...})` |
| `file` | `FileModel` (Upload/Download/Delete) | S3, LocalFile | `file.File.Upload({key: ..., body: ...})` |
| `queue` | Singleton Pub/Sub (Publish/Subscribe) | PostgreSQL, Memory | `@publish "topic" {payload}` |

## Runtime Testing

`fullend gen`은 OpenAPI 스펙에서 [Hurl](https://hurl.dev) 테스트를 생성한다.

```bash
hurl --test --variable host=http://localhost:8080 artifacts/my-project/tests/*.hurl
```

- **smoke.hurl** — 엔드포인트 스모크 테스트 (자동 생성)
- **scenario-*.hurl** — 비즈니스 시나리오 테스트 (직접 작성)
- **invariant-*.hurl** — 교차 엔드포인트 불변성 테스트 (직접 작성)

## Architecture

SSaC와 STML은 fullend에 `internal/ssac/`, `internal/stml/`로 통합되어 있다. [SSaC](https://github.com/park-jun-woo/ssac), [STML](https://github.com/park-jun-woo/stml) 리포는 fullend에서 복사한 미러.

모든 SSOT는 CLI 호출 당 한 번 `ParseAll()`로 파싱되어 validate, gen, status, chain 파이프라인에서 공유된다.

## Acknowledgments

fullend는 이 프로젝트들 위에 만들어졌다.

### SSOT Foundations

- [OpenAPI Initiative](https://www.openapis.org/) — 프론트엔드와 백엔드를 잇는 API 명세 표준
- [sqlc](https://sqlc.dev/) — SQL-first Go 코드 생성. fullend의 DDL 기반 모델 접근은 sqlc 철학에서 직접 영감
- [Open Policy Agent](https://www.openpolicyagent.org/) — 코드로서의 정책. Rego가 fullend의 인가 레이어를 구동
- [Mermaid](https://mermaid.js.org/) — 코드로서의 다이어그램. 상태 다이어그램이 런타임 상태 머신이 된다

### Code Generation & Validation

- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) — OpenAPI → Go 서버/타입 코드 생성
- [kin-openapi](https://github.com/getkin/kin-openapi) — Go에서 OpenAPI 3.x 파싱 및 검증
- [Hurl](https://hurl.dev/) — 플레인텍스트 HTTP 테스트

### Generated Code Runtime

- [React](https://react.dev/), [React Router](https://reactrouter.com/), [TanStack Query](https://tanstack.com/query), [React Hook Form](https://react-hook-form.com/)
- [Vite](https://vite.dev/), [Tailwind CSS](https://tailwindcss.com/), [TypeScript](https://www.typescriptlang.org/)
- [Gin](https://gin-gonic.com/), [lib/pq](https://github.com/lib/pq)

## License

MIT — [LICENSE](LICENSE) 참조.
