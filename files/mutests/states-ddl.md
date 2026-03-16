# Mutation Test — States ↔ DDL

### MUT-STATES-DDL-001: @state 상태 필드가 DDL에 없음
- 대상: `specs/gigbridge/db/gigs.sql`
- 변경: `status VARCHAR(50) NOT NULL DEFAULT 'draft'` 컬럼 삭제
- 기대: ERROR — stateDiagram의 상태 필드(status)가 DDL 테이블에 존재하지 않음
- 결과: PASS — crosscheck/states에서 @state Inputs 필드 → DDL 컬럼 매핑 검증

### MUT-STATES-DDL-002: DDL DEFAULT 값이 초기 상태와 불일치
- 대상: `specs/gigbridge/db/gigs.sql`
- 변경: `DEFAULT 'draft'` → `DEFAULT 'open'`
- 기대: WARNING — stateDiagram 초기 상태가 "draft"인데 DDL DEFAULT는 "open"
- 결과: SKIP — 현재 자동화된 검증 미구현 (향후 추가 예정)
