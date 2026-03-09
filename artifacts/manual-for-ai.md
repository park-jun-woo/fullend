# fullend — AI SSOT Integration Guide

> Rules for writing 10 SSOTs (fullend.yaml, OpenAPI, SQL DDL, SSaC, STML, Mermaid stateDiagram, OPA Rego, Gherkin Scenario, Func Spec, Terraform) in a single project.
> Does not explain OpenAPI/SQL DDL/Terraform syntax. Covers only fullend.yaml/SSaC/STML/stateDiagram/OPA Rego/Gherkin/Func syntax and cross-SSOT connection rules.

## Project Directory Structure

```
<project-root>/
├── fullend.yaml                  # Project config (required)
├── api/openapi.yaml              # OpenAPI 3.x (with x- extensions)
├── db/
│   ├── *.sql                     # DDL (CREATE TABLE, CREATE INDEX)
│   └── queries/*.sql             # sqlc queries (-- name: Method :cardinality)
├── service/*.go                  # SSaC declarations (Go comment DSL)
├── model/*.go                    # Go interfaces (component definitions)
├── func/<pkg>/*.go               # Custom func implementations (optional)
├── states/*.md                   # Mermaid stateDiagram (state transitions)
├── policy/*.rego                 # OPA Rego (authorization policies)
├── scenario/*.feature            # Gherkin scenarios (fixed-pattern)
├── frontend/
│   ├── *.html                    # STML declarations (HTML5 + data-*)
│   ├── *.custom.ts               # Frontend computed functions (optional)
│   └── components/*.tsx          # React component wrappers (optional)
└── terraform/*.tf                # HCL infrastructure declarations
```

## fullend.yaml — Project Configuration

Required file at specs root. Kubernetes-style declarative YAML.

```yaml
apiVersion: fullend/v1
kind: Project

metadata:
  name: <project-name>

backend:
  lang: go                          # Backend language
  framework: gin                    # Backend framework
  module: github.com/org/project    # Go module path (used in go.mod, imports)
  middleware:                        # Middleware list (must match OpenAPI securitySchemes names)
    - bearerAuth

frontend:
  lang: typescript                  # Frontend language
  framework: react                  # Frontend framework
  bundler: vite                     # Bundler
  name: project-web                 # npm package name

deploy:
  image: ghcr.io/org/project        # Container image (optional)
  domain: project.example.com       # Service domain (optional)
```

### Required Fields

| Field | Description |
|-------|-------------|
| `apiVersion` | Must be `fullend/v1` |
| `kind` | Must be `Project` |
| `metadata.name` | Project identifier |
| `backend.module` | Go module path |

### Cross-validation Rules

| Rule | Level |
|------|-------|
| `backend.middleware` names must match OpenAPI `securitySchemes` keys | ERROR |
| OpenAPI `securitySchemes` keys must exist in `backend.middleware` | ERROR |
| Endpoint `security` references must exist in `backend.middleware` | ERROR |

## SSaC — Service Logic Declarations

### Syntax

```go
// @sequence <type>        — block start
// @model <Model.Method>   — resource model.method
// @param <Name> <source> [-> column]  — source: request | currentUser | varName | "literal". -> column: explicit DDL column mapping
// @result <var> <Type>    — result binding
// @message "msg" [STATUS] — custom error message (optional, STATUS default 500)
// @var <name>             — variable to return in response
// @action @resource @id   — authorize only (all 3 required)
// @component              — call only (component reference)
// @func <pkg.funcName>    — call only (package-level function reference)
```

### 10 Sequence Types

| Type | Purpose | Required Tags |
|---|---|---|
| authorize | Permission check | @action, @resource, @id |
| get | Single/list query | @model, @result |
| guard nil | Return error if nil | target variable name |
| guard exists | Return error if not nil | target variable name |
| guard state | Check state transition validity | target stateDiagramID, @param entity.Field |
| post | Create | @model, @result |
| put | Update | @model |
| delete | Delete | @model |
| call | External component/function call | @component or @func |
| response | JSON response return | (none, @var optional) |

### @func — Package-Level Function Call

`@func` references a package-level function with a standardized signature: `func(In) (Out, error)`.

