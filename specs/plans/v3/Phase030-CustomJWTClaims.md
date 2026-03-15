✅ 완료

# Phase030: 커스텀 JWT Claims 지원 (BUG021)

## 목표

`fullend.yaml`에 정의한 모든 claims(표준 + 커스텀)로부터 타입 안전한 auth/middleware 코드를 생성한다.

설계 원칙:
1. **SSOT 1회 정의** — claims 매핑은 `fullend.yaml`에만 존재
2. **타입 확정 생성** — `map[string]interface{}` 없이, claims 기반 struct를 생성
3. **`//fullend:preserve` 허용** — 생성된 auth/middleware를 사용자가 커스터마이즈 가능
4. **pkg/auth 분리** — claims 의존 함수(IssueToken, VerifyToken, RefreshToken)는 gluegen 산출물로 이동, 범용 유틸(HashPassword, VerifyPassword 등)은 `pkg/auth`에 유지

## 동기

현재 `pkg/auth`의 `IssueToken`/`VerifyToken`이 `UserID`, `Email`, `Role` 3필드를 하드코딩.
미들웨어 codegen의 `claimToVerifyField` 맵도 3개만 등록.
커스텀 claim(`OrgID: org_id`)을 추가하면 `// OrgID: unknown claim key "org_id"` 주석 처리.

## 설계

### 1단계: `fullend.yaml` auth 섹션 확장

```yaml
auth:
  type: jwt               # 인증 방식 지정 (현재 jwt만 지원, 향후 확장 가능)
  secret_env: JWT_SECRET
  claims:
    ID: user_id:int64
    Email: email:string
    Role: role:string
    OrgID: org_id:int64   # 커스텀 claim도 동일한 형식
  roles: [admin, member]
```

- `type` 필드 추가: 인증 방식 명시. `jwt` 외 값은 validate ERROR.
- claims value 형식: `claim_key:go_type`. 타입 미지정 시 `string` 기본값.

### 2단계: `projectconfig` 파서 수정

```go
type ClaimDef struct {
    Key    string // JWT claim key (예: "org_id")
    GoType string // Go type (예: "int64"), 기본값 "string"
}

type Auth struct {
    Type      string              `yaml:"type"`       // "jwt" (필수)
    SecretEnv string              `yaml:"secret_env"`
    Claims    map[string]ClaimDef // FieldName → ClaimDef (파싱 후 변환)
    RawClaims map[string]string   `yaml:"claims"`     // YAML 원본
    Roles     []string            `yaml:"roles"`
}
```

`Load()` 후처리에서 `RawClaims` → `Claims` 변환.

### 3단계: validate — claims 검증

`ProjectConfig.Validate()` 확장:

| 규칙 | 레벨 | 조건 |
|---|---|---|
| auth.type 필수 | ERROR | auth 섹션이 있으면 type 필수 |
| auth.type 지원 확인 | ERROR | `jwt` 외 값 |
| claims 필수 | ERROR | auth 섹션이 있으면 claims 최소 1개 필수 |
| 타입 제한 | ERROR | GoType이 `string`, `int64`, `bool` 외 |
| JWT 예약 키 충돌 | ERROR | claim key가 `exp`, `iat`, `sub`, `iss`, `aud`, `nbf`, `jti` |
| claim key 중복 | ERROR | 서로 다른 필드가 동일한 claim key 사용 |
| 타입 미지정 | WARNING | 타입 미지정 시 string 기본값 적용 안내 |

claims 필수화로 `domain.go:102-105`의 깨진 fallback(`pkg/middleware.CurrentUser` 참조) 제거.
auth 섹션이 있으면 반드시 claims가 있으므로 항상 inline CurrentUser struct 생성.

### 4단계: `pkg/auth` — claims 의존 함수 제거, 범용 유틸 유지

**삭제** (gluegen 산출물로 이동):
- `pkg/auth/issue_token.go` — claims 구조 의존
- `pkg/auth/verify_token.go` — claims 구조 의존
- `pkg/auth/refresh_token.go` — claims 구조 의존

