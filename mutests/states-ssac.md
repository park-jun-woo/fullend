# Mutation Test — States ↔ SSaC

### MUT-STATES-SSAC-001: 상태 전이 함수명 오타
- 대상: `specs/gigbridge/states/gig.md`
- 변경: `draft --> open : PublishGig` → `draft --> open : Publishgig`
- 기대: ERROR — SSaC function "PublishGig"과 불일치
- 결과: PASS — States↔SSaC + States↔OpenAPI 3건 검출


### MUT-STATES-SSAC-002: transition 있지만 @state 가드 없음
- 대상: `specs/gigbridge/service/gig/publish_gig.ssac`
- 변경: `@state` 시퀀스 삭제 (diagram에는 PublishGig transition 유지)
- 기대: WARNING — 상태 전이가 정의되어 있지만 함수에 @state 가드 없음
- 결과: 미실행
