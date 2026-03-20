# Mutation Test — OpenAPI 단독

### MUT-OPENAPI-001: path parameter 이름 충돌
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: `/gigs/{GigID}/proposals` → `/gigs/{ID}/proposals` (기존 `/gigs/{GigID}`와 같은 위치에 다른 파라미터명)
- 기대: ERROR — segment[1]에 {GigID}와 {ID} 혼재
- 결과: PASS — checkPathParamConflicts에서 동일 세그먼트 위치 파라미터명 충돌 검출

### MUT-OPENAPI-002: OpenAPI 파싱 오류
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: `paths:` 이하 들여쓰기 파괴
- 기대: ERROR — OpenAPI 파싱 실패
- 결과: PASS — openapi3.NewLoader().LoadFromFile 단계에서 구조 오류 검출