**유지** (`pkg/auth`에 잔류):
- `pkg/auth/hash_password.go` — bcrypt, claims 무관
- `pkg/auth/verify_password.go` — bcrypt, claims 무관
- `pkg/auth/generate_reset_token.go` — crypto/rand, claims 무관

유지되는 함수들은 기존과 동일하게 `FullendPkgSpecs`로 수집되어 `@call auth.HashPassword` 등에서 참조 가능.

### 5단계: gluegen — `internal/auth/` 생성

`fullend.yaml` claims로부터 타입 확정된 JWT 함수를 프로젝트의 `internal/auth/`에 생성한다.

**`internal/auth/issue_token.go`** (generated):

```go
//fullend:gen ssot=fullend.yaml contract=abc1234
package auth

import (
    "os"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type IssueTokenRequest struct {
    ID    int64
    Email string
    Role  string
    OrgID int64
}

type IssueTokenResponse struct {
    AccessToken string `json:"access_token"`
}

//fullend:gen ssot=fullend.yaml contract=abc1234
func IssueToken(req IssueTokenRequest) (IssueTokenResponse, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        secret = "secret"
    }
    claims := jwt.MapClaims{
        "user_id": req.ID,
        "email":   req.Email,
        "role":    req.Role,
        "org_id":  req.OrgID,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString([]byte(secret))
    return IssueTokenResponse{AccessToken: signed}, err
}
```

**`internal/auth/verify_token.go`** (generated):

```go
//fullend:gen ssot=fullend.yaml contract=def5678
package auth

import (
    "fmt"
    "github.com/golang-jwt/jwt/v5"
)

type VerifyTokenRequest struct {
    Token  string
    Secret string
}

type VerifyTokenResponse struct {
    ID    int64
    Email string
    Role  string
    OrgID int64
}

//fullend:gen ssot=fullend.yaml contract=def5678
func VerifyToken(req VerifyTokenRequest) (VerifyTokenResponse, error) {
    token, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
        return []byte(req.Secret), nil
    })
    if err != nil {
        return VerifyTokenResponse{}, err
    }
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok || !token.Valid {
        return VerifyTokenResponse{}, fmt.Errorf("invalid token")
    }
    userID, _ := claims["user_id"].(float64)
    email, _ := claims["email"].(string)
    role, _ := claims["role"].(string)
    orgID, _ := claims["org_id"].(float64)
    return VerifyTokenResponse{
        ID:    int64(userID),
        Email: email,
        Role:  role,
        OrgID: int64(orgID),
    }, nil
}
```

`//fullend:preserve`로 전환하면 JWT 만료 시간 변경, 서명 알고리즘 교체 등 커스터마이즈 가능.

### 6단계: gluegen — `internal/middleware/` 생성

**`internal/middleware/bearerauth.go`** (generated):

```go
//fullend:gen ssot=fullend.yaml contract=ghi9012
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "{modulePath}/internal/auth"
    "{modulePath}/internal/model"
)

//fullend:gen ssot=fullend.yaml contract=ghi9012
func BearerAuth(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if !strings.HasPrefix(header, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            return
        }
        token := strings.TrimPrefix(header, "Bearer ")
        out, err := auth.VerifyToken(auth.VerifyTokenRequest{Token: token, Secret: secret})
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
            return
        }
        c.Set("currentUser", &model.CurrentUser{
            ID:    out.ID,
            Email: out.Email,
            Role:  out.Role,
            OrgID: out.OrgID,
        })
        c.Next()
    }
}
```

타입 어설션 없이 struct 필드 직접 대입. `//fullend:preserve`로 커스텀 인증 로직 추가 가능.

### 7단계: gluegen — CurrentUser struct 생성

**`internal/gen/gogin/domain.go`의 `generateAuthStubWithDomains`** 수정.

`inferClaimGoType(fieldName)` 삭제, `ClaimDef.GoType` 사용:

