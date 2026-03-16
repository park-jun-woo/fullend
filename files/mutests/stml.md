# Mutation Test — STML 단독

### MUT-STML-001: data-fetch operationId 불일치
- 대상: `specs/gigbridge/frontend/gig-list.html`
- 변경: `data-fetch="ListGigs"` → `data-fetch="listgigs"` (대소문자 변경)
- 기대: ERROR — OpenAPI operationId "ListGigs"와 불일치
- 결과: PASS — STML validator가 operationId 존재 여부 검증

### MUT-STML-002: data-action에 GET 메서드 사용
- 대상: `specs/gigbridge/frontend/gig-detail.html`
- 변경: `data-action="GetGig"` (GET 메서드인 operationId를 action에 사용)
- 기대: ERROR — data-action은 non-GET 메서드만 허용
- 결과: PASS — STML validator가 action 블록의 HTTP 메서드 검증

### MUT-STML-003: 존재하지 않는 컴포넌트 참조
- 대상: `specs/gigbridge/frontend/gig-list.html`
- 변경: `<x-component name="GigCard">` → `<x-component name="NonExistent">`
- 기대: ERROR — 컴포넌트 파일 미존재
- 결과: PASS — STML validator가 컴포넌트 파일 존재 여부 검증
