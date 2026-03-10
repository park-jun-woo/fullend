✅ 완료

# Phase 010: scenario-gen Hurl 생성 버그 수정

## 목표

`GenerateScenarioHurl`이 생성하는 scenario/invariant hurl 파일의 3가지 버그를 수정한다.

## 배경

Gherkin `.feature` → `.hurl` 변환은 구조적으로 잘 동작하나, 실제 실행 시 실패하는 3가지 문제가 있다:

1. JSON 리터럴 값이 path 변수로 잘못 치환됨
2. assertion의 subField가 PascalCase (`ID`)로 생성되나 JSON 응답은 소문자 (`id`)
3. capture 미정의 변수가 path에 사용되어 런타임 에러

### 구체적 사례

**negative-auth.feature:**
```gherkin
When POST EnrollCourse {"CourseID": 1, "PaymentMethod": "card"}
Then status == 401
```

생성 결과 (버그):
```hurl
POST {{host}}/courses/{{course_id}}/enroll   ← course_id 미정의
```

기대 결과:
```hurl
POST {{host}}/courses/1/enroll               ← 리터럴 1 직접 삽입
```

**course-lifecycle.feature:**
```gherkin
And response.courses contains course.ID
```

생성 결과 (버그):
```hurl
jsonpath "$.courses[*].ID" includes {{course_id}}   ← ID (PascalCase)
```

기대 결과:
```hurl
jsonpath "$.courses[*].id" includes {{course_id}}    ← id (JSON 태그)
```

## 변경 사항

### 1. `hurl_scenario.go` `buildScenarioURL()` — 리터럴 값 직접 삽입

현재 `findJSONVarRef()`는 변수 참조(`course.ID`)만 반환하고, 리터럴(`1`)이면 빈 문자열을 반환 → fallback이 `{{course_id}}`로 변환.

수정: `findJSONVarRef()` 대신 `findJSONValue()`로 변경. 변수 참조와 리터럴 모두 반환하되 구분:

```go
func buildScenarioURL(pathTemplate, json string, captures map[string]bool) string {
    // ...
    // Try to find value in JSON body.
    val, isVarRef := findJSONValue(json, paramName)
    if val != "" {
        if isVarRef {
            hurlVar := varRefToHurl(val)
            result = result[:pos] + "{{" + hurlVar + "}}" + result[pos+closeBrace+1:]
        } else {
            // Literal value — insert directly.
            result = result[:pos] + val + result[pos+closeBrace+1:]
        }
    } else {
        // No JSON match — use snake_case variable (existing capture).
        hurlVar := pascalToSnakeHurl(paramName)
        result = result[:pos] + "{{" + hurlVar + "}}" + result[pos+closeBrace+1:]
    }
    // ...
}
```

`findJSONValue()` 반환: `(value string, isVarRef bool)`
- `"CourseID": course.ID` → `("course.ID", true)`
- `"CourseID": 1` → `("1", false)`
- `"CourseID": "abc"` → `("abc", false)` (따옴표 제거)

### 2. `hurl_scenario.go` `writeAssertLineV2()` — subField 소문자 변환

`contains`/`excludes` assertion에서 subField를 소문자로 변환:

```go
// 변경 전
buf.WriteString(fmt.Sprintf("jsonpath \"$.%s[*].%s\" includes %s\n", a.Field, subField, val))

// 변경 후
buf.WriteString(fmt.Sprintf("jsonpath \"$.%s[*].%s\" includes %s\n", a.Field, strings.ToLower(subField), val))
```

`ID` → `id`, `Email` → `email` 등 JSON 태그 기준.

### 3. `hurl_scenario.go` `buildScenarioBody()` — 리터럴 path param도 body에서 제거

현재 path param 필드를 body에서 제거하는 로직은 있으나, 리터럴 값일 때도 동일하게 제거되는지 확인. 기존 로직이 필드명 기준으로 제거하므로 추가 수정 불필요 (확인만).

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/gluegen/hurl_scenario.go` | 수정 — `findJSONVarRef` → `findJSONValue` 리팩터, assertion subField 소문자 변환 |

## 의존성

없음 — scenario-gen 내부 수정.

## 검증 방법

```bash
go build ./cmd/fullend/
./fullend gen specs/dummy-lesson artifacts/dummy-lesson

# 생성된 hurl 확인:
# 1. scenario-negative-auth.hurl에서 courses/1/enroll (리터럴)
# 2. scenario-course-lifecycle.hurl에서 $.courses[*].id (소문자)
# 3. invariant-course-deletion.hurl에서 $.courses[*].id (소문자)

# 서버 기동 후 전체 hurl 테스트:
cd artifacts/dummy-lesson/backend
JWT_SECRET=test-secret-key go run ./cmd/ -dsn "postgres://postgres:test1224@localhost:15432/dummy_lesson?sslmode=disable" &
hurl --test --variable host=http://localhost:8080 ../tests/*.hurl
```