```go
for _, field := range fields {
    def := claims[field]
    b.WriteString(fmt.Sprintf("\t%s %s\n", field, def.GoType))
}
```

claims 미정의 fallback (`pkg/middleware.CurrentUser` 참조) 삭제 — 3단계에서 claims 필수화했으므로 도달 불가.

### 8단계: SSaC codegen — import 경로 분기

SSaC `@call auth.IssueToken` 생성 시 import 경로 처리:

현재 SSaC codegen(`go_templates.go`)은 `@call` 패키지명으로 import를 생성한다:
```go
// 현재: 항상 pkg/auth import
"github.com/geul-org/fullend/pkg/auth"
```

변경: `auth.IssueToken`, `auth.VerifyToken`, `auth.RefreshToken`은 generated `internal/auth` 참조, 나머지(`auth.HashPassword` 등)는 기존 `pkg/auth` 참조.

**분기 로직** (SSaC codegen 또는 gogin transform에서):
- claims 의존 함수 목록: `IssueToken`, `VerifyToken`, `RefreshToken`
- `@call auth.{위 함수}` → import `"{modulePath}/internal/auth"`
- `@call auth.{그 외}` → import `"github.com/geul-org/fullend/pkg/auth"`

동일 파일에서 양쪽 모두 사용하면 import alias 필요:
```go
import (
    auth "{modulePath}/internal/auth"           // IssueToken, VerifyToken
    pkgauth "github.com/geul-org/fullend/pkg/auth"  // HashPassword, VerifyPassword
)
```

실제로 `login.ssac`은 `auth.VerifyPassword` + `auth.IssueToken` 둘 다 사용하므로 이 케이스 처리 필수.

**대안**: 같은 패키지명 `auth`이므로, generated `internal/auth/`에 `pkg/auth`의 범용 함수를 re-export하면 import 하나로 해결:

```go
// generated: internal/auth/reexport.go
package auth

import pkgauth "github.com/geul-org/fullend/pkg/auth"

// Re-export pkg/auth utilities for unified import.
var HashPassword = pkgauth.HashPassword
var VerifyPassword = pkgauth.VerifyPassword
var GenerateResetToken = pkgauth.GenerateResetToken

type HashPasswordRequest = pkgauth.HashPasswordRequest
type HashPasswordResponse = pkgauth.HashPasswordResponse
type VerifyPasswordRequest = pkgauth.VerifyPasswordRequest
type VerifyPasswordResponse = pkgauth.VerifyPasswordResponse
type GenerateResetTokenRequest = pkgauth.GenerateResetTokenRequest
type GenerateResetTokenResponse = pkgauth.GenerateResetTokenResponse
```

이러면:
- SSaC codegen의 import 로직 변경 불필요 — 항상 `{modulePath}/internal/auth`만 import
- `auth.HashPassword`, `auth.IssueToken` 모두 같은 import로 해결
- `pkg/auth` import alias 불필요

**이 대안 채택.**

#### import 치환 구현

SSaC 파일은 기존대로 `import "github.com/geul-org/fullend/pkg/auth"` 유지한다.
codegen이 생성한 `.go` 파일에서 `transformSource()` (`internal/gen/gogin/gogin.go`)가 치환:

```go
// Fix auth import: "github.com/geul-org/fullend/pkg/auth" → "{modulePath}/internal/auth"
if strings.Contains(src, "\"github.com/geul-org/fullend/pkg/auth\"") {
    src = strings.ReplaceAll(src, "\"github.com/geul-org/fullend/pkg/auth\"",
        fmt.Sprintf("\"%s/internal/auth\"", modulePath))
}
```

기존 `transformSource()`에 동일한 패턴이 이미 존재 (`authz` → `pkg/authz`, `queue` → `pkg/queue` 등).
이 치환으로 `reexport.go`와 결합되어 `auth.HashPassword`(re-export), `auth.IssueToken`(generated) 모두 단일 import로 해결.

