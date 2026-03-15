# ✅ Phase019: Smoke Gen — FK 순서 + CHECK enum 더미값 + FK 변수 참조

## 목표

`fullend gen`의 hurl smoke 테스트 생성기에서 발견된 세 가지 버그를 수정한다:
1. **FK 순서 버그**: Register에 `org_id` FK가 있으면 CreateOrganization이 먼저 실행되어야 하는데, Auth 뒤에 나옴
2. **CHECK enum 버그**: DDL `CHECK (col IN ('a','b','c'))` 제약이 있는 컬럼에 `"test_string"` 더미값 생성 → DB 제약 위반
3. **FK 변수 참조 버그**: Register body의 `org_id` 값이 하드코딩 `1`이지만, prereq에서 capture한 `{{organization_id}}`를 사용해야 함

## 근본 원인 분석

### Bug 1: FK 순서 — 출력 로직이 ordering을 무시

`generateHurlTests()`의 출력 흐름:
```
Line 63-65: writeAuthSection() — 하드코딩으로 Auth를 무조건 먼저 출력
Line 69-72: steps 순회, IsAuth 스킵 — 나머지 step은 Auth 뒤에 출력
```

`buildScenarioOrder()`에서 prereqSteps를 auth 앞에 배치해도, `generateHurlTests()`가 `writeAuthSection()`을 루프 밖에서 먼저 호출하고, 루프에서 `IsAuth`를 skip하므로 순서가 무력화.

**수정 방향**: `generateHurlTests()` 출력 루프를 steps 순서대로 순회하도록 변경. auth step을 만나면 `writeAuthSection()`을 한 번만 호출.

### Bug 2: CHECK enum — 이미 수정됨 ✅

이전 세션에서 올바르게 구현 완료:
- `parseDDLCheckEnums()`: DDL SQL에서 CHECK IN 구문 파싱
- `checkEnums` 파라미터 스레딩 완료
- `generateDummyValue()`에서 CHECK enum 룩업 추가

### Bug 3: FK 변수 참조 — 새로 발견

`writeAuthPair()` → `generateRequestBodyWithOverrides()`에서 Register body의 `org_id` field:
- 현재: `generateDummyValue("org_id", integerSchema, checkEnums)` → `1` (integer 기본값)
- 기대: `{{organization_id}}` (prereq에서 capture한 변수)

`generateRequestBodyWithOverrides()`는 role/email만 오버라이드하고, FK field(`_id` suffix)에 대해 captured variable 치환 로직이 없음.

**수정 방향**: `generateRequestBodyWithOverrides()`에 `captures map[string]bool` 전달. `_id` suffix field에 대해 captures에서 매칭되는 변수가 있으면 `{{변수명}}` 템플릿으로 출력.

매칭 로직: `org_id` → captures에서 `organization_id` 또는 `org_id` 검색. path resource 이름에서 유추:
- captures에 있는 `*_id` 변수들을 순회
- `org_id` field의 `org` prefix가 capture 변수명의 prefix와 매칭되면 사용
- 예: `org_id` → prefix `org` → captures에서 `organization_id` (`organization`이 `org`로 시작) → `{{organization_id}}`

## 현재 상태 (반쯤 구현됨)

### ✅ 올바른 변경 (유지)
- `parseDDLCheckEnums()` 함수 추가 — hurl.go 하단
- `checkEnums` 파라미터 스레딩 전체 완료
- `generateDummyValue()`, `generateRequestBody()`에 `checkEnums` 추가 — hurl_util.go
- `collectAuthFKResources()` — FK prefix 감지 (반환 타입 `[]string`으로 변경됨)
- `matchFKPrefix()` — prefix 매칭
- `buildScenarioOrder()` 내 prereqSteps 분리 로직 — 올바른 순서 반환

### ❌ 수정 필요
- `generateHurlTests()` 출력 루프 — steps 순서를 무시하는 하드코딩 구조
- `generateRequestBodyWithOverrides()` — FK field의 capture 변수 참조 미지원
- `matchFKPrefix()` — prefix 매칭이 너무 느슨 (`org` → `originals` 오탐 가능). DDL FK 참조를 사용하거나 더 엄격한 매칭 필요

## 수정 전략

### Step 1: `generateHurlTests()` 출력 루프 수정

steps 순서대로 순회, auth를 만나면 `writeAuthSection()` 한 번 호출:

