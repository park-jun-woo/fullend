# ZenFlow Add-on #02 — 웹훅 알림 시스템

## 개요
워크플로우 실행 완료 시 `@publish`로 이벤트 발행, `@subscribe`로 웹훅 URL에 알림 전송. fullend의 pub/sub 파이프라인 검증.

## fullend.yaml 변경
- `queue.backend: postgres` 추가

## 신규 엔드포인트
- **POST /webhooks** (`CreateWebhook`): 조직별 웹훅 URL 등록
- **GET /webhooks** (`ListWebhooks`): 조직의 웹훅 목록 조회
- **DELETE /webhooks/{id}** (`DeleteWebhook`): 웹훅 삭제

## DDL 추가
- `webhooks` 테이블: id, org_id (FK), url, event_type, created_at

## SSaC 변경
- `execute_workflow.ssac`: 기존 @response 전에 `@publish "workflow.executed" {WorkflowID: wf.ID, OrgID: wf.OrgID, Status: "completed"}` 추가
- `on_workflow_executed.ssac` (신규): `@subscribe "workflow.executed"` → 해당 org의 웹훅 URL 조회 → `@call webhook.Deliver({URL: ..., Payload: ...})` 호출

## Custom Functions
- `webhook.Deliver(URL, Payload)`: HTTP POST 시뮬레이션 (Func purity → 실제 전송 불가, 시뮬레이션만)

## 검증 포인트
- `@publish` / `@subscribe` 시퀀스 타입 (11개 중 유일 미검증)
- `queue.backend: postgres` 설정 (fullend_queue 테이블)
- 비동기 이벤트 처리 패턴
- crosscheck: `@publish topic → @subscribe exists` WARNING 규칙

## E2E Scenario
- 조직 생성 → 웹훅 등록 → 워크플로우 생성/활성화/실행 → 실행 로그 확인 + 웹훅 전달 확인

## 참고
- `@subscribe` 함수는 HTTP 트리거가 아닌 큐 이벤트 트리거. `@response` 사용 불가.
- message struct를 .ssac 파일 내에 선언해야 함.