### 9단계: SSaC codegen — IssueToken 호출 코드

SSaC의 `@call auth.IssueToken`은 기존 `@call`과 동일하게 input 인자로 필드 매핑을 지정한다:

```
// @call auth.IssueToken({ID: user.ID, Email: user.Email, Role: user.Role, OrgID: me.OrgID}) → token
```

4단계에서 생성된 `IssueTokenRequest`는 claims 기반 struct이므로, SSaC input이 그대로 struct 초기화 코드가 된다. 기존 `@call` 코드 생성 템플릿이 그대로 동작.

**crosscheck 추가**: `@call auth.IssueToken`의 input field가 `fullend.yaml` claims의 필드명과 일치하는지 검증.

### 10단계: `FullendPkgSpecs` 수집 경로 수정

`findFullendPkgRoot()`가 `pkg/` 전체를 스캔하는데, `issue_token.go`, `verify_token.go`, `refresh_token.go`가 삭제되면 자연히 수집되지 않음. 별도 수정 불필요.

다만 crosscheck `CheckFuncs()`에서 `@call auth.IssueToken` 검증 시:
- 현재: `FullendPkgSpecs`에서 `auth.issueToken` 찾음
- 변경 후: `FullendPkgSpecs`에 없음 → funcspec 검증 실패

**수정**: `CheckFuncs()`에서 claims 의존 함수(`IssueToken`, `VerifyToken`, `RefreshToken`)는 `auth.type: jwt`일 때 built-in으로 간주하여 검증 스킵. 또는 gluegen이 `internal/auth/`에 생성한 funcspec을 `ProjectFuncSpecs`로 수집하도록 확장.

**후자 채택**: generated `internal/auth/*.go`에도 `@func` 어노테이션을 포함시켜, 프로젝트 빌드 후 funcspec으로 인식. 다만 `fullend validate` 시점에는 아직 생성 전이므로, claims 의존 함수는 validate에서 **known built-in 목록**으로 화이트리스트:

```go
var jwtBuiltinFuncs = map[string]bool{
    "auth.issueToken":  true,
    "auth.verifyToken": true,
    "auth.refreshToken": true,
}
```

### 11단계: crosscheck, contract 타입 맞춤

**`internal/crosscheck/types.go`**:
```go
Claims map[string]ClaimDef // map[string]string → map[string]ClaimDef
```

**`internal/crosscheck/claims.go`**:
- `CheckClaims`: `claims[field]` 존재 확인 — 동일 로직, 타입만 변경
- `CheckClaimsRego`: `claimValues[v]` → `claimValues[def.Key]`로 변경

**`internal/orchestrator/validate.go`**:
- claims 추출 타입 맞춤

**`internal/contract/hash.go`**:
- `HashClaims(claims map[string]ClaimDef)` — `k+":"+def.Key+":"+def.GoType`

### 12단계: 정리

- `generateServerStructWithDomains`에서 미사용 claims 파라미터 제거
- `inferClaimGoType` 함수 삭제
- `claimToVerifyField` 맵 삭제
- `domain.go:102-105` 깨진 fallback 삭제
- fullend `go.mod`에서 `golang-jwt/jwt/v5` 의존 유지 (pkg/auth 나머지 함수가 사용하지 않으므로 삭제 가능하나, 테스트용으로 유지)

## 변경 파일

