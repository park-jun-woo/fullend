# Mutation Test — States ↔ SSaC

### MUT-STATES-SSAC-001: 상태 전이 함수명 오타
- 대상: `specs/gigbridge/states/gig.md`
- 변경: `draft --> open : PublishGig` → `draft --> open : Publishgig`
- 기대: ERROR — SSaC function "PublishGig"과 불일치
- 결과: PASS — States↔SSaC + States↔OpenAPI 3건 검출

