# ZenFlow Add-on #03 — 워크플로우 템플릿 마켓플레이스

## 개요
조직이 워크플로우를 "템플릿"으로 공개. 다른 조직이 검색/복제. cursor pagination, @exists, @dto 모델 검증.

## 신규 엔드포인트
- **POST /templates** (`PublishTemplate`): 워크플로우를 템플릿으로 공개
- **GET /templates** (`ListTemplates`): cursor pagination + x-filter로 템플릿 검색
- **POST /templates/{id}/clone** (`CloneTemplate`): 템플릿을 자기 조직 워크플로우로 복제
- **GET /templates/{id}** (`GetTemplate`): 템플릿 상세 조회

## DDL 추가
- `templates` 테이블: id, source_workflow_id (FK), org_id (FK), title, description, category, clone_count (INT DEFAULT 0), created_at
- `CREATE UNIQUE INDEX idx_templates_source ON templates(source_workflow_id)` — 중복 공개 방지

## @dto 모델
- `TemplateDetail`: 템플릿 + 작성 조직명 + 액션 수 등 조합 정보 (DDL에 없는 비정규화 뷰)

## OpenAPI x- 확장
- `x-pagination: { style: cursor, defaultLimit: 20, maxLimit: 100 }` — cursor pagination
- `x-filter: { allowed: [category] }` — 카테고리 필터
- cursor는 id DESC 고정 (x-sort 없음)

## SSaC 설계
- PublishTemplate: `@get Workflow` → `@exists Template "Already published" 409` → `@post Template`
- CloneTemplate: `@get Template` → 대상 조직에 워크플로우 복제 + 액션 복사 + clone_count 증가

## 검증 포인트
- **cursor pagination**: Page[T] 대신 Cursor[T] 사용, LIMIT+1 방식
- **@exists 가드**: 중복 등록 방지 (nil이 아니면 409)
- **@dto 모델**: DDL 테이블과 무관한 순수 DTO 타입
- **OpenAPI response**: items + next_cursor + has_next 형식

## E2E Scenario
- Org A: 워크플로우 생성 → 템플릿 공개 → 중복 공개 시도 → 409
- Org B: 템플릿 목록 조회 (cursor pagination) → 템플릿 복제 → 자기 워크플로우로 확인
