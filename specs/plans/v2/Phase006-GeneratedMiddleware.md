✅ 완료

# Phase 006: 프로젝트별 BearerAuth 미들웨어 생성

## 목표

`pkg/middleware/bearerauth.go` 공용 미들웨어 대신, gluegen이 `fullend.yaml` claims config를 기반으로 프로젝트별 미들웨어를 생성한다. `*model.CurrentUser`를 직접 생성하므로 타입 불일치 panic이 해소된다.

## 배경

Phase 005에서 claims 기반 `model.CurrentUser` 독립 struct을 도입했으나, `pkg/middleware.BearerAuth`는 여전히 `*middleware.CurrentUser`를 gin context에 넣어 타입 불일치 panic 발생:
```
interface conversion: interface {} is *middleware.CurrentUser, not *model.CurrentUser
```

## 변경 사항

### 1. gluegen — 미들웨어 생성 함수 추가

**새 파일**: `internal/gluegen/middlewaregen.go`

claims config를 받아 `internal/middleware/bearerauth.go`를 생성:
- `pkg/auth.VerifyToken` 호출로 JWT 파싱 (기존 로직 재사용)
- claims key → `model.CurrentUser` 필드 매핑을 fullend.yaml 그대로 반영
- int64 필드 (ID 등)는 `float64 → int64` 변환, string 필드는 직접 대입
- 토큰 없거나 무효하면 빈 `model.CurrentUser{}` 세팅 (abort 안 함 — 기존 동작 유지)

생성 결과 예시:
```go
package middleware

import (
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/park-jun-woo/fullend/pkg/auth"
    "github.com/park-jun-woo/dummy-lesson/internal/model"
)

func BearerAuth(secret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        header := c.GetHeader("Authorization")
        if !strings.HasPrefix(header, "Bearer ") {
            c.Set("currentUser", &model.CurrentUser{})
            c.Next()
            return
        }
        token := strings.TrimPrefix(header, "Bearer ")
        out, err := auth.VerifyToken(auth.VerifyTokenRequest{Token: token, Secret: secret})
        if err != nil {
            c.Set("currentUser", &model.CurrentUser{})
            c.Next()
            return
        }
        c.Set("currentUser", &model.CurrentUser{
            Email: out.Email,
            ID:    out.UserID,
            Role:  out.Role,
        })
        c.Next()
    }
}
```

### 2. gluegen/domain.go — server.go import 변경

`generateCentralServer`에서:
- `"github.com/park-jun-woo/fullend/pkg/middleware"` → `"{modulePath}/internal/middleware"` import 변경
- `middleware.BearerAuth("secret")` 호출은 그대로 (함수명 동일)

### 3. gluegen/gluegen.go — Generate 흐름에 미들웨어 생성 추가

`Generate()` 함수에서 claims config가 있으면 `generateMiddleware(intDir, modulePath, claims)` 호출.

### 4. pkg/middleware/bearerauth.go — 유지 (삭제 안 함)

claims 없는 프로젝트의 fallback으로 남겨둔다. claims config가 없으면 기존처럼 `pkg/middleware` import.

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/gluegen/middlewaregen.go` | 신규 — claims 기반 미들웨어 생성 |
| `internal/gluegen/gluegen.go` | 수정 — Generate()에 generateMiddleware 호출 추가 |
| `internal/gluegen/domain.go` | 수정 — server.go import를 claims 유무에 따라 분기 |

## 의존성

- `pkg/auth.VerifyToken` — JWT 파싱 로직 재사용 (변경 없음)
- `pkg/auth.VerifyTokenResponse` — `UserID`, `Email`, `Role` 필드 사용

## 검증 방법

```bash
go build ./cmd/fullend/
./fullend gen specs/dummy-lesson artifacts/dummy-lesson
cd artifacts/dummy-lesson/backend && go mod tidy && go build ./...
# 서버 기동 후 hurl smoke 테스트
JWT_SECRET=test-secret-key go run ./cmd/ -dsn "postgres://postgres:test1224@localhost:15432/dummy_lesson?sslmode=disable" &
hurl --test --variable host=http://localhost:8080 tests/smoke.hurl
```
