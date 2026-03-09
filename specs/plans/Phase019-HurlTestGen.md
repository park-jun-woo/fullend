✅ 완료

# Phase 019: OpenAPI → Hurl 테스트 자동 생성

## 목표

OpenAPI 스펙에서 Hurl(`.hurl`) HTTP 시나리오 테스트를 자동 생성한다. `fullend gen` 파이프라인에 `hurl-gen` 단계를 추가하여, 코드젠 산출물에 런타임 검증 테스트를 포함시킨다.

```
현재:
  fullend gen → 빌드 가능한 코드 (컴파일 검증만)

목표:
  fullend gen → 빌드 가능한 코드 + Hurl 테스트 (HTTP 런타임 검증)
  hurl --test artifacts/dummy-lesson/tests/*.hurl
```

---

## Hurl 개요

- Apache License 2.0 (Orange-OpenSource/hurl)
- 평문 파일로 HTTP 요청/응답을 선언적으로 기술
- 변수 캡처 → 후속 요청에 주입 (시나리오 체이닝)
- jsonpath assertion으로 응답 구조/값 검증
- CLI: `hurl --test --variable host=http://localhost:8080 *.hurl`

---

## 생성 전략

### 1. 시나리오 구성

OpenAPI의 엔드포인트를 분석하여 하나의 `.hurl` 파일에 전체 시나리오를 구성한다.

#### 순서 결정 규칙

```
1단계: 인증 (securitySchemes가 있으면)
  - Register → Login → token 캡처

2단계: 리소스 CRUD (의존 순서)
  - 부모 리소스 먼저 (path depth 순)
  - 각 리소스 내: POST → GET(list) → GET(single) → PUT → DELETE
  - POST 응답에서 ID 캡처 → 후속 요청에 주입

3단계: 종속 리소스
  - path parameter가 부모 ID를 참조하면 부모 다음에 배치
  - 예: /courses/{CourseID}/lessons → courses CRUD 후에 배치
```

#### 의존 그래프 예시 (dummy-lesson)

```
register → login(token)
  → courses.create(token, course_id)
    → courses.get(course_id)
    → courses.update(token, course_id)
    → courses.publish(token, course_id)
    → lessons.create(token, course_id, lesson_id)
      → lessons.update(token, lesson_id)
      → lessons.delete(token, lesson_id)
    → reviews.create(token, course_id, review_id)
      → reviews.list(course_id)
      → reviews.delete(token, review_id)
    → enroll(token, course_id)
      → enrollments.list(token)
      → payments.list(token)
    → courses.delete(token, course_id)
```

### 2. 더미값 생성

OpenAPI schema의 type + format에서 파생:

| type | format | 더미값 |
|---|---|---|
| `string` | (없음) | `"test_string"` |
| `string` | `email` | `"test@example.com"` |
| `string` | `date-time` | `"2025-01-01T00:00:00Z"` |
| `integer` | (없음) | `1` |
| `integer` | `int64` | `1` |
| `boolean` | — | `true` |
| `number` | — | `1.0` |

필드명 기반 힌트:
- `*Password*` → `"Password1234!"` (예측 가능한 테스트 비밀번호)
- `*Price*`, `*Amount*` → `10000`
- `*Rating*` → `5`
- `*URL*` → `"https://example.com/test"`

### 3. Assertion 생성

response schema에서 파생:

```
- 최상위 필드 존재: jsonpath "$.course" exists
- 배열 타입: jsonpath "$.courses" isCollection
- 필드 값 일치 (POST로 보낸 값): jsonpath "$.course.Title" == "test_string"
- ID 존재: jsonpath "$.course.ID" exists
```

### 4. 인증 처리

```
securitySchemes.bearerAuth 존재
  → Register + Login을 시나리오 앞에 삽입
  → token 캡처: [Captures] token: jsonpath "$.token.AccessToken"
  → security 필요 엔드포인트에 Authorization: Bearer {{token}} 주입
```

### 5. x- 확장 활용

