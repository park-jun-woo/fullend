# Mutation Test — Config ↔ OpenAPI

### MUT-CONFIG-OPENAPI-001: Middleware 제거
- 대상: `specs/gigbridge/fullend.yaml`
- 변경: `middleware: [bearerAuth]` → `middleware: []`
- 기대: ERROR — OpenAPI securitySchemes bearerAuth 존재하지만 미들웨어 미선언
- 결과: PASS — Phase014: middleware nil 조건 제거로 검출

### MUT-CONFIG-OPENAPI-002: endpoint security가 middleware에 미정의
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: endpoint에 `security: [{apiKey: []}]` 추가 (middleware에 apiKey 없음)
- 기대: ERROR — endpoint가 참조하는 security name이 middleware에 없음
- 결과: 미실행