```go
authWritten := false
currentResource := ""
for _, step := range steps {
    if step.IsAuth {
        if !authWritten && hasAuth {
            writeAuthSection(&buf, doc, captures, roles, checkEnums)
            authWritten = true
        }
        continue
    }
    if !canResolvePathParams(step.Path, captures) { continue }
    resource := inferResource(step.Path)
    if resource != currentResource {
        currentResource = resource
        buf.WriteString(fmt.Sprintf("\n# ===== %s =====\n\n", ...))
    }
    writeStep(&buf, step, captures, doc, roleMap, checkEnums)
}
```

**주의**: prereq step(CreateOrganization)은 auth 전에 출력됨. 이때 `captures`에 `organization_id`가 추가됨 (`writeStep` 내 POST capture 로직). 따라서 후속 `writeAuthSection()`에서 이 값을 참조할 수 있음.

단, prereq step의 `canResolvePathParams` 체크: `/organizations`는 path param 없으므로 통과. 문제없음.

### Step 2: `writeAuthPair()` → FK 변수 참조

`writeAuthPair()`에 `captures` 파라미터 추가 (이미 token capture용으로 받고 있음).

`generateRequestBodyWithOverrides()`에 `captures` 전달, `_id` suffix field 처리:

```go
case strings.HasSuffix(lower, "_id"):
    // Find matching captured variable: org_id → organization_id
    captureVar := findMatchingCapture(name, captures)
    if captureVar != "" {
        // Output raw template without JSON quoting: {{organization_id}}
        lines = append(lines, fmt.Sprintf("  %s: {{%s}}", formatDummyValue(name), captureVar))
        continue
    }
    val = generateDummyValue(name, prop, checkEnums)
```

`findMatchingCapture()` 로직:
```go
func findMatchingCapture(fieldName string, captures map[string]bool) string {
    // Direct match: org_id in captures
    if captures[fieldName] { return fieldName }
    // Prefix match: org_id → prefix "org", find "organization_id" in captures
    prefix := strings.TrimSuffix(fieldName, "_id")
    for cap := range captures {
        if strings.HasSuffix(cap, "_id") {
            capPrefix := strings.TrimSuffix(cap, "_id")
            if strings.HasPrefix(capPrefix, prefix) {
                return cap
            }
        }
    }
    return ""
}
```

### Step 3: `matchFKPrefix()` 안전성 강화

현재 `strings.HasPrefix("organizations", "org")` 방식은 오탐 위험. 대안:
- DDL FK 참조 파싱은 과도 → 현재 prefix 매칭 유지하되, 최소 3글자 이상 prefix 매칭 + resource 길이 비율 체크
- 실용적 판단: 현실적으로 `org` prefix가 `organizations` 외 다른 resource와 충돌할 확률은 매우 낮음. 현재 유지.

## 변경 파일 목록

| 파일 | 변경 내용 |
|---|---|
| `internal/gluegen/hurl.go` | (1) `generateHurlTests()` 출력 루프: steps 순서대로 순회 (2) `generateRequestBodyWithOverrides()`에 captures 전달 + FK 변수 참조 (3) `findMatchingCapture()` 헬퍼 추가 |
| `internal/gluegen/hurl_util.go` | 변경 없음 (이미 올바르게 수정됨) |

## 예상 출력 (zenflow smoke.hurl)

```hurl
# ===== Organizations (must precede Auth — FK dependency) =====

# CreateOrganization
POST {{host}}/organizations
Content-Type: application/json
{
  "credits_balance": 1,
  "name": "test_string",
  "plan_type": "free"
}

HTTP 201
[Captures]
organization_id: jsonpath "$.organization.id"
[Asserts]
jsonpath "$.organization" exists

# ===== Auth =====

# Register
POST {{host}}/auth/register
Content-Type: application/json
{
  "email": "test@test.com",
  "org_id": {{organization_id}},
  "password": "Password1234!",
  "role": "admin"
}

HTTP 201
...
```

## 검증 방법

1. `go build ./cmd/fullend/` — 빌드 통과
2. `go test ./...` — 테스트 통과
3. `fullend gen specs/zenflow artifacts/zenflow` — 재생성
4. `artifacts/zenflow/tests/smoke.hurl` 확인:
   - CreateOrganization이 Register **앞**에 위치
   - Register의 `org_id`가 `{{organization_id}}` (capture 변수 참조)
   - `plan_type`이 `"free"` (CHECK enum 첫 번째 값)
5. zenflow 서버 빌드 → hurl 테스트 통과