```go
// Value form — captures result
// @sequence call
// @func auth.hashPassword
// @param Password request
// @result hashedPassword string

// Guard form — no result, error = rejection
// @sequence call
// @func auth.verifyPassword
// @param user.PasswordHash
// @param Password request
// @message "비밀번호가 일치하지 않습니다" 401
```

`@message` on a call sequence:
- With `@message`: error → responds with that message and status code (guard form)
- Without `@message`: error → responds with `"funcName 호출 실패" 500` (default)

### Example: All Sequence Types

```go
// @sequence authorize
// @action update
// @resource course
// @id CourseID
//
// @sequence get
// @model Course.FindByID
// @param CourseID request
// @result course Course
//
// @sequence guard nil course
// @message "Course not found"
//
// @sequence call
// @func auth.verifyPassword
// @param user.PasswordHash
// @param Password request
// @message "Wrong password" 401
//
// @sequence post
// @model Enrollment.Create
// @param CourseID request
// @param UserID currentUser
// @result enrollment Enrollment
//
// @sequence put
// @model Course.IncrementEnrollCount
// @param CourseID request
//
// @sequence call
// @component notification
// @param enrollment
//
// @sequence response json
// @var enrollment
func EnrollCourse(w http.ResponseWriter, r *http.Request) {}
```

### @param Source Rules

| Source | Meaning | Codegen |
|---|---|---|
| `request` | HTTP request body/query | `r.FormValue("Name")` |
| `currentUser` | Authenticated user info | `currentUser.Name` |
| variable name | @result variable from previous sequence | Direct reference |
| `var.Field` | Field of a previous @result variable | `var.Field` |
| `"literal"` | Hardcoded string | Used as-is |

`-> column` mapping: `@param PaymentMethod request -> method` — explicit DDL column mapping instead of auto snake_case conversion.

### Function Name = operationId

SSaC function names must match OpenAPI operationId exactly. This is the key that connects frontend (STML) and backend (SSaC).

```
OpenAPI: operationId: EnrollCourse
SSaC:    func EnrollCourse(...)
STML:    data-action="EnrollCourse"
```

## Func Spec — External Function Declarations

`func/<pkg>/*.go` files define custom function implementations. Each file follows a fixed pattern:

```go
package auth

import "golang.org/x/crypto/bcrypt"

// @func hashPassword
// @description 평문 비밀번호를 bcrypt 해시로 변환한다

type HashPasswordRequest struct {
    Password string
}

type HashPasswordResponse struct {
    HashedPassword string
}

func HashPassword(req HashPasswordRequest) (HashPasswordResponse, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    return HashPasswordResponse{HashedPassword: string(hash)}, err
}
```

### Rules

- **`@func`**: Function identifier (matches SSaC `@func pkg.funcName`)
- **`@description`**: Natural language one-liner (LLM uses this to implement the body)
- **Request/Response struct**: Go structs are the spec. No additional annotations needed.
- **Signature**: Always `func FuncName(req FuncNameRequest) (FuncNameResponse, error)`
- **Package-level function**: No Service struct dependency

### Fallback Chain

1. `specs/<project>/func/<pkg>/` — Project custom (highest priority)
2. `pkg/<pkg>/` — fullend default (fallback)
3. Neither → ERROR with skeleton suggestion

### fullend Default Functions (pkg/)

fullend ships with built-in default implementations in `pkg/`:

#### auth — 인증

| Function | Description | Request | Response |
|---|---|---|---|
| `hashPassword` | bcrypt 해싱 | `{ Password }` | `{ HashedPassword }` |
| `verifyPassword` | bcrypt 검증 | `{ PasswordHash, Password }` | `{}` (error=불일치) |
| `issueToken` | JWT 액세스 토큰 발급 (24h) | `{ UserID, Email, Role }` | `{ AccessToken }` |
| `verifyToken` | JWT 검증 → claims 추출 | `{ Token, Secret }` | `{ UserID, Email, Role }` |
| `refreshToken` | 리프레시 토큰 발급 (7일) | `{ UserID, Email, Role }` | `{ RefreshToken }` |
| `generateResetToken` | 비밀번호 리셋용 랜덤 토큰 | `{}` | `{ Token }` |

#### crypto — 암호화

