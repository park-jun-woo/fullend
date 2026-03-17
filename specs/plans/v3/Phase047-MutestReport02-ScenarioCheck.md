# Phase047: Scenario 교차 검증 재확인

## 목표

mutest-report02 FAIL 중 scenario-check 도메인 2건을 재확인하고, 필요 시 수정한다.

## 사전 재확인 (구현 전 수동 실행)

2건 모두 기존 코드에 매칭 로직이 이미 존재한다.

| ID | 기존 로직 | 재확인 방법 |
|---|---|---|
| MUT-SCENARIO-OPENAPI-001 | `validateHurlEntry:13` — `findMatchingRoute(hurlSegs, e.Method, routes)` path 매칭 | zenflow openapi.yaml에서 `/auth/login` → `/auth/Login` 변경 → validate → hurl의 `/auth/login`이 "path not found in OpenAPI" ERROR인지 확인 |
| MUT-SCENARIO-OPENAPI-002 | `validateHurlEntry:23` — `matchedRoute == nil` method 매칭 | zenflow scenario-happy-path.hurl에서 `POST /auth/register` → `GET /auth/register` 변경 → validate → "method not defined" ERROR인지 확인 |

- **PASS로 확인되면**: agent 실행 오류. 수정 불필요
- **FAIL로 확인되면**: `findMatchingRoute` 또는 `normalizeHurlPath`의 정규화 로직에 버그. 아래 수정 방향 적용

## 조건부 수정 (재확인 FAIL 시에만)

### MUT-SCENARIO-OPENAPI-001 — path 매칭 실패

- **chain**: `CheckHurlFiles` → `parseHurlFile` + `buildHurlRoutes` + `validateHurlEntry` → `findMatchingRoute` → `normalizeHurlPath` + `segmentsMatch`
- **예상 원인**: `normalizeHurlPath`나 `normalizeOpenAPIPath`에서 경로를 소문자로 정규화하여 `/auth/Login`과 `/auth/login`이 동일하게 매칭될 가능성
- **수정**: 정규화 과정에서 path segment의 대소문자를 보존하도록 변경

### MUT-SCENARIO-OPENAPI-002 — method 매칭 실패

- **예상 원인**: `buildHurlRoutes`가 OpenAPI에서 라우트를 만들 때 method를 대소문자 정규화하거나, `findMatchingRoute`에서 method 비교가 case-insensitive일 가능성
- **수정**: method 비교가 정확한 대소문자 매칭인지 확인

## 의존성

- 추가 외부 패키지 없음

## 검증 방법

1. 사전 재확인 2건 수동 실행
2. FAIL 확인 시: `go test ./internal/crosscheck/...` 통과
3. `go test ./...` 통과
4. zenflow-try05 대상 뮤테이션 재실행으로 PASS 전환 확인
