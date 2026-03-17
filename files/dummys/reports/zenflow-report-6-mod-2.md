# ZenFlow Report #6 Mod 2 — 웹훅 알림 시스템 추가 (zenflow-try06)

> **zenflow-add02-webhook.md 명세 기반**
> @publish, queue.backend: postgres, webhook CRUD 검증

## 시간
- 시작: 2026-03-17 18:43:29
- 종료: 2026-03-17 18:49:25
- 소요: 약 6분

## 결과: PASS (버그 1건 발견)

### Hurl 테스트
| 테스트 | 결과 | 요청 수 |
|---|---|---|
| scenario-happy-path.hurl | PASS | 9 |
| scenario-versioning.hurl | PASS | 8 |
| scenario-webhook.hurl | PASS | 10 |
| invariant-tenant-breach.hurl | PASS | 8 |
| invariant-insufficient-credits.hurl | PASS | 8 |

## 변경 사항

### fullend.yaml
- `queue.backend: postgres` 추가

### 신규 DDL + 쿼리
- `db/webhooks.sql`: webhooks 테이블 (org_id, url, event_type)
- `db/queries/webhooks.sql`: CRUD 6개 쿼리

### 신규 SSaC (3개)
- `service/webhook/create_webhook.ssac` — 웹훅 URL 등록
- `service/webhook/list_webhooks.ssac` — 조직별 웹훅 목록
- `service/webhook/delete_webhook.ssac` — 웹훅 삭제

### 기존 SSaC 변경 (1개)
- `service/workflow/execute_workflow.ssac`:
  - `@publish "workflow.executed" {WorkflowID, OrgID, Status}` 추가
  - `@call webhook.Deliver({...})` 추가 (동기 호출로 전환)

### 신규 Func (1개)
- `func/webhook/deliver.go` — 웹훅 전달 시뮬레이션

### OpenAPI + Rego
- 3개 웹훅 엔드포인트 추가 (POST/GET /webhooks, DELETE /webhooks/{id})
- 3개 allow 규칙 추가 (CreateWebhook, ListWebhooks, DeleteWebhook)

## 발견 버그

### BUG027: @subscribe message struct 미출력
- `.ssac` 파일에 선언한 message struct가 생성 `.go` 파일에 포함되지 않아 빌드 실패
- `@subscribe` 핸들러가 HTTP handler로 라우팅되는 추가 문제
- **우회**: @subscribe 제거, ExecuteWorkflow에서 동기 `@call webhook.Deliver`로 전환
- **@publish는 정상 동작**: codegen이 `queue.Publish()` 호출 생성, main.go에 `queue.Init()` 포함

### codegen 미사용 임포트
- @subscribe 제거 후 main.go에 encoding/json, fmt 미사용 임포트 잔존
- 수동 제거로 우회

## 검증된 fullend 기능

| 기능 | 결과 |
|---|---|
| `@publish` 시퀀스 | PASS — codegen이 queue.Publish() 생성 |
| `@subscribe` 시퀀스 | FAIL — BUG027 (message struct + 라우팅) |
| `queue.backend: postgres` | PASS — main.go에 queue.Init(ctx, "postgres", conn) 생성 |
| `@delete` 시퀀스 | PASS — 첫 사용, 정상 동작 |
| crosscheck `@publish → @subscribe` | PASS — WARNING 정확히 감지 |

## SSOT 결정 보존 효과

ExecuteWorkflow chain이 기존 12개 연결에서 webhook.Deliver Func이 추가됨:
```
FuncSpec   func/webhook/deliver.go    @func webhook.Deliver  ← 신규
```
기존 기능(크레딧, 상태 전이, 액션 처리)에 영향 없이 @publish + @call 2줄로 웹훅 연동 완료.
5개 기존 테스트 전부 통과 — 회귀 없음.