| Function | Description | Request | Response |
|---|---|---|---|
| `encrypt` | AES-256-GCM 암호화 | `{ Plaintext, Key(hex) }` | `{ Ciphertext(base64) }` |
| `decrypt` | AES-256-GCM 복호화 | `{ Ciphertext(base64), Key(hex) }` | `{ Plaintext }` |
| `generateOTP` | TOTP 시크릿 + QR URL 생성 | `{ Issuer, AccountName }` | `{ Secret, URL }` |
| `verifyOTP` | TOTP 코드 검증 | `{ Code, Secret }` | `{}` (error=불일치) |

#### storage — S3 호환 파일 스토리지

| Function | Description | Request | Response |
|---|---|---|---|
| `uploadFile` | 파일 업로드 | `{ Bucket, Key, Data, ContentType, Endpoint, Region }` | `{ URL }` |
| `deleteFile` | 파일 삭제 | `{ Bucket, Key, Endpoint, Region }` | `{}` |
| `presignURL` | 서명된 다운로드 URL | `{ Bucket, Key, ExpiresIn, Endpoint, Region }` | `{ URL }` |

#### mail — 이메일

| Function | Description | Request | Response |
|---|---|---|---|
| `sendEmail` | SMTP 평문 이메일 | `{ Host, Port, Username, Password, From, To, Subject, Body }` | `{}` |
| `sendTemplateEmail` | Go 템플릿 HTML 이메일 | `{ Host, Port, Username, Password, From, To, Subject, TemplateName, Data }` | `{}` |

#### text — 텍스트 처리

| Function | Description | Request | Response |
|---|---|---|---|
| `generateSlug` | URL-safe slug 생성 | `{ Text }` | `{ Slug }` |
| `sanitizeHTML` | XSS 방지 HTML 정제 | `{ HTML }` | `{ Sanitized }` |
| `truncateText` | 유니코드 안전 텍스트 자르기 | `{ Text, MaxLength, Suffix }` | `{ Truncated }` |

#### image — 이미지 처리

| Function | Description | Request | Response |
|---|---|---|---|
| `ogImage` | OG 이미지 생성 (1200x630, PNG) | `{ Data }` | `{ Data }` |
| `thumbnail` | 썸네일 생성 (200x200, PNG) | `{ Data }` | `{ Data }` |

### SSaC Usage

SSaC에서 `@func pkg.funcName` 으로 참조:

```go
// @sequence call
// @func auth.hashPassword
// @param Password request
// @result hashedPassword string
```

생성 코드:
```go
out, err := auth.HashPassword(auth.HashPasswordRequest{Password: password})
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "hashPassword 호출 실패"})
    return
}
hashedPassword := out.HashedPassword
```

### Missing Func Error

구현이 없는 @func 참조 시, fullend validate가 스켈레톤을 자동 제안:

```
ERROR: @func billing.calculateRefund — 구현 없음

다음 파일을 작성하세요: func/billing/calculate_refund.go

package billing

// @func calculateRefund
// @description <이 함수가 무엇을 하는지 한 줄로 설명>

type CalculateRefundRequest struct {
    Reservation Reservation
}

type CalculateRefundResponse struct {
    Refund Refund
}

func CalculateRefund(req CalculateRefundRequest) (CalculateRefundResponse, error) {
    // TODO: implement
    return CalculateRefundResponse{}, nil
}
```

LLM에 이 에러 메시지를 그대로 전달하면 `@description`만 채우고 본체를 구현할 수 있다.

## Middleware — gin Middleware (pkg/middleware/)

fullend provides built-in gin middleware. Automatically wired based on OpenAPI `securitySchemes`.

### BearerAuth Middleware

```go
// pkg/middleware/bearerauth.go
func BearerAuth(secret string) gin.HandlerFunc
```

- `fullend.yaml` `backend.middleware`에 `bearerAuth` 선언 + OpenAPI `securitySchemes`에 `bearerAuth` 존재 시 적용
- `Authorization: Bearer <token>` → `auth.VerifyToken` → `c.Set("currentUser", &CurrentUser{...})`
- abort하지 않음 — 토큰 없거나 유효하지 않으면 빈 `CurrentUser{}` 세팅. authorize 시퀀스가 권한 검사 담당.

### CurrentUser Type

```go
// pkg/middleware/bearerauth.go
type CurrentUser struct {
    UserID int64
    Email  string
    Role   string
}
```

