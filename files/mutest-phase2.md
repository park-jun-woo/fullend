# Mutation Test Phase 2 — crosscheck 심화 검출력 검증

## 목적

Phase 1(MUT-01~15)에서 주요 경로를 커버했다. Phase 2는 아직 안 건드린 crosscheck 규칙을 검증한다.

## 방법

1. 변경 적용
2. `go run ./cmd/fullend validate specs/gigbridge` 실행
3. 검출 여부 기록 (PASS = 잡음, FAIL = 못 잡음)
4. `git checkout -- specs/gigbridge/` 로 되돌림

---

## 시나리오

### MUT-16: Claims 필드명 변경 (Claims ↔ SSaC)
- 대상: `specs/gigbridge/fullend.yaml`
- 변경: `claims.ID: user_id` → `claims.ID: userId`
- 기대: ERROR — SSaC currentUser.ID의 claim key 불일치

### MUT-17: Middleware 제거 (Middleware ↔ OpenAPI)
- 대상: `specs/gigbridge/fullend.yaml`
- 변경: `middleware: [bearerAuth]` → `middleware: []`
- 기대: ERROR — OpenAPI securitySchemes bearerAuth 존재하지만 미들웨어 미선언

### MUT-18: OpenAPI 에러 응답 삭제 (ErrStatus ↔ OpenAPI)
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: GetGig의 `404` 응답 삭제
- 기대: ERROR — SSaC @empty 기본 404에 대응하는 OpenAPI 응답 없음

### MUT-19: Policy ownership 컬럼 변경 (Policy ↔ DDL)
- 대상: `specs/gigbridge/policy/authz.rego`
- 변경: `@ownership gig gigs.client_id` → `@ownership gig gigs.owner_id`
- 기대: ERROR — DDL gigs 테이블에 owner_id 컬럼 부재

### MUT-20: DDL 고아 테이블 (DDL ↔ SSaC coverage)
- 대상: `specs/gigbridge/db/`
- 변경: `CREATE TABLE audit_logs (id BIGSERIAL PRIMARY KEY, action TEXT);` 파일 추가
- 기대: WARNING — SSaC에서 audit_logs 테이블을 사용하지 않음

### MUT-21: OpenAPI 유령 property (OpenAPI ↔ DDL)
- 대상: `specs/gigbridge/api/openapi.yaml`
- 변경: Gig schema에 `rating: { type: integer }` 추가
- 기대: WARNING — DDL gigs에 rating 컬럼 없음

### MUT-22: Sensitive 컬럼명 변경 (Sensitive ↔ DDL)
- 대상: `specs/gigbridge/db/users.sql`
- 변경: `password_hash` → `pw_hash`
- 기대: 기존 @sensitive 패턴(`password`, `secret`, `token` 등)에 안 걸림 — FAIL 예상

### MUT-23: Hurl HTTP method 변경 (Hurl ↔ OpenAPI)
- 대상: `specs/gigbridge/tests/scenario-gig-lifecycle.hurl`
- 변경: `POST {{host}}/auth/login` → `GET {{host}}/auth/login`
- 기대: ERROR — OpenAPI /auth/login에 GET 메서드 미정의

---

## 결과 기록

| ID | 변경 | 기대 | 실제 | 비고 |
|---|---|---|---|---|
| MUT-16 | claims.ID: user_id→userId | PASS | ~~FAIL~~ → PASS | Phase014: CheckClaimsRego 추가 — Rego claims 참조 ↔ fullend.yaml claims 대조 |
| MUT-17 | middleware 제거 | PASS | ~~FAIL~~ → PASS | Phase014: middleware nil 조건 제거 — OpenAPI 있으면 항상 검증 |
| MUT-18 | GetGig 404 응답 삭제 | PASS | PASS | @empty→OpenAPI 404 누락 정확히 검출 |
| MUT-19 | ownership gigs.client_id→owner_id | PASS | PASS | 재테스트 PASS — 초회 sed 패턴 오류 (콜론 누락) |
| MUT-20 | 고아 테이블 audit_logs 추가 | PASS | PASS | DDL→SSaC coverage WARNING 검출 |
| MUT-21 | Gig에 rating property 추가 | PASS | ~~FAIL~~ → PASS | Phase014: checkGhostProperties 추가 — OpenAPI→DDL 역방향 검증 (ERROR) |
| MUT-22 | password_hash→pw_hash | FAIL | FAIL→PASS(hash) | pw_hash는 "hash" 패턴으로 검출. Phase014: sensitive 패턴 20개로 확장 |
| MUT-23 | POST→GET /auth/login | PASS | PASS | Hurl method↔OpenAPI method 불일치 검출 |

**Phase014 후 결과: 8/8 검출 (100%)** — MUT-22는 "hash" 서브스트링 매칭으로 검출

### 미검출 분석

| ID | 원인 | 수정 난이도 | 가치 |
|---|---|---|---|
| MUT-16 | claims key는 코드젠 전용, crosscheck 대상 아님 | 중 | 낮음 — 런타임 JWT 파싱 에러로 즉시 발견 |
| MUT-17 | YAML 값 없음→nil→검증 스킵. 빈 배열이면 정상 검출 | 하 | 높음 — security 빠지면 모든 API 무방비 |
| MUT-21 | OpenAPI schema→DDL 역방향 property 대조 없음 | 중 | 낮음 — 유령 필드지만 동작에 영향 없음 |
| MUT-22 | sensitive 패턴 한계 (pw_, pwd_ 등 변형) | 하 | 낮음 — @sensitive 수동 태깅으로 보완 가능 |

### Phase014 수정 완료

- **MUT-16**: `CheckClaimsRego` 추가 — Rego `input.claims.xxx` ↔ fullend.yaml claims values
- **MUT-17**: `crosscheck.go` middleware nil 조건 제거 → OpenAPI 있으면 항상 검증
- **MUT-21**: `checkGhostProperties` 추가 — OpenAPI schema property → DDL column (ERROR)
- **MUT-22**: `sensitive.go` 패턴 4개 → 20개로 확장 (credential, otp, pin, ssn 등)