| x- 확장 | Hurl 생성 |
|---|---|
| `x-pagination` | `?limit=2` 파라미터 + `jsonpath "$.total"` assertion |
| `x-sort` | `?sort=created_at&direction=desc` 파라미터 |
| `x-filter` | `?category=dev` 파라미터 (더미 필터값) |
| `x-include` | `?include=instructor` 파라미터 + 포함 필드 assertion |

---

## 생성 예시

```hurl
# ===== Auth =====

# Register
POST {{host}}/register
Content-Type: application/json
{
  "Email": "test@example.com",
  "Password": "Password1234!",
  "Name": "test_string"
}

HTTP 200
[Asserts]
jsonpath "$.user" exists
jsonpath "$.user.ID" exists
jsonpath "$.user.Email" == "test@example.com"

# Login
POST {{host}}/login
Content-Type: application/json
{
  "Email": "test@example.com",
  "Password": "Password1234!"
}

HTTP 200
[Captures]
token: jsonpath "$.token.AccessToken"
[Asserts]
jsonpath "$.token" exists

# ===== Courses =====

# CreateCourse
POST {{host}}/courses
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "Title": "test_string",
  "Category": "test_string",
  "Level": "test_string",
  "Price": 10000
}

HTTP 200
[Captures]
course_id: jsonpath "$.course.ID"
[Asserts]
jsonpath "$.course" exists
jsonpath "$.course.Title" == "test_string"

# ListCourses (with pagination + sort + filter + include)
GET {{host}}/courses?limit=2&sort=created_at&direction=desc&category=test_string&include=instructor

HTTP 200
[Asserts]
jsonpath "$.courses" isCollection
jsonpath "$.total" exists

# GetCourse (with include)
GET {{host}}/courses/{{course_id}}?include=instructor

HTTP 200
[Asserts]
jsonpath "$.course.ID" == {{course_id}}
jsonpath "$.course" exists

# UpdateCourse
PUT {{host}}/courses/{{course_id}}
Authorization: Bearer {{token}}
Content-Type: application/json
{
  "Title": "updated_string",
  "Category": "test_string",
  "Level": "test_string",
  "Price": 20000
}

HTTP 200

# DeleteCourse
DELETE {{host}}/courses/{{course_id}}
Authorization: Bearer {{token}}

HTTP 200
```

---

## 구현

### 새 파일

| 파일 | 역할 |
|---|---|
| `artifacts/internal/gluegen/hurl.go` | 1단계 시나리오 구성 (OpenAPI → smoke.hurl) |
| `artifacts/internal/gluegen/hurl_util.go` | 공용 유틸: 더미값 생성, assertion 생성, 인증 헤더 등 (2단계에서 재사용) |

### 수정 파일

| 파일 | 변경 |
|---|---|
| `artifacts/internal/gluegen/gluegen.go` | `Generate()`에 `generateHurlTests()` 호출 추가 |
| `artifacts/internal/orchestrator/gen.go` | hurl-gen 성공 리포트 출력 |

### 출력

```
<artifacts-dir>/
  tests/
    smoke.hurl             # 1단계: OpenAPI 스모크 테스트 (자동 생성, 매번 덮어씀)
    # 차후 2단계 추가 시:
    # usecase-course-lifecycle.hurl   ← UseCase SSOT에서 생성
    # usecase-student-enrollment.hurl
```

1단계는 `smoke.hurl` 단일 파일. 시나리오가 순서 의존적이므로 하나의 파일이 자연스럽다. 2단계 UseCase 시나리오는 usecase별로 별도 파일로 생성되며, `hurl --test tests/*.hurl`로 전체 실행 시 함께 돌아간다.

### 핵심 함수

```go
// generateHurlTests는 OpenAPI doc에서 Hurl 시나리오를 생성한다.
func generateHurlTests(doc *openapi3.T, outDir string) error

// buildScenarioOrder는 엔드포인트를 의존 순서로 정렬한다.
// 1. auth (register, login)
// 2. 부모 리소스 CRUD (path depth 순)
// 3. 자식 리소스 CRUD
func buildScenarioOrder(doc *openapi3.T) []scenarioStep

// generateDummyValue는 schema type+format에서 더미값을 생성한다.
func generateDummyValue(fieldName string, schema *openapi3.Schema) interface{}

// generateAssertions는 response schema에서 jsonpath assertion을 생성한다.
func generateAssertions(schema *openapi3.Schema, prefix string) []string
```