생성 프로젝트의 `model/auth.go`에서 타입 앨리어스로 연결:
```go
type CurrentUser = middleware.CurrentUser
```

SSaC codegen이 `c.MustGet("currentUser").(*model.CurrentUser)` 생성 → 미들웨어가 세팅한 값을 핸들러에서 사용.

### Route Grouping (OpenAPI security)

OpenAPI `security` 필드가 라우트 그룹 결정의 SSOT:

```yaml
paths:
  /login:
    post:
      operationId: Login
      # security 없음 → public group (미들웨어 없음)
  /courses:
    post:
      operationId: CreateCourse
      security:
        - bearerAuth: []    # → auth group (JWT 미들웨어 적용)
```

생성 코드:
```go
r := gin.Default()
auth := r.Group("/")
auth.Use(middleware.BearerAuth("secret"))

r.POST("/login", s.Auth.Login)              // public
auth.POST("/courses", s.Course.CreateCourse) // JWT 미들웨어 적용
```

## STML — UI Declarations

### Core data-* Attributes (8)

| Attribute | Value | Purpose | Placement |
|---|---|---|---|
| `data-fetch` | operationId | GET binding | Container element |
| `data-action` | operationId | POST/PUT/DELETE binding | form/button element |
| `data-field` | field name | Request body field | Inside data-action |
| `data-bind` | field name (dot notation) | Response field display | Inside data-fetch |
| `data-param-*` | `route.ParamName` | Path/query parameter | data-fetch or data-action element |
| `data-each` | array field name | List iteration | Inside data-fetch |
| `data-state` | condition expression | Conditional rendering | Anywhere |
| `data-component` | component name | React component delegation | Anywhere |

### Infrastructure data-* Attributes (3)

| Attribute | Value | Requirement |
|---|---|---|
| `data-paginate` | (no value, boolean) | Requires x-pagination in OpenAPI |
| `data-sort` | `column` or `column:desc` | Requires x-sort in OpenAPI |
| `data-filter` | `col1,col2` | Requires x-filter in OpenAPI |

### data-state Suffix Rules

| Pattern | Meaning | Codegen |
|---|---|---|
| `items.empty` | Array is empty | `{data.items?.length === 0 && ...}` |
| `items.loading` | Loading | `{isLoading && ...}` |
| `items.error` | Error occurred | `{isError && ...}` |
| `canEdit` | Boolean field | `{data.canEdit && ...}` |

### custom.ts Rules

When data-bind references a field not in the OpenAPI response schema, exporting a function with the same name in `<page>.custom.ts` passes validation.

```ts
// login-page.custom.ts
export function formattedDate(data) {
  return new Date(data.CreatedAt).toLocaleDateString()
}
```

### Example: Complex Page

```html
<main>
  <section data-fetch="ListCourses" data-paginate data-sort="created_at:desc" data-filter="category,level">
    <ul data-each="courses">
      <li>
        <h3 data-bind="title"></h3>
        <span data-bind="price"></span>
        <div data-component="RatingStars" data-bind="averageRating"></div>
      </li>
    </ul>
    <p data-state="courses.empty">No courses found</p>
    <div data-state="courses.loading">Loading...</div>
  </section>

  <form data-action="CreateCourse">
    <input data-field="Title" placeholder="Course title" />
    <input data-field="Price" type="number" placeholder="Price" />
    <select data-field="Category">
      <option value="dev">Development</option>
      <option value="design">Design</option>
    </select>
    <button type="submit">Create Course</button>
  </form>
</main>
```

## OpenAPI x- Extensions

Declare infrastructure parameters on OpenAPI endpoints. SSaC specs declare only business parameters; infrastructure parameters go in x- extensions only.

```yaml
/courses:
  get:
    operationId: ListCourses
    x-pagination:
      style: offset           # offset | cursor
      defaultLimit: 20
      maxLimit: 100
    x-sort:
      allowed: [created_at, price, rating]
      default: created_at
      direction: desc          # asc | desc
    x-filter:
      allowed: [category, level, instructor_id]
    x-include:
      allowed: [instructor_id:users.id]   # FKColumn:RefTable.RefColumn
```

### x-pagination

