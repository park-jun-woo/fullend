# Phase 027: gin 프레임워크 전환 + 스모크 테스트 통과

## 목표

생성 백엔드를 `net/http` → `gin` 프레임워크로 전환한다.
`fullend gen` → `go build` → `hurl --test smoke.hurl` 전 과정이 수동 패치 없이 통과.

## 선행 조건

- SSaC 수정지시서012 완료 (✅ call codegen 수정)
- SSaC 수정지시서013 완료 (✅ gin 템플릿 전환)

## 변경 항목

### A. SSaC 의존 (수정지시서013)

SSaC generator의 Go 코드 생성 템플릿을 gin 대상으로 변경.
- 함수 시그니처: `(w http.ResponseWriter, r *http.Request)` → `(c *gin.Context)`
- 경로 파라미터: 함수 인자 → `c.Param()` + 타입 변환
- 요청 파싱: `json.NewDecoder` → `c.ShouldBindJSON`
- 에러 응답: `http.Error` → `c.JSON(status, gin.H{"error": msg})`
- 성공 응답: `json.NewEncoder` → `c.JSON(200, gin.H{...})`
- authorize: `currentUser` → `c.MustGet("currentUser")`

### B. fullend 미들웨어 패키지 (`pkg/middleware/`)

#### 1. `pkg/middleware/jwt.go` — JWT 미들웨어

fullend가 기본 제공하는 gin 미들웨어. `pkg/auth`와 동일하게 패키지로 제공.

```go
package middleware

type CurrentUser struct {
    UserID int64
    Email  string
    Role   string
}

func JWT(secret string) gin.HandlerFunc {
    // Authorization: Bearer <token> → auth.VerifyToken → c.Set("currentUser", ...)
    // abort하지 않고 빈 CurrentUser 세팅 — authorize 시퀀스가 권한 검사 담당
}
```

생성 프로젝트의 `model/auth.go`에서 타입 앨리어스로 연결:
```go
type CurrentUser = middleware.CurrentUser
```

### C. fullend glue-gen 변경

#### 1. 라우터: ServeMux → gin.Engine

```go
r := gin.Default()

// OpenAPI securitySchemes에서 bearerAuth 감지 → 미들웨어 자동 연결
auth := r.Group("/")
auth.Use(middleware.JWT("secret"))

r.POST("/login", s.Auth.Login)           // security 없음 → public
auth.POST("/courses", s.Course.CreateCourse)  // security: [{bearerAuth}] → auth 그룹
```

- 경로 파라미터: `{CourseID}` → `:CourseID` (gin 문법)
- path param 파싱/전달 코드 제거 — SSaC codegen이 함수 안에서 `c.Param()` 처리
- **라우트 그룹 결정: OpenAPI `security` 필드가 SSOT** (SSaC authorize 시퀀스 대신)

#### 2. model/auth.go 생성

```go
import "github.com/park-jun-woo/fullend/pkg/middleware"

type CurrentUser = middleware.CurrentUser

type Authorizer interface {
    Check(user *CurrentUser, action, resource string, id interface{}) (bool, error)
}
```

- `CurrentUserFunc` 타입 제거
- `service/auth.go` (DefaultCurrentUser) 더 이상 생성하지 않음

#### 3. main.go 생성

- `http.ListenAndServe` → `r.Run(addr)`
- Authz 초기화 + handler에 `Authz: az` 연결

```go
az, err := authz.New(conn)
server := &service.Server{
    Course: &coursesvc.Handler{
        CourseModel: model.NewCourseModel(conn),
        Authz:       az,
    },
}
r := service.SetupRouter(server)
log.Fatal(r.Run(*addr))
```

#### 4. handler struct 변경

`CurrentUser model.CurrentUserFunc` 필드 제거 (미들웨어가 context에 저장).

#### 5. SSaC codegen 리시버 패칭

```go
// 변경 전
src = strings.ReplaceAll(src, "authz.Check(currentUser,", rcv+".Authz.Check("+rcv+".CurrentUser(r),")

// 변경 후 — currentUser는 SSaC codegen이 c.MustGet으로 이미 추출
src = strings.ReplaceAll(src, "authz.Check(currentUser,", rcv+".Authz.Check(currentUser,")
```

#### 6. QueryOpts gin 전환

`ParseQueryOpts(r *http.Request, ...)` → `ParseQueryOpts(c *gin.Context, ...)`
`r.URL.Query().Get(...)` → `c.Query(...)`

### D. dummy-lesson spec 업데이트

#### 1. import 경로 추가 (수정지시서012 반영)

```go
import (
    "net/http"
    _ "github.com/park-jun-woo/fullend/pkg/auth"
)
```

#### 2. @result 필드명 명시

```go
// @result token IssueTokenResponse.AccessToken
// @result hashedPassword HashPasswordResponse.HashedPassword
```

#### 3. @param -> 매핑

```go
// @param user.ID -> UserID
// @param user.Email -> Email
// @param user.Role -> Role
```

### E. dummy-study spec 동일 적용

- `service/login.go`: import 추가

### F. go.mod 의존성

생성되는 `go.mod`에 `github.com/gin-gonic/gin` 추가.
fullend 자체 `go.mod`에도 gin 의존성 추가 (`pkg/middleware`용).

## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| `pkg/middleware/jwt.go` | **신규** — JWT gin 미들웨어 + CurrentUser 타입 |
| `artifacts/internal/gluegen/domain.go` | gin 라우터, OpenAPI security 기반 그룹, model/auth.go, main.go, handler struct |
| `artifacts/internal/gluegen/gluegen.go` | 리시버 패칭 authz 단순화, QueryOpts r→c |
| `artifacts/internal/gluegen/queryopts.go` | ParseQueryOpts gin.Context 전환 |
| `specs/dummy-lesson/service/auth/login.go` | import + @result 필드명 + @param 매핑 |
| `specs/dummy-lesson/service/auth/register.go` | import + @result 필드명 |
| `specs/dummy-study/service/login.go` | import 추가 |
| `go.mod` | gin 의존성 추가 |

## 검증 방법

```bash
# 1. fullend 빌드
go build ./artifacts/cmd/fullend/

# 2. dummy-lesson 코드 산출
fullend gen --skip terraform specs/dummy-lesson /tmp/gen-lesson

# 3. 생성된 서버 빌드 (수동 패치 없이)
cd /tmp/gen-lesson/backend && go build -o server ./cmd/main.go

# 4. DB 초기화 + 서버 기동
./server -addr :18080 -dsn "postgres://..."

# 5. Hurl 스모크 테스트
hurl --test --variable host=http://localhost:18080 /tmp/gen-lesson/tests/smoke.hurl
```

전체 Hurl 시나리오 통과가 목표.
