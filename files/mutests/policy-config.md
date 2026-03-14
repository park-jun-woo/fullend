# Mutation Test — Policy ↔ Config

### MUT-POLICY-CONFIG-001: Policy role명 변경
- 대상: `specs/gigbridge/policy/authz.rego`
- 변경: `input.role == "client"` → `input.role == "Client"` (PublishGig 규칙)
- 기대: ERROR — fullend.yaml roles에 "Client" 없음
- 결과: PASS — Phase013 구현 후 검출 성공

### MUT-POLICY-CONFIG-002: Claims 필드명 변경
- 대상: `specs/gigbridge/fullend.yaml`
- 변경: `claims.ID: user_id` → `claims.ID: userId`
- 기대: ERROR — Rego input.claims 참조와 fullend.yaml claims 불일치
- 결과: PASS — Phase014: CheckClaimsRego 추가로 검출
