# ZenFlow Add-on #06 — 감사 로그 (cache + offset pagination + x-sort)

## 개요
모든 주요 조작(워크플로우 생성/활성화/실행, 템플릿 공개/복제 등)에 감사 로그 기록. 최근 조회는 cache에서, 전체 이력은 DB에서 offset pagination + x-sort로 조회.

## 검증 포인트 (미검증 fullend 기능)
- **pkg/cache**: `cache.Cache.Set/Get` — 패키지 프리픽스 @model
- **cache.backend: postgres** 설정
- **offset pagination + x-sort 조합**: allowed 2개 이상, runtime sort switching
- **기존 SSaC에 @post 추가**: 여러 기존 서비스에 감사 로그 기록 1줄 삽입 — SSOT 변경 전파 테스트

## fullend.yaml 변경
- `cache.backend: postgres` 추가

## 신규 엔드포인트
- **GET /audit-logs** (`ListAuditLogs`): offset pagination + x-sort(created_at, action) + x-filter(action, actor_id)
- **GET /audit-logs/recent** (`GetRecentAuditLogs`): cache에서 최근 N건 조회 (빠른 대시보드용)

## DDL 추가
- `audit_logs` 테이블: id, org_id (FK), actor_id (FK users), action VARCHAR(100), resource_type VARCHAR(50), resource_id BIGINT, detail TEXT, created_at
- 인덱스: org_id, created_at, action

## SSaC 설계
- ListAuditLogs: `@get Page[AuditLog] page = AuditLog.ListByOrgID({OrgID: currentUser.OrgID, Query: query})` — offset pagination
- GetRecentAuditLogs: `@get cache.Cache.Get({key: cacheKey})` → cache hit이면 반환, miss면 DB 조회
- 기존 SSaC 변경 (3~5개 파일): `@post AuditLog.Create({...})` 1줄 추가
  - CreateWorkflow, ActivateWorkflow, ExecuteWorkflow, PublishTemplate, CloneTemplate 등

## OpenAPI x- 확장
```yaml
x-pagination:
  style: offset
  defaultLimit: 20
  maxLimit: 100
x-sort:
  allowed: [created_at, action]
  default: created_at
  direction: desc
x-filter:
  allowed: [action, actor_id]
```

## E2E Scenario
- 워크플로우 생성 → 활성화 → 실행 → 감사 로그 조회 (3건 이상) → sort by action 확인 → filter by action 확인
- 기존 테스트 회귀 확인 (기존 SSaC에 @post 추가해도 기존 response/status 무변경)
