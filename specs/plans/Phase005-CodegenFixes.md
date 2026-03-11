# ✅ Phase 005: 코드젠 빌드 에러 수정 (authz·state·claims·gluegen)

## 목표

Phase 004 스모크 테스트에서 발견된 fullend 쪽 코드젠 빌드 에러를 수정한다.
SSaC 수정지시서002(guard-style @call, 구조체 필드 매핑, validateCurrentUserType 제거)와 병행하여, fullend 측 수정으로 `go build` 통과를 달성한다.

## 배경

`fullend gen` → `go build` 시 다음 에러 발생:

| 에러 | 원인 | 수정 대상 |
|------|------|----------|
| `authz.Input{}` undefined | SSaC codegen이 `authz.Input{}` 생성, fullend authz-gen은 `id interface{}` 사용 | fullend authz-gen |
| `coursestate.Input{}` undefined | SSaC codegen이 `coursestate.Input{}` 생성, fullend state-gen은 `CanTransition(bool, string)` 사용 | fullend state-gen |
| `currentUser.ID` undefined | `middleware.CurrentUser.UserID` vs spec `currentUser.ID` 필드명 불일치 | fullend.yaml claims 설정 + codegen |
| `Authz.Check returns 2 values` | `Check` 반환 `(bool, error)`, SSaC codegen은 `err :=` 단일 할당 | authz-gen |
| `authz` bare import | ✅ Phase 004에서 수정 완료 |
| `"model"` bare import | ✅ Phase 004에서 수정 완료 |
| `@call` 모델 오수집 | ✅ Phase 004에서 수정 완료 |

## 아키텍처 결정 사항

### CurrentUser는 fullend.yaml claims 설정이 소유한다

- **SSaC는 `currentUser`를 예약 소스로만 인식** — 타입 정의·검증은 SSaC의 책임이 아님
- **SSaC의 `validateCurrentUserType` 제거** (수정지시서002 문제3)
- **fullend.yaml `backend.auth.claims`** 섹션에서 필드-to-claim 매핑을 명시적으로 선언
- **fullend codegen**이 claims 설정으로부터 `CurrentUser` 타입 + 미들웨어 코드를 생성
- **fullend crosscheck**가 SSaC spec의 `currentUser.X` 필드가 claims에 정의되어 있는지 검증
- `specs/dummy-lesson/model/current_user.go`는 불필요 — fullend.yaml에서 대체

### fullend.yaml claims 설정 예시

```yaml
backend:
  auth:
    secret_env: JWT_SECRET
    claims:
      ID: user_id        # claim "user_id" → CurrentUser.ID (int64)
      Email: email        # claim "email"   → CurrentUser.Email (string)
      Role: role          # claim "role"    → CurrentUser.Role (string)
```

## SSaC가 `currentUser`를 소비하는 전체 경로

fullend가 CurrentUser를 책임지므로, SSaC가 currentUser를 소비하는 모든 지점을 fullend가 검증·생성해야 한다.

### SSaC codegen이 생성하는 코드 (3곳)

| 생성 코드 | SSaC 위치 | fullend 대응 |
|-----------|-----------|-------------|
| `currentUser := c.MustGet("currentUser").(*model.CurrentUser)` | `go_templates.go` "currentUser" 템플릿 | **codegen**: claims 기반 `model.CurrentUser` 타입 생성 |
| `currentUser.ID`, `currentUser.Email` 등 필드 접근 | `go_target.go` `argToCode()` — `a.Source == "currentUser"` | **crosscheck**: 필드가 claims에 정의됐는지 검증 |
| `"model"` import 추가 | `go_target.go` `needsCU == true` → `seen["model"]` | **gluegen**: bare `"model"` → full path 변환 (✅ Phase 004 완료) |

### SSaC가 currentUser 사용을 감지하는 조건 (`needsCurrentUser()`)

1. `seq.Type == "auth"` (모든 @auth 시퀀스)
2. `arg.Source == "currentUser"` (인자로 currentUser.X 전달)
3. `input`이 `"currentUser."` prefix 사용 (@auth, @state의 input)

### fullend가 체크해야 할 항목

| # | 체크 항목 | 체크 대상 | 에러 수준 | 구현 위치 |
|---|----------|-----------|----------|-----------|
| 1 | claims 설정 존재 | `currentUser`를 사용하는 SSaC spec이 있으면 `backend.auth.claims`가 반드시 존재해야 함 | ERROR | `crosscheck/claims.go` |
| 2 | 필드 정의 확인 | SSaC spec의 모든 `currentUser.X` 필드가 claims에 정의돼 있는지 | ERROR | `crosscheck/claims.go` |
| 3 | @auth 정합성 | `@auth` 시퀀스 사용 시 claims 설정 존재 필수 (currentUser가 scope에 주입되므로) | ERROR | `crosscheck/claims.go` |
| 4 | 생성 타입 일치 | codegen된 `model.CurrentUser` 구조체 필드가 claims 키와 1:1 매핑 | (codegen 보장) | `gluegen/` |
| 5 | 미들웨어 일치 | codegen된 미들웨어가 claims claim-key를 정확히 추출 | (codegen 보장) | `gluegen/` |
| 6 | IssueToken 정합성 | `pkg/auth/issue_token.go`의 claims 필드가 fullend.yaml claims의 claim-key와 매핑 | WARNING | `crosscheck/claims.go` |