| Field | Type | Description |
|---|---|---|
| `style` | string | `offset` (Limit/Offset) or `cursor` (cursor-based) |
| `defaultLimit` | int | Default page size |
| `maxLimit` | int | Maximum page size |

### x-sort

| Field | Type | Description |
|---|---|---|
| `allowed` | string[] | Sortable columns (snake_case) |
| `default` | string | Default sort column |
| `direction` | string | `asc` or `desc` |

### x-filter

| Field | Type | Description |
|---|---|---|
| `allowed` | string[] | Filterable columns (snake_case) |

### x-include

| Field | Type | Description |
|---|---|---|
| `allowed` | string[] | Forward FK includes. Format: `FKColumn:RefTable.RefColumn` |

Syntax (single format only):
- `instructor_id:users.id` — courses.instructor_id -> users.id FK relation to include User
- Runtime include name: Remove `_id` from FK column (`instructor_id` -> `instructor`)
- Reverse FK (1:N) not supported — use separate endpoints

### x- Extension Codegen Effects

- SSaC: Operations with x- extensions get `opts QueryOpts` parameter auto-added to model methods
- SSaC: `:many` + x-pagination -> return type `([]T, int, error)` (includes total count)
- STML: `data-paginate` -> `useState(page, limit)` + prev/next buttons
- STML: `data-sort` -> `useState(sortBy, sortDir)` + toggle buttons
- STML: `data-filter` -> `useState(filters)` + filter inputs

## sqlc Query Rules

```sql
-- name: FindByID :one
SELECT * FROM courses WHERE id = $1;

-- name: List :many
SELECT * FROM courses ORDER BY created_at DESC;

-- name: Create :one
INSERT INTO courses (title, price, instructor_id)
VALUES ($1, $2, $3) RETURNING *;

-- name: SoftDelete :exec
UPDATE courses SET deleted_at = NOW() WHERE id = $1;
```

| Cardinality | SSaC @result Type | Codegen Return |
|---|---|---|
| `:one` | `*Type` | `(*Course, error)` |
| `:many` | `[]Type` | `([]Course, error)` |
| `:exec` | (none) | `error` |

Model name derived from sqlc query filename: `courses.sql` -> `Course`
Singularization rules: `ies`->`y`, `sses`->`ss`, `xes`->`x`, otherwise remove trailing `s`

## model/*.go Rules

Defines reference targets for SSaC `@component`.

```go
// model/notification.go
package model

// NotificationService is a notification component.
type NotificationService interface {
    Send(userID int64, message string) error
}
```

- `type XxxInterface interface` -> referenceable via `@component xxx`
- Structs with `// @dto` comment -> skip DDL table matching (for pure DTOs like Token, Refund)

## Gherkin Scenario — Cross-Endpoint Test Declarations

`scenario/*.feature` files declare cross-endpoint business scenarios and invariants using a constrained Gherkin syntax (fixed patterns, machine-parseable).

### Tags

| Tag | Meaning | Hurl Output |
|---|---|---|
| `@scenario` | Business scenario | `scenario-{feature}.hurl` |
| `@invariant` | Invariant verification | `invariant-{feature}.hurl` |

### Action Steps (Given/When/Then/And)

```
METHOD operationId {JSON} → result     # request + capture
METHOD operationId {JSON}              # request only
METHOD operationId → result            # no-body request + capture
METHOD operationId                     # no-body request only
```

- `METHOD`: `GET`, `POST`, `PUT`, `DELETE`
- `operationId`: OpenAPI operationId (PascalCase)
- `{JSON}`: Request parameters. Unquoted `var.Field` = variable reference
- `→ result`: Capture response as variable. `→ token` auto-injects Authorization header

### Assertion Steps (Then/And)

```
status == CODE                         # HTTP status code
response.field exists                  # field existence
response.field == value                # value equality
response.array contains var.Field      # array inclusion
response.array excludes var.Field      # array exclusion
response.array count > N               # array size
```

### Example

```gherkin
@scenario
Feature: Instructor creates and publishes a course

  Scenario: Full course lifecycle
    Given POST Register {"Email": "inst@test.com", "Password": "Pass1234!", "Name": "Instructor"} → user
    And POST Login {"Email": "inst@test.com", "Password": "Pass1234!"} → token
    When POST CreateCourse {"Title": "Go 101", "Category": "dev", "Level": "beginner", "Price": 10000} → course
    And PUT PublishCourse {"CourseID": course.ID}
    Then GET ListCourses → courses
    And response.courses contains course.ID
    And status == 200
```

