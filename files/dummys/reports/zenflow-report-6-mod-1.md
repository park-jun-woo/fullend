# ZenFlow Report #6 Mod 1 — 워크플로우 버전 관리 추가 (zenflow-try06)

> **gigbridge-try04 패턴 참고** + **zenflow-try06 기존 SSOT 위에 유지보수**
> SSOT 결정 보존 효과 검증 목적

## 시간
- 시작: 2026-03-17 18:20:49 (zenflow-try06 완료 직후)
- 종료: 2026-03-17 18:32:32
- 소요: 약 12분

## 결과: PASS

### Validate
- WARNING 5건 (기존과 동일, 무해)

### Build
- `go build` 성공

### Hurl 테스트
| 테스트 | 결과 | 요청 수 |
|---|---|---|
| smoke.hurl | PASS (서버 기동 시 자동) | — |
| scenario-happy-path.hurl | PASS | 9 |
| scenario-versioning.hurl | PASS | 8 |
| invariant-tenant-breach.hurl | PASS | 8 |
| invariant-insufficient-credits.hurl | PASS | 8 |

## 변경 사항

### 신규 파일 (5개)
| 파일 | 역할 |
|---|---|
| service/workflow/create_workflow_version.ssac | 기존 워크플로우 복제 + 버전 증가 + 액션 복사 |
| service/workflow/list_workflow_versions.ssac | root_workflow_id 기반 버전 목록 조회 |
| func/workflow/resolve_root_id.go | 조건분기: root_workflow_id가 0이면 자기 ID, 아니면 기존 root |
| func/workflow/next_version.go | 버전 번호 +1 계산 |
| tests/scenario-versioning.hurl | 버전 생성 + 액션 복사 + 버전 목록 검증 |

### 수정 파일 (5개)
| 파일 | 변경 내용 |
|---|---|
| db/workflows.sql | `version BIGINT DEFAULT 1`, `root_workflow_id BIGINT DEFAULT 0` 컬럼 추가 |
| db/queries/workflows.sql | `WorkflowCreateVersion`, `WorkflowListVersions` 쿼리 추가 |
| db/queries/actions.sql | `ActionCopyToWorkflow` (INSERT...SELECT) 쿼리 추가 |
| api/openapi.yaml | 2 엔드포인트 추가 + Workflow 스키마에 version, root_workflow_id 필드 |
| policy/authz.rego | CreateWorkflowVersion (admin), ListWorkflowVersions (any) 규칙 추가 |
| frontend/workflow-detail.html | version 바인딩 + New Version 버튼 추가 |

## Feature Chain 결과

```
── Feature Chain: CreateWorkflowVersion ──
  OpenAPI    POST /workflows/{id}/new-version
  SSaC       @get @empty @auth @call @post @put @response
  DDL        actions, workflows
  Rego       resource: workflow
  FuncSpec   workflow.ResolveRootID, workflow.NextVersion
  STML       data-action="CreateWorkflowVersion"

── Feature Chain: ListWorkflowVersions ──
  OpenAPI    GET /workflows/{id}/versions
  SSaC       @get @empty @auth @response
  DDL        workflows
  Rego       resource: workflow
```

## 설계 결정 기록

1. **root_workflow_id 패턴**: parent_id(직계 부모) 대신 root_workflow_id(최초 버전 ID)를 사용. 재귀 탐색 없이 `WHERE root_workflow_id = X OR id = X`로 전체 버전 조회 가능.

2. **Func으로 조건분기 위임**: SSaC가 `if`를 지원하지 않으므로, "root_workflow_id가 0이면 자기 ID 사용" 로직을 `resolveRootID` Func에 위임. SSaC의 `@call workflow.ResolveRootID` 라인이 이 결정의 근거를 보존.

3. **INSERT...SELECT로 액션 복사**: SSaC에 루프가 없으므로, sqlc 쿼리 `ActionCopyToWorkflow`가 DB 레벨에서 일괄 복사. SSaC에서는 `@put Action.CopyToWorkflow`로 한 줄.

4. **기존 테스트 무변경**: DDL 컬럼 추가(version, root_workflow_id)에 DEFAULT 값이 있으므로 기존 CreateWorkflow 흐름에 영향 없음. 4개 기존 테스트 전부 통과.

## SSOT 결정 보존 효과 평가

| 항목 | 결과 |
|---|---|
| 변경 영향 범위 파악 | `fullend chain`으로 즉시 파악 — 어떤 SSOT가 연결되는지 한눈에 |
| 정합성 검증 | `fullend validate`가 NOT NULL 누락, required 누락 등 즉시 감지 |
| 기존 기능 회귀 | 기존 4개 테스트 전부 통과 — DDL 변경이 기존에 영향 없음 확인 |
| 맥락 추적 | "왜 nextVersion이 별도 Func인가?" → SSaC @call 라인이 답 보존 |
| 변경 비용 | 5개 신규 + 6개 수정 = 총 11파일, 약 12분 소요 |

## 버그 리포트
없음.
