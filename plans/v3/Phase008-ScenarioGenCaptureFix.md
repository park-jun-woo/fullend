✅ 완료

# Phase008: scenario-gen 캡처 변수 불일치 수정 (BUG013)

## 목표

scenario-gen이 생성하는 Hurl 파일에서 캡처 변수명과 URL 템플릿 변수명이 일치하도록 수정한다.

## 문제 분석

### 근본 원인

두 개의 독립된 네이밍 시스템이 서로를 모르고 동작:

1. **캡처 네이밍** (`inferScenarioCaptureNamed`, line 435): `captureName + "_" + responseProperty + "_id` → `project_project_id`
2. **URL 변수 네이밍** (`buildScenarioURL`, line 281): JSON body에 var ref 없으면 `pascalToSnakeHurl(paramName)` → `id`

결과: `[Captures]`에는 `project_project_id`로 저장하지만 URL에서는 `{{id}}`로 참조 → `Undefined variable: id`

### 부차 문제

- **중복 Register**: 한 파일 내 여러 시나리오가 동일 이메일로 Register → unique 제약 위반
- **중복 접두사**: `project_project_id`는 `{capture}_{response_field}_{id_field}` 규칙에서 capture와 response_field가 같을 때 의미 없는 중복 발생

## 설계

### 현재 .feature 문법

```gherkin
When POST CreateProject {...} → projectResult
# 캡처 이름만 지정, 어떤 필드를 캡처할지는 암시적
```

### 현재 동작의 5가지 암시 규칙

| # | 규칙 | 위치 |
|---|---|---|
| 1 | 응답에서 object.id를 찾아 캡처 | `inferScenarioCaptureNamed` |
| 2 | 캡처 변수명 = `{capture}_{objName}_{idField}` | `inferScenarioCaptureNamed` |
| 3 | token 문자열 포함 → 토큰 캡처 | `writeActionHurlV2` |
| 4 | URL 파라미터 값을 JSON body에서 탐색 | `buildScenarioURL` |
| 5 | JSON에 값 없으면 paramName → snake_case 변환 | `buildScenarioURL` |

문제는 **규칙 2**와 **규칙 5**가 같은 변수를 다른 이름으로 생성하는 것.

### 수정 방안: URL fallback에서 captures 맵 활용

**핵심**: `buildScenarioURL`의 fallback(규칙 5)에서 captures 맵을 참조하여 이미 캡처된 변수 중 해당 path param에 매칭되는 것을 찾는다.

```
# 변경 전 (line 281)
hurlVar := pascalToSnakeHurl(paramName)  // → "id"

# 변경 후
hurlVar := findCapturedVar(paramName, captures)  // captures에서 "_id" 접미사 매칭 → "project_project_id"
if hurlVar == "" {
    hurlVar = pascalToSnakeHurl(paramName)  // 최후 fallback
}
```

### 캡처 변수 네이밍 단순화

`inferScenarioCaptureNamed`에서 capture 이름과 response field 이름이 같으면 중복 제거:

```
# 변경 전
captureName + "_" + name + "_id"  // → "project_project_id"

# 변경 후
capture="project", response field="project" → "project_id"  (중복 제거)
capture="gig",     response field="gig"     → "gig_id"
capture="result",  response field="project" → "result_project_id"  (다르면 유지)
```

### 이메일 유니크화

`deriveEmailPrefix`가 이미 파일 단위 접두사를 생성하지만, 같은 파일 내 여러 시나리오 간 충돌은 미처리.

수정: 시나리오 인덱스를 이메일에 반영.

```go
// renderFeatureHurl에서 시나리오 순회 시
emailPrefix := fmt.Sprintf("%s%d", deriveEmailPrefix(f.File), i)
```

### 크로스체크 강화

현재 `checkCaptureRefs`는 루트 변수명(`gigResult`)만 검증. 실제 생성 시 `inferScenarioCaptureNamed`가 만드는 Hurl 변수명(`gig_gig_id`)과 URL에서 참조하는 변수명이 일치하는지는 미검증.

#### 추가 검증 1: 캡처 변수 → 응답 스키마 경로 유효성

`.feature`에서 `→ gig` 캡처 후 `gig.gig.id`를 참조할 때, OpenAPI 응답 스키마에 `gig` object와 `id` field가 실제로 존재하는지 검증.

```go
// checkCapturePathValidity: 캡처 변수의 dotted 참조 경로가 OpenAPI 응답 스키마와 일치하는지 확인
// e.g. "gig.gig.id" → CreateGig 응답에 $.gig.id 존재? ✓
//      "gig.gig.name" → CreateGig 응답에 $.gig.name 존재? ✗ → WARNING
```

#### 추가 검증 2: 동일 파일 내 중복 이메일 검출

같은 `.feature` 파일의 서로 다른 시나리오에서 동일 이메일로 Register 시 WARNING.

```go
// checkDuplicateEmails: 파일 내 모든 시나리오의 Register JSON에서 이메일 추출,
// 시나리오 간 중복 발견 시 WARNING
```

## 변경 파일

| 파일 | 변경 |
|---|---|
| `internal/gluegen/hurl_scenario.go` | `inferScenarioCaptureNamed` 중복 접두사 제거 |
| `internal/gluegen/hurl_scenario.go` | `buildScenarioURL` 시그니처에 `captures` 추가, fallback에서 캡처 변수 매칭 |
| `internal/gluegen/hurl_scenario.go` | `renderFeatureHurl` 시나리오별 emailPrefix 생성 |
| `internal/gluegen/hurl_scenario.go` | `writeActionHurlV2`에서 `buildScenarioURL` 호출에 captures 전달 |
| `internal/crosscheck/scenario.go` | `checkCapturePathValidity` 추가 — dotted ref 경로가 OpenAPI 응답 스키마에 존재하는지 검증 |
| `internal/crosscheck/scenario.go` | `checkDuplicateEmails` 추가 — 파일 내 시나리오 간 중복 이메일 WARNING |
| `internal/crosscheck/scenario.go` | `CheckScenarios`에서 두 검증 호출 추가 |
| `internal/gluegen/hurl_scenario_test.go` | 캡처-URL 일치 테스트, 중복 접두사 제거 테스트, 이메일 유니크 테스트 |
| `internal/crosscheck/scenario_test.go` | 캡처 경로 유효성 테스트, 중복 이메일 테스트 |

## 검증

1. `go test ./internal/gluegen/...` — 유닛 테스트 통과
2. `go test ./internal/crosscheck/...` — 크로스체크 테스트 통과
3. `fullend validate specs/gigbridge` — 검증 통과 (새 검증 규칙 포함)
4. `fullend gen specs/gigbridge artifacts/gigbridge` — 생성된 hurl 파일에서 캡처 변수와 URL 변수 일치 확인
5. 생성된 `scenario-*.hurl` 수동 확인: `[Captures]`의 변수명과 이후 `{{변수명}}`이 동일한지 검증

## 의존성

없음. fullend 내부 변경만.

## 상태: REJECT — Phase009(.scene 포맷)로 대체
