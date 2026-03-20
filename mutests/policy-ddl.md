# Mutation Test — Policy ↔ DDL

### MUT-POLICY-DDL-001: Policy ownership 컬럼 변경
- 대상: `specs/gigbridge/policy/authz.rego`
- 변경: `@ownership gig gigs.client_id` → `@ownership gig gigs.owner_id`
- 기대: ERROR — DDL gigs 테이블에 owner_id 컬럼 부재
- 결과: PASS — 재테스트 PASS (초회 sed 패턴 오류)

### MUT-POLICY-DDL-002: @ownership via join table 미존재
- 대상: `specs/gigbridge/policy/authz.rego`
- 변경: `@ownership gig gigs.client_id` → `@ownership gig gig_members.user_id via gig_members` (join table 추가, DDL에 미존재)
- 기대: ERROR — DDL에 "gig_members" join 테이블 없음
- 결과: 미실행

### MUT-POLICY-DDL-003: @ownership join column 미존재
- 대상: `specs/gigbridge/policy/authz.rego`
- 변경: join table 어노테이션에 존재하지 않는 join column 지정
- 기대: ERROR — join 테이블에 해당 컬럼 없음
- 결과: 미실행