| 파일 | 변경 |
|---|---|
| `internal/projectconfig/projectconfig.go` | `Auth.Type` 추가, `ClaimDef` 타입, `Auth.Claims` 타입 변경, YAML 후처리, `Validate()` 확장 (claims 필수, type 검증) |
| `pkg/auth/issue_token.go` | **삭제** |
| `pkg/auth/verify_token.go` | **삭제** |
| `pkg/auth/refresh_token.go` | **삭제** |
| `pkg/auth/hash_password.go` | 유지 (변경 없음) |
| `pkg/auth/verify_password.go` | 유지 (변경 없음) |
| `pkg/auth/generate_reset_token.go` | 유지 (변경 없음) |
| `internal/gen/gogin/middleware.go` | `claimToVerifyField` 삭제, `internal/auth/` + `internal/middleware/` + `internal/auth/reexport.go` 생성 로직으로 전면 재작성 |
| `internal/gen/gogin/domain.go` | `inferClaimGoType` 삭제 → `ClaimDef.GoType`, 미사용 claims 파라미터 제거, 깨진 fallback 삭제 |
| `internal/gen/gogin/gogin.go` | claims 추출 타입 맞춤 (`map[string]string` → `map[string]ClaimDef`), `transformSource()`에 `pkg/auth` → `internal/auth` import 치환 추가 |
| `internal/crosscheck/types.go` | `Claims` 필드 타입 변경 |
| `internal/crosscheck/claims.go` | `CheckClaims`, `CheckClaimsRego` 시그니처 + 내부 로직 맞춤 |
| `internal/crosscheck/func.go` | JWT built-in 화이트리스트 추가 |
| `internal/orchestrator/validate.go` | claims 추출 타입 맞춤 |
| `internal/contract/hash.go` | `HashClaims` 시그니처 + 해시 입력 변경 |
| `artifacts/manual-for-ai.md` | `auth.type`, claims 타입 힌트 문법, 생성 산출물 구조, import 규칙 문서화 |

## 산출물 구조

`fullend gen` 실행 시 프로젝트 `artifacts/<name>/backend/` 아래 생성:

```
internal/
├── auth/
│   ├── issue_token.go          # //fullend:gen — claims 기반 (preserve 전환 가능)
│   ├── verify_token.go         # //fullend:gen — claims 기반 (preserve 전환 가능)
│   ├── refresh_token.go        # //fullend:gen — claims 기반 (preserve 전환 가능)
│   └── reexport.go             # //fullend:gen — pkg/auth 범용 함수 re-export
├── middleware/
│   └── bearerauth.go           # //fullend:gen (preserve 전환 가능)
└── model/
    └── auth.go                 # CurrentUser struct (claims 기반)
```

## 의존성

- Phase028(genapi 분리) 이후. Phase029와 독립.
- 외부 패키지: `github.com/golang-jwt/jwt/v5` — 산출물의 `go.mod`에 추가 필요.

## 검증

1. `go test ./internal/projectconfig/...` — `type: jwt`, `OrgID: org_id:int64` 파싱 + 검증 규칙
2. `go test ./internal/gen/gogin/...` — auth/middleware/reexport 코드 생성 테스트
3. `go test ./internal/crosscheck/...` — claims 타입 변경 + JWT built-in 화이트리스트
4. `go test ./internal/ssac/...` — IssueToken 타입 확정 struct 생성
5. `fullend validate specs/dummys/zenflow` — claims 검증 통과
6. `fullend gen specs/dummys/zenflow artifacts/zenflow` — `OrgID: org_id:int64` 추가 후 빌드 성공
7. `go vet ./...` 통과
8. `go build ./...` — pkg/auth 부분 삭제 후 fullend 자체 빌드 확인

## 리스크

- `pkg/auth` 부분 삭제로 기존 SSaC `import "github.com/geul-org/fullend/pkg/auth"` 유지 가능 (HashPassword 등은 잔류). JWT 함수만 generated `internal/auth`로 이동하며, `reexport.go`가 통합 import 보장.
- `auth.type` 필수화 — 기존 fullend.yaml에 `type` 없으면 validate ERROR. 마이그레이션: `type: jwt` 한 줄 추가.
- `golang-jwt/jwt/v5`가 산출물 의존성 — 생성된 프로젝트의 `go.mod`에 자동 추가 필요. gluegen의 `go.mod` 생성 로직 확인.
- `reexport.go`가 type alias(`=`)를 사용하므로 `pkg/auth`의 Request/Response struct 변경 시 산출물 재생성 필요 — contract hash로 감지.