### dummy-lesson에서 사용되는 currentUser 필드

| 필드 | 사용 파일 | 사용 횟수 |
|------|----------|-----------|
| `currentUser.ID` | create_course, create_review, enroll_course(×2), list_my_enrollments, list_my_payments, FindByCourseAndUser(×2) | 9 |
| `currentUser.Email` | (현재 사용 없음) | 0 |
| `currentUser.Role` | (현재 사용 없음, @auth가 내부적으로 사용) | 0 |

## 변경 항목

### A. fullend.yaml claims 설정 파싱 (`internal/projectconfig/`)

`fullend.yaml`에 `backend.auth` 섹션을 추가하고 파싱한다.

```go
type Backend struct {
    Lang       string   `yaml:"lang"`
    Framework  string   `yaml:"framework"`
    Module     string   `yaml:"module"`
    Middleware []string `yaml:"middleware"`
    Auth       *Auth    `yaml:"auth"`       // 추가
}

type Auth struct {
    SecretEnv string            `yaml:"secret_env"`
    Claims    map[string]string `yaml:"claims"` // FieldName → claim key
}
```

### B. CurrentUser 타입 codegen (`internal/gluegen/`)

claims 설정으로부터 `internal/model/auth.go`에 `CurrentUser` 구조체를 생성한다.

현재 하드코딩:
```go
// domain.go generateAuthStubWithDomains()
type CurrentUser = middleware.CurrentUser  // ← pkg/middleware에 의존
```

변경 후:
```go
// claims 기반 생성
type CurrentUser struct {
    ID    int64   // ← claims["ID"] 존재, IssueToken의 UserID claim → int64 추론
    Email string  // ← claims["Email"]
    Role  string  // ← claims["Role"]
}
```

`pkg/middleware/bearerauth.go`는 더 이상 `CurrentUser` 타입을 정의하지 않는다. 미들웨어도 claims 설정 기반으로 생성되거나, 생성된 프로젝트의 `internal/middleware/`에 codegen된다.

### C. 미들웨어 codegen (`internal/gluegen/`)

claims 설정 + `secret_env`로부터 JWT 미들웨어를 생성한다.

```go
// 생성될 internal/middleware/bearerauth.go
func BearerAuth(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        // ... JWT 파싱 ...
        c.Set("currentUser", &model.CurrentUser{
            ID:    claims["user_id"].(int64),    // ← claims 매핑에서 파생
            Email: claims["email"].(string),
            Role:  claims["role"].(string),
        })
        c.Next()
    }
}
```

### D. fullend crosscheck — `currentUser ↔ claims` 검증 (`internal/crosscheck/claims.go`)

새 crosscheck 규칙:

```go
// CheckClaims validates currentUser usage against fullend.yaml claims config.
func CheckClaims(
    serviceFuncs []ssacparser.ServiceFunc,
    claims map[string]string, // nil if no auth config
) []CrossError {
    // 1. SSaC spec 전체에서 currentUser 사용 필드 수집
    usedFields := collectCurrentUserFields(serviceFuncs)
    // → {"ID": ["create_course.go:4", "enroll_course.go:8", ...], ...}

    // 2. currentUser를 사용하는데 claims 설정이 없으면 ERROR
    if len(usedFields) > 0 && claims == nil {
        // ERROR: "currentUser를 사용하지만 fullend.yaml에 backend.auth.claims가 정의되지 않았습니다"
    }

    // 3. 각 필드가 claims에 정의돼 있는지 검증
    for field, locations := range usedFields {
        if _, ok := claims[field]; !ok {
            // ERROR: "currentUser.{field} 사용 ({locations}) — claims에 미정의"
        }
    }
}

// collectCurrentUserFields scans all SSaC sequences for currentUser.X references.
func collectCurrentUserFields(funcs []ssacparser.ServiceFunc) map[string][]string {
    // 수집 대상:
    // 1. seq.Args에서 a.Source == "currentUser" → a.Field
    // 2. seq.Inputs에서 "currentUser.X" prefix → X
    // 3. @auth 시퀀스는 currentUser 암묵 사용 (scope 주입) → claims 존재 필수
}
```

수집 대상 상세:

| SSaC 소비 지점 | 코드 위치 | 수집 방법 |
|---------------|-----------|-----------|
| `@get/@post/@put/@delete` Args | `seq.Args[i].Source == "currentUser"` | `arg.Field` 수집 |
| `@auth` Inputs | `seq.Inputs["id"]` 등에서 `"currentUser.X"` | prefix split |
| `@state` Inputs | `seq.Inputs["status"]` 등에서 `"currentUser.X"` | prefix split |
| `@auth` 암묵 사용 | `seq.Type == "auth"` | claims 존재 필수 (필드 수집 불필요) |

