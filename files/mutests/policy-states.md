# Mutation Test — Policy ↔ States

### MUT-POLICY-STATES-001: @ownership 테이블이 DDL에 없음
- 대상: `specs/gigbridge/policy/authz.rego`
- 변경: `@ownership gigs.user_id` → `@ownership nonexistent.user_id`
- 기대: ERROR — @ownership 테이블 "nonexistent"가 DDL에 없음
- 결과: PASS — crosscheck/policy에서 @ownership 테이블 존재 검증

### MUT-POLICY-STATES-002: @ownership 컬럼이 DDL 테이블에 없음
- 대상: `specs/gigbridge/policy/authz.rego`
- 변경: `@ownership gigs.user_id` → `@ownership gigs.nonexistent_id`
- 기대: ERROR — gigs 테이블에 "nonexistent_id" 컬럼 없음
- 결과: PASS — crosscheck/policy에서 @ownership 컬럼 존재 검증