## SSOT Connection Map

```
         OpenAPI (operationId)
           |               |
    SSaC (funcName)    STML (data-fetch/action)
      |    |    |          |             |
  DDL  Func  States   Policy    Scenario (.feature)
      |
  sqlc queries (model.method)
```

### Name Matching Rules

| Source | Target | Matching |
|---|---|---|
| stateDiagram transition event | SSaC funcName / OpenAPI operationId | Identical (PascalCase) |
| SSaC function name | OpenAPI operationId | Identical (PascalCase) |
| STML data-fetch/action | OpenAPI operationId | Identical (PascalCase) |
| SSaC @model Model | DDL table name | PascalCase -> snake_case + plural (`Course` -> `courses`) |
| SSaC @model .Method | sqlc query `-- name:` | Identical (`FindByID` = `FindByID`) |
| SSaC @func pkg.name | Func spec @func name | Identical (`hashPassword` = `hashPassword`) |
| x-sort/filter allowed | DDL column name | Identical snake_case |
| x-include allowed | DDL FK relation | `FKColumn:RefTable.RefColumn` -> DDL FK mapping |

## fullend Cross-Validation Rules

After individual tools (ssac validate, stml validate) run their own checks, fullend catches cross-layer mismatches:

| Rule | Validation | Level |
|---|---|---|
| x-sort <-> DDL | Column exists in table | ERROR |
| x-sort <-> DDL index | Column has an index | WARNING |
| x-filter <-> DDL | Column exists in table | ERROR |
| x-include <-> DDL FK | Tables connected by FK relation | WARNING |
| SSaC @result <-> DDL | Result type has corresponding table | WARNING |
| SSaC @param <-> DDL | Parameter has corresponding column | WARNING |
| SSaC funcName -> operationId | SSaC function has corresponding operationId | ERROR |
| operationId -> SSaC funcName | operationId has corresponding SSaC function | WARNING |
| States transition -> SSaC | Transition event has corresponding SSaC function | ERROR |
| States transition -> OpenAPI | Transition event has corresponding operationId | ERROR |
| SSaC guard state -> States | Referenced stateDiagram exists | ERROR |
| States transition -> SSaC guard state | Function with transition has guard state | WARNING |
| guard state field -> DDL | State field exists as DDL column | ERROR |
| Policy <-> SSaC authorize | SSaC authorize (action, resource) -> Rego allow rule exists | WARNING |
| Policy <-> SSaC authorize | Rego allow (action, resource) -> SSaC authorize exists | WARNING |
| Policy @ownership -> DDL | @ownership table.column exists in DDL | ERROR |
| Policy @ownership via -> DDL | via join table.fk exists in DDL | ERROR |
| Policy <-> States | Transition event with authorize -> Rego allow rule exists | WARNING |
| Scenario -> OpenAPI operationId | Scenario step operationId exists in OpenAPI | ERROR |
| Scenario -> OpenAPI method | Scenario step METHOD matches OpenAPI method | ERROR |
| Scenario -> OpenAPI fields | Scenario JSON fields exist in request schema | ERROR |
| Scenario -> States | Scenario step order follows state transitions | WARNING |
| Func -> SSaC @func | @func reference has matching implementation | ERROR |
| Func body | Function body is not a TODO stub | WARNING |
| Func param count | @param count = Request field count | ERROR |
| Func param type | i-th @param type (DDL/OpenAPI) = i-th Request field type | ERROR |
| Func result/response | @result exists ↔ Response fields exist | ERROR/WARNING |
| Func source var | @param source variable defined in prior @result | WARNING |

## Mermaid stateDiagram — State Transition Declarations

`states/*.md` files declare resource state transitions using Mermaid stateDiagram.

### Syntax

```markdown
# CourseState

​```mermaid
stateDiagram-v2
    [*] --> unpublished
    unpublished --> published: PublishCourse
    published --> deleted: DeleteCourse
    unpublished --> deleted: DeleteCourse
​```
```

### Rules

