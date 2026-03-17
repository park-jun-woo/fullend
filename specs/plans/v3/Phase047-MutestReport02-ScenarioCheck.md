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
reHurlRequest = regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH)\s+\{\{host\}\}(.+)`)
```
zenflow hurl 파일은 `http://localhost:8080/auth/register` 형식 → 정규식 미매칭 → entries 빈 배열 → 검증 미실행.

## 변경 파일 목록

### 1. hurl 파서 URL 패턴 확장

- **파일**: `internal/crosscheck/parse_hurl_file.go`
- **chain**: `parseHurlFile` ← `CheckHurlFiles`
- **현상**: `reHurlRequest`가 `{{host}}` 필수. `http://localhost:8080` 등 절대 URL 미지원
- **수정**: 정규식을 `{{host}}` 또는 `http(s)://호스트:포트` 패턴 모두 지원하도록 확장. path 부분만 추출
- **테스트**: `hurl_test.go`에 `POST http://localhost:8080/auth/register` 파싱 케이스 추가

## 의존성

- 추가 외부 패키지 없음

## 검증 방법

1. `go test ./internal/crosscheck/...` 통과
2. `go test ./...` 통과
3. zenflow-try05 대상 2건 뮤테이션 재실행으로 PASS 전환 확인
