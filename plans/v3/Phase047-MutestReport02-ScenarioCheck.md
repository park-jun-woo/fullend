# Phase047: Scenario 교차 검증 — hurl 파서 URL 패턴 확장 ✅ 완료

## 목표

mutest-report02 FAIL 중 scenario-check 도메인 2건을 해결한다.

## 사전 재확인 결과

| ID | 계획 시 판단 | 재확인 결과 |
|---|---|---|
| MUT-SCENARIO-OPENAPI-001 | path 매칭 case-insensitive 의심 | **진짜 FAIL** — 근본 원인은 path 매칭이 아님 |
| MUT-SCENARIO-OPENAPI-002 | method 매칭 case-insensitive 의심 | **진짜 FAIL** — 근본 원인은 method 매칭이 아님 |

**공통 근본 원인**: `parse_hurl_file.go:13`의 정규식이 `{{host}}` 패턴만 지원:
```go
// 변경 전
reHurlRequest = regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH)\s+\{\{host\}\}(.+)`)
```
zenflow hurl 파일은 `http://localhost:8080/auth/register` 형식 → 정규식 미매칭 → entries 빈 배열 → 검증 미실행.

## 수정 내용

### 1. hurl 파서 URL 패턴 확장

- **파일**: `internal/crosscheck/parse_hurl_file.go`
- **chain**: `parseHurlFile` ← `CheckHurlFiles`
- **수정**: 정규식을 `{{host}}` 또는 `http(s)://호스트:포트` 패턴 모두 지원하도록 확장
```go
// 변경 후
reHurlRequest = regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH)\s+(?:\{\{host\}\}|https?://[^/]*)(/.+)`)
```
- **테스트**: `hurl_test.go`에 `TestParseHurlFile_AbsoluteURL` 추가 (http, https, 포트 포함 3건)

## 검증 결과

- `go test ./...` 전체 통과
- zenflow-try05 직접 검증:
  - 존재하지 않는 path (`/nonexistent/path`) → "path not found in OpenAPI" ERROR 검출 ✓
  - method 불일치 (`GET /auth/register` vs OpenAPI `POST`) → "method not defined in OpenAPI" ERROR 검출 ✓
- 원래 mutest 변이(OpenAPI `/auth/login` → `/auth/Login`)는 zenflow hurl에 `/auth/login` 요청이 없어서 직접 검출 불가. hurl 파서 자체는 정상 작동 확인