- Filename = stateDiagram ID (e.g., `course.md` -> `course`)
- Transition label = SSaC function name = OpenAPI operationId (PascalCase)
- `[*]` -> initial state (must match DDL DEFAULT value)
- One stateDiagram per file

### Usage in SSaC

```go
// @sequence guard state course
// @param course.Published
```

- `course`: stateDiagram ID (`states/course.md`)
- `course.Published`: State field from previous @result variable
- Function name is used as transition event (PublishCourse function -> PublishCourse transition)

### Codegen Output

```go
// guard state -> 409 Conflict if transition not allowed
if !coursestate.CanTransition(course.Published, "PublishCourse") {
    c.JSON(http.StatusConflict, gin.H{"error": "state transition not allowed"})
    return
}
```

State machine package (`states/<id>state/<id>state.go`) is auto-generated by fullend gen.

## OPA Rego — Authorization Policy Declarations

`policy/*.rego` files declare authorization policies using OPA Rego. fullend parses, cross-validates, and auto-generates an OPA Go SDK-based Authorizer implementation.

### Input Schema (Fixed)

| Field | Source | Description |
|---|---|---|
| `input.user.id` | `CurrentUser.UserID` | Authenticated user ID |
| `input.user.role` | `CurrentUser.Role` | User role |
| `input.action` | SSaC `@action` | Action to perform |
| `input.resource` | SSaC `@resource` | Target resource |
| `input.resource_id` | SSaC `@id` | Resource identifier |
| `input.resource_owner` | @ownership -> DB lookup | Resource owner ID |

### 5 Allow Patterns

```rego
# 1. Unconditional (authenticated only)
allow if { input.action == "create"; input.resource == "course" }

# 2. Role-based
allow if { input.action == "create"; input.resource == "course"; input.user.role == "instructor" }

# 3. Owner-based
allow if { input.action == "update"; input.resource == "course"; input.user.id == input.resource_owner }

# 4. Role + Owner
allow if { input.action == "delete"; input.resource == "course"; input.user.role == "instructor"; input.user.id == input.resource_owner }

# 5. Multiple actions (set)
allow if { input.action in {"update", "delete", "publish"}; input.resource == "course"; input.user.id == input.resource_owner }
```

### @ownership Annotations

Declare owner lookup methods at the top of Rego files:

```rego
# @ownership course: courses.instructor_id
# @ownership lesson: courses.instructor_id via lessons.course_id
# @ownership review: reviews.user_id
```

| Format | Meaning |
|---|---|
| `resource: table.column` | Direct lookup |
| `resource: table.column via join_table.fk` | JOIN lookup |

### Relationship with SSaC authorize

SSaC `authorize` sequence `@action`/`@resource` maps to Rego `allow` rule `input.action`/`input.resource`. Independent concern from stateDiagram (state vs permission).

### Codegen Output

`authz/authz.go` — OPA Go SDK-based Authorizer implementation (evaluates embedded .rego file)
`authz/<name>.rego` — Copied from specs (for go:embed)

## Runtime Testing (Hurl)

`fullend gen` auto-generates Hurl tests from OpenAPI specs and Gherkin scenarios.

```bash
# After starting the server:
hurl --test --variable host=http://localhost:8080 artifacts/<project>/tests/*.hurl
```

Generated tests include:
- **smoke.hurl** — OpenAPI endpoint smoke tests (auto-generated)
- **scenario-*.hurl** — Business scenario tests (from .feature files)
- **invariant-*.hurl** — Cross-endpoint invariant tests (from .feature files)

## fullend CLI

```bash
fullend validate [--skip kind,...] <specs-dir>                 # Individual validation + cross-validation
fullend gen      [--skip kind,...] <specs-dir> <artifacts-dir> # validate -> codegen + Hurl tests + state machines + OPA Authorizer
fullend status   <specs-dir>                                   # SSOT summary
```

9 required SSOTs (OpenAPI, DDL, SSaC, Model, STML, States, Policy, Scenario, Terraform) cause an ERROR if missing. Func is optional (detected only when `func/` exists). Use `--skip` to explicitly exclude:

```bash
fullend validate --skip states,terraform,scenario specs/my-project
```

Skip kinds: `openapi`, `ddl`, `ssac`, `model`, `stml`, `states`, `policy`, `scenario`, `func`, `terraform`
