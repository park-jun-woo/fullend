# Phase009: 시나리오 테스트 — 자체 DSL 폐기, Hurl 직접 작성

## 목표

Gherkin `.feature` 파서와 시나리오 Hurl 변환기를 제거한다. 시나리오 테스트는 사용자가 `.hurl`을 직접 작성하며, 이것이 시나리오 SSOT가 된다.

## 배경

기존 흐름: `.feature` (자체 DSL) → fullend 변환 → `.hurl` → hurl 실행

문제:
- 변환 과정에서 캡처 변수 불일치 (BUG013)
- 추론 함수 3개, 암시 규칙 5개 — 복잡하고 깨지기 쉬움
- 자체 DSL을 유지보수하는 비용 대비 이득 없음

## 변경 후 구조

```
자동 생성 (fullend gen):
  OpenAPI → 스모크 테스트 .hurl (엔드포인트별 단건 호출, 기존 hurl-gen)

시나리오 SSOT (사용자 직접 작성):
  tests/scenario-*.hurl

검증 (fullend validate):
  scenario-*.hurl → OpenAPI 단방향 크로스체크
```

## Hurl 시나리오 SSOT 규칙

사용자가 직접 작성하는 시나리오 `.hurl` 파일:

- 위치: `tests/scenario-*.hurl` (시나리오), `tests/invariant-*.hurl` (불변식)
- 내용: 표준 Hurl 문법 그대로. 자체 확장 없음
- 예시:

```hurl
# === Happy Path: Full Gig Lifecycle ===

# Register
POST {{host}}/auth/register
Content-Type: application/json
{
  "email": "client@test.com",
  "password": "pass123",
  "role": "client",
  "name": "Test Client"
}

HTTP 200

# Login
POST {{host}}/auth/login
Content-Type: application/json
{
  "email": "client@test.com",
  "password": "pass123"
}

HTTP 200
[Captures]
token: jsonpath "$.access_token"

# CreateGig
POST {{host}}/gigs
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "title": "Build Website",
  "description": "Need a website built",
  "budget": 5000
}

HTTP 200
[Captures]
gig_id: jsonpath "$.gig.id"
[Asserts]
jsonpath "$.gig.status" == "draft"

# PublishGig
PUT {{host}}/gigs/{{gig_id}}
Authorization: Bearer {{token}}

HTTP 200
[Asserts]
jsonpath "$.gig.status" == "open"
```

## 크로스체크: Scenario → OpenAPI (단방향)

시나리오는 API의 소비자. 시나리오에서 호출하는 경로/메서드가 OpenAPI에 존재하는지만 확인한다. 역방향(모든 API가 시나리오에 있는가)은 커버리지 문제이므로 검증하지 않는다.

### 검증 내용

| 규칙 | 설명 | 수준 |
|---|---|---|
| 경로 존재 | `.hurl`의 URL path가 OpenAPI에 정의되어 있는가 | ERROR |
| 메서드 일치 | 해당 path의 HTTP 메서드가 OpenAPI에 정의되어 있는가 | ERROR |
| 상태코드 정의 | 기대하는 HTTP 상태코드가 OpenAPI responses에 있는가 | WARNING |

### 구현

라인 단위 정규식으로 `.hurl` 파일에서 요청/응답 쌍을 추출한다. 무거운 Hurl 파서 불필요.

```go
// 요청 라인: GET {{host}}/gigs/{{gig_id}}
reRequest = regexp.MustCompile(`^(GET|POST|PUT|DELETE|PATCH)\s+\{\{host\}\}(.+)`)

// 응답 라인: HTTP 200
reResponse = regexp.MustCompile(`^HTTP\s+(\d+)`)
```

**URL 매칭**: Hurl `{{변수}}`와 OpenAPI `{param}`은 이름이 다를 수 있으므로 세그먼트 단위 비교.

```
Hurl:    /gigs/{{gig_id}}   → ["gigs", ":param"]
OpenAPI: /gigs/{id}         → ["gigs", ":param"]
매칭 ✓
```

변환 로직:
- Hurl: `{{host}}` 제거, `{{...}}` → `:param`
- OpenAPI: `{...}` → `:param`
- 세그먼트 배열 비교로 경로 매칭

기존 `crosscheck/scenario.go`의 States, Policy 크로스체크는 제거한다. operationId 기반이었으므로 `.hurl` URL path에서 역매핑이 필요한데, 복잡도 대비 가치가 낮다.

## 삭제 대상

| 대상 | 이유 |
|---|---|
| `internal/scenario/` 전체 | Gherkin 파서 불필요 |
| `gluegen/hurl_scenario.go` 전체 | 시나리오 Hurl 변환 불필요 |
| `gluegen/hurl_util.go` 시나리오 전용 함수 | `inferScenarioCaptureNamed`, `inferScenarioCapture`, `inferCaptureField` |
| `crosscheck/scenario.go` 전체 | 새 `.hurl` 크로스체크로 대체 |
| `crosscheck/scenario_test.go` 전체 | 동일 |
| `specs/*/scenario/*.feature` | `.hurl`로 대체 |

## 변경/추가 파일

| 파일 | 변경 |
|---|---|
| `internal/crosscheck/hurl.go` | 신규 — `.hurl` → OpenAPI 단방향 크로스체크 |
| `internal/crosscheck/hurl_test.go` | 신규 — 크로스체크 테스트 |
| `internal/crosscheck/crosscheck.go` | `CheckScenarios` 제거, `CheckHurlFiles` 호출 추가 |
| `internal/orchestrator/validate.go` | `.feature` 탐지 → `.hurl` 탐지로 변경 |
| `internal/orchestrator/gen.go` | 시나리오 Hurl 생성 호출 제거 |
| `internal/orchestrator/detect.go` | SSOT 탐지에서 scenario kind 변경 |
| `artifacts/manual-for-ai.md` | 시나리오 섹션 — Hurl 직접 작성 안내로 교체 |

## .feature 처리

`.feature` 파일이 남아있으면 validate 시 **ERROR**:

```
ERROR  scenario/gig_lifecycle.feature: .feature is no longer supported. Delete this file.
       Write scenario tests directly in Hurl format: tests/scenario-*.hurl
       See: https://hurl.dev/docs/manual.html
```

## 검증

1. `go test ./internal/crosscheck/...` — `.hurl` 크로스체크 테스트
2. `go test ./...` — 삭제 후 전체 빌드/테스트 통과
3. `fullend validate specs/gigbridge` — `.hurl` 시나리오 크로스체크 동작 확인
4. `fullend gen specs/gigbridge artifacts/gigbridge` — 스모크 `.hurl`만 생성, 시나리오 `.hurl` 미생성 확인

## BUG013

이 Phase로 BUG013은 해당 없음(Won't Fix). 변환 레이어 자체가 제거되므로 버그 원인이 소멸한다.

## 의존성

없음. fullend 내부 변경만. Phase008(REJECT)을 대체한다.

## 상태: ✅ 완료
