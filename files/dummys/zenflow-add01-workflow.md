# ZenFlow Add-on #01 — 워크플로우 버전 관리

## 개요
기존 워크플로우를 복제해 새 draft 버전을 생성하고, 액션도 함께 복사. 버전 목록 조회 지원.

## 신규 엔드포인트
- **POST /workflows/{id}/new-version** (`CreateWorkflowVersion`): 기존 워크플로우 복제 + 버전 증가 + 액션 일괄 복사
- **GET /workflows/{id}/versions** (`ListWorkflowVersions`): root_workflow_id 기반 전체 버전 목록 조회

## DDL 변경
- `workflows` 테이블에 `version BIGINT NOT NULL DEFAULT 1`, `root_workflow_id BIGINT NOT NULL DEFAULT 0` 추가

## sqlc 쿼리 추가
- `WorkflowCreateVersion :one` — 버전/루트 지정 INSERT
- `WorkflowListVersions :many` — `WHERE (root_workflow_id = $1 OR id = $1) AND org_id = $2`
- `ActionCopyToWorkflow :exec` — `INSERT...SELECT`로 액션 일괄 복사

## Custom Functions
- `resolveRootID(WorkflowID, RootWorkflowID)`: 조건분기 — root_workflow_id가 0이면 자기 ID, 아니면 기존 root (SSaC에 if 없으므로 Func 위임)
- `nextVersion(CurrentVersion)`: 버전 번호 +1 계산

## 설계 결정
1. **root_workflow_id 패턴**: parent_id(직계 부모) 대신 root_workflow_id(최초 버전)를 사용. 재귀 탐색 없이 OR 조건으로 전체 버전 조회.
2. **INSERT...SELECT 액션 복사**: SSaC 루프 미지원 → DB 레벨 일괄 복사로 해결.
3. **DEFAULT 값 설정**: 기존 CreateWorkflow 흐름에 영향 없음 (version=1, root_workflow_id=0).

## Authorization
- CreateWorkflowVersion: admin only
- ListWorkflowVersions: any authenticated user (org 격리는 쿼리 레벨)

## E2E Scenario
- 워크플로우 생성 → 액션 2개 추가 → 새 버전 생성 → 버전 목록에 v1, v2 확인