### E. crosscheck 통합 (`internal/crosscheck/crosscheck.go`)

```go
// crosscheck.go Run() 에 추가
// Claims ↔ SSaC currentUser
if input.ServiceFuncs != nil {
    errs = append(errs, CheckClaims(input.ServiceFuncs, input.Claims)...)
}
```

`CrossValidateInput`에 `Claims map[string]string` 필드 추가.

### F. fullend authz-gen — `authz.Input` 타입 생성 + `Check` 시그니처 조정

SSaC codegen이 `authz.Input{ID: courseID}` 형태로 생성하므로, fullend authz-gen이 `Input` 구조체를 생성하고 `Check` 시그니처를 맞춘다.

```go
// 생성할 타입
type Input struct {
    ID interface{}
}

// Check 시그니처 조정: (bool, error) → error
func (a *OPAAuthorizer) Check(user *model.CurrentUser, action, resource string, input Input) error {
    // ... OPA eval ...
    if !allowed {
        return fmt.Errorf("forbidden")
    }
    return nil
}
```

반환값을 `error`만으로 변경하면 SSaC codegen의 `err :=` 패턴과 일치한다.

`model/auth.go`의 `Authorizer` 인터페이스도 동기화:
```go
type Authorizer interface {
    Check(user *CurrentUser, action, resource string, input authz.Input) error
}
```

### G. fullend state-gen — `coursestate.Input` 타입 생성 + `CanTransition` 시그니처 조정

SSaC codegen이 `coursestate.Input{Status: course.Published}` 형태로 생성하므로, fullend state-gen이 `Input` 구조체를 생성하고 `CanTransition` 시그니처를 맞춘다.

```go
// 생성할 타입 (state diagram의 inputs 키에서 파생)
type Input struct {
    Status interface{}
}

// CanTransition 시그니처 조정
func CanTransition(input Input, event string) bool {
    current := resolveState(input.Status)
    _, ok := transitions[transitionKey{from: current, event: event}]
    return ok
}
```

### H. `specs/dummy-lesson/model/current_user.go` 제거

fullend.yaml claims 설정이 CurrentUser를 대체하므로, SSaC model에서 CurrentUser 정의를 제거한다.

### I. `pkg/middleware/bearerauth.go` 역할 축소

기존: `CurrentUser` 타입 정의 + 미들웨어 함수 제공.
변경: `CurrentUser` 타입 정의 제거, 미들웨어는 codegen 기반으로 이전.

`pkg/middleware/bearerauth.go`는 범용 JWT 파싱 유틸리티로 축소하거나, claims 기반 codegen이 완전히 대체하면 삭제 후보.

## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| `internal/projectconfig/projectconfig.go` | `Auth` 구조체 + claims 파싱 추가 |
| `specs/dummy-lesson/fullend.yaml` | `backend.auth` 섹션 추가 |
| `internal/crosscheck/claims.go` | **신규** — `CheckClaims()` + `collectCurrentUserFields()` |
| `internal/crosscheck/crosscheck.go` | `CrossValidateInput`에 `Claims` 필드, `Run()`에 `CheckClaims` 호출 추가 |
| `internal/crosscheck/types.go` | (필요 시 CrossError 관련 추가) |
| `internal/gluegen/domain.go` | `generateAuthStubWithDomains` → claims 기반 CurrentUser 구조체 생성 |
| `internal/gluegen/gluegen.go` | authz-gen: `Input` 타입 + `Check` 반환 `error`만 |
| `internal/gluegen/gluegen.go` | state-gen: `Input` 타입 + `CanTransition(Input, string)` |
| `internal/gluegen/` | 미들웨어 codegen (claims 기반 JWT 미들웨어 생성) |
| `pkg/middleware/bearerauth.go` | `CurrentUser` 타입 제거, 역할 축소 |
| `specs/dummy-lesson/model/current_user.go` | 삭제 |
| `internal/orchestrator/` | claims config를 crosscheck에 전달하는 배관 추가 |

## 의존성

- SSaC 수정지시서002 완료 (guard-style @call, 구조체 필드 매핑, validateCurrentUserType 제거)
- Phase 004 진행 중 (validate 통과 완료, gen 완료, build 에러 수정 중)

## 검증 방법

```bash
# 1. fullend 빌드
go build ./cmd/fullend/

# 2. fullend 테스트
go test ./... -count=1

# 3. crosscheck 검증 — claims 미정의 시 ERROR
# fullend.yaml에서 backend.auth.claims 제거 후 validate → currentUser 사용 ERROR 확인

# 4. crosscheck 검증 — 필드 미정의 시 ERROR
# claims에서 ID 제거 후 validate → currentUser.ID 미정의 ERROR 확인

# 5. gen + build
./fullend gen specs/dummy-lesson artifacts/dummy-lesson
cd artifacts/dummy-lesson/backend && go mod tidy && go build ./...
```
