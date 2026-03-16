# Mutation Test — OpenAPI ↔ DDL

### MUT-OPENAPI-DDL-001: DDL 컬럼명 언더스코어
- 대상: `specs/gigbridge/db/gigs.sql`
- 변경: `client_id` → `clientid` (언더스코어 제거)
- 기대: ERROR — OpenAPI x-model Gig의 `client_id`와 불일치
- 결과: PASS — x-include FK + Policy ownership 컬럼 부재 검출

### MUT-OPENAPI-DDL-002: DDL 컬럼 추가 누락
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: Gig schema에 `budget` property 제거
- 기대: WARNING — DDL gigs.budget이 OpenAPI에 없음
- 결과: PASS — SSaC @post budget 필드 OpenAPI 부재 검출

### MUT-OPENAPI-DDL-003: OpenAPI 유령 property
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: Gig schema에 `rating: { type: integer }` 추가
- 기대: WARNING — DDL gigs에 rating 컬럼 없음
- 결과: PASS — Phase014: checkGhostProperties 추가로 검출 (ERROR)

### MUT-OPENAPI-DDL-004: cursor pagination + x-sort allowed 2개 이상
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: ListGigs에 `x-pagination: { style: cursor }` + `x-sort: { allowed: [created_at, title] }`
- 기대: ERROR — cursor pagination에서 런타임 정렬 전환은 cursor를 깨뜨림
- 결과: 미실행

### MUT-OPENAPI-DDL-005: cursor pagination + non-UNIQUE x-sort default
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: ListGigs cursor pagination의 x-sort default를 UNIQUE 인덱스 없는 컬럼으로 지정
- 기대: ERROR — 중복값 시 cursor가 깨짐
- 결과: 미실행

### MUT-OPENAPI-DDL-006: x-include 잘못된 형식
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: `x-include: [client_id:users.id]` → `x-include: [client_id_users_id]` (콜론 누락)
- 기대: ERROR — x-include 형식 "column:table.column" 위반
- 결과: 미실행

### MUT-OPENAPI-DDL-007: x-include 대상 테이블 미존재
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: `x-include: [client_id:users.id]` → `x-include: [client_id:nonexistent.id]`
- 기대: ERROR — DDL에 "nonexistent" 테이블 없음
- 결과: 미실행
