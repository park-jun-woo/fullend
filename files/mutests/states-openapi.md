# Mutation Test — States ↔ OpenAPI

### MUT-STATES-OPENAPI-001: transition event가 SSaC에 없음
- 대상: `specs/gigbridge/states/gig.md`
- 변경: `draft --> open : PublishGig` → `draft --> open : PublishGigTYPO`
- 기대: ERROR — transition event "PublishGigTYPO"에 대응하는 SSaC 함수 없음
- 결과: PASS — crosscheck/states에서 transition event → SSaC 함수 매칭 검증

### MUT-STATES-OPENAPI-002: SSaC @state가 존재하지 않는 diagram 참조
- 대상: `specs/gigbridge/service/gig/publish_gig.ssac`
- 변경: `@state gig` → `@state nonexistent`
- 기대: ERROR — "nonexistent" diagram이 states/ 디렉토리에 없음
- 결과: PASS — crosscheck/states에서 @state → stateDiagram 존재 검증
