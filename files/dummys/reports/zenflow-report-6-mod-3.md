# ZenFlow Report #6 Mod 3 — 템플릿 마켓플레이스 추가 (zenflow-try06)

> **zenflow-add03-template.md 명세 기반**
> cursor pagination, @exists(BUG028로 우회), cross-org 복제 검증

## 시간
- 시작: 2026-03-17 18:52:10
- 종료: 2026-03-17 18:57:37
- 소요: 약 5분

## 결과: PASS (버그 1건 발견)

### Hurl 테스트
| 테스트 | 결과 | 요청 수 |
|---|---|---|
| scenario-happy-path.hurl | PASS | 9 |
| scenario-versioning.hurl | PASS | 8 |
| scenario-webhook.hurl | PASS | 10 |
| scenario-template.hurl | PASS | 15 |
| invariant-tenant-breach.hurl | PASS | 8 |
| invariant-insufficient-credits.hurl | PASS | 8 |
| **합계** | **6/6 PASS** | **58 requests** |

## 변경 사항

### 신규 DDL + 쿼리
- `db/templates.sql`: templates 테이블 + UNIQUE INDEX on source_workflow_id
- `db/queries/templates.sql`: 5개 쿼리 (Create, FindByID, FindBySourceWorkflowID, List, IncrementCloneCount)

### 신규 SSaC (4개)
- `service/template/publish_template.ssac` — 워크플로우를 템플릿으로 공개
- `service/template/list_templates.ssac` — cursor pagination으로 템플릿 목록
- `service/template/get_template.ssac` — 템플릿 상세 + 조직명
- `service/template/clone_template.ssac` — 템플릿→워크플로우 복제 + 액션 복사 + clone_count 증가

### OpenAPI
- 4개 엔드포인트 추가
- `x-pagination: { style: cursor }` + `x-filter: { allowed: [category] }`
- ListTemplates, GetTemplate은 public (bearerAuth 불필요)

### Rego
- PublishTemplate: admin only
- CloneTemplate: any authenticated user

## 발견 버그

### BUG028: @empty 강제 규칙과 @exists 충돌
- FK 참조 @get 후 validator가 @empty를 강제 요구
- @exists는 "not nil이면 에러" — @empty와 논리적으로 모순
- **우회**: @exists 제거, DB UNIQUE 제약에 의존 (500 반환, 409가 아님)
- 테스트에서 HTTP 500으로 중복 감지 확인

## 검증된 fullend 기능

| 기능 | 결과 |
|---|---|
| **cursor pagination** | PASS — Cursor[Template], items/next_cursor/has_next 정상 |
| **x-filter** | PASS — category 필터 |
| **@exists 가드** | FAIL — BUG028 (@empty 강제 규칙 충돌) |
| **cross-org 복제** | PASS — Org A 템플릿 → Org B 워크플로우 복제 |
| **INSERT...SELECT 재활용** | PASS — 기존 ActionCopyToWorkflow 쿼리 재사용 |
| **public endpoint** | PASS — ListTemplates/GetTemplate에 bearerAuth 없이 동작 |

## 누적 zenflow-try06 현황

| 항목 | 값 |
|---|---|
| DDL 테이블 | 7개 |
| OpenAPI 엔드포인트 | 21개 |
| SSaC 서비스 함수 | 21개 |
| Func | 6개 |
| Rego 규칙 | 16개 |
| Hurl 테스트 | 6개, 58 requests |
| 발견 버그 | BUG027 (@subscribe), BUG028 (@exists) |