---

## 문서 업데이트

### CLAUDE.md

`CLI 명령어` > `fullend gen` 단계에 추가:

```
7. hurl-gen                        ← OpenAPI → Hurl 시나리오 테스트
```

`교차 검증 규칙` 섹션 하단에 추가:

```
### 런타임 검증
fullend gen이 생성하는 Hurl 테스트로 HTTP 런타임 검증을 수행한다.
hurl --test --variable host=http://localhost:8080 artifacts/<project>/tests/*.hurl
```

### README.md

Cross-Validation 섹션 아래에 Runtime Testing 섹션 추가:

```markdown
## Runtime Testing

`fullend gen` also generates [Hurl](https://hurl.dev) test scenarios from OpenAPI specs.

\`\`\`bash
# Start your server, then:
hurl --test --variable host=http://localhost:8080 artifacts/my-project/tests/*.hurl
\`\`\`
```

### manual-for-ai.md

`fullend CLI` 섹션의 gen 설명에 hurl-gen 단계 추가. `5-SSOT 연결 맵` 하단에 Hurl 테스트 위치 명시:

```
OpenAPI → Hurl 시나리오 (.hurl)
  - 엔드포인트별 요청/응답 검증
  - 인증 흐름 자동 포함
  - x-pagination/sort/filter/include 파라미터 반영
```

---

## 2단계 확장 설계

Phase019는 1단계(스모크 테스트)만 구현하지만, 2단계(UseCase SSOT 시나리오)가 추가될 때 자연스럽게 확장 가능한 구조로 설계한다.

### 파일 구조

```
gluegen/
  hurl_util.go       ← 공용 (1단계 + 2단계 공유)
                        generateDummyValue, generateAssertions,
                        writeAuthPreamble, writeCaptureID
  hurl.go            ← 1단계: OpenAPI → smoke.hurl
  hurl_usecase.go    ← 2단계 (차후): UseCase SSOT → usecase-*.hurl
```

### 2단계가 hurl_util.go에서 재사용하는 것

| 유틸 | 1단계 사용 | 2단계 사용 |
|---|---|---|
| `generateDummyValue` | 스키마에서 더미 request body | 시나리오 step의 파라미터 값 |
| `generateAssertions` | response 스키마에서 assertion | invariant에서 assertion |
| `writeAuthPreamble` | securitySchemes에서 인증 흐름 | policy.role에서 다중 사용자 인증 |
| `writeCaptureID` | POST 응답 ID 캡처 | 시나리오 step 간 데이터 흐름 |

### 2단계에서 추가되는 것

- UseCase SSOT 파서 (YAML → 구조체)
- 네거티브 테스트 생성 (policy 위반, guard 조건, 상태 전이 위반)
- 다중 사용자 시나리오 (instructor + student 각각 register/login)
- invariant 검증 (삭제 후 목록 확인 등)

---

## 의존성

- **Phase018 완료** ✅
- **Hurl CLI**: 런타임에 `hurl` 설치 필요 (생성은 fullend만으로 가능, 실행 시에만 hurl 필요)
- **kin-openapi**: 이미 사용 중

## 검증

```bash
# 1. fullend gen으로 Hurl 파일 생성
fullend gen specs/dummy-lesson artifacts/dummy-lesson

# 2. 생성된 파일 확인
cat artifacts/dummy-lesson/tests/smoke.hurl

# 3. 시나리오 구조 확인
# - register → login 선행
# - token 캡처 → 후속 요청 주입
# - CRUD 순서 (POST → GET → PUT → DELETE)
# - 부모 → 자식 순서
# - x- 확장 파라미터 반영

# 4. (서버 실행 후) 실제 테스트
hurl --test --variable host=http://localhost:8080 artifacts/dummy-lesson/tests/smoke.hurl
```
