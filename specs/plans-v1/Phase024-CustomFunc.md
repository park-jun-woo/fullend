✅ 완료

# Phase 024: Custom Func — pkg Go 모듈 + gen/validate 통합

## 목표

SSaC `@sequence call`이 참조하는 외부 함수를 fullend가 관리한다.
기본 구현은 `pkg` Go 모듈로 퍼블리시하여 `go get`으로 가져온다.
프로젝트 커스텀이 있으면 로컬 패키지가 우선한다.

```
현재:
  @sequence call → Handler에 func 필드 → main.go에서 수동 주입 (nil panic 위험)

목표:
  @sequence call → @func auth.hashPassword
  커스텀 없으면 → go get github.com/geul-org/fullend → import하여 사용
  커스텀 있으면 → specs/<project>/func/auth/ → internal/auth/로 복사
```

## 전제 조건

- ssac 수정지시서 011 완료 필수:
  - `@sequence password` 폐지 → `@sequence call` 통합
  - `@func package.funcName` 파싱 (Package 필드 추가)
  - `func(In) (Out, error)` 시그니처 코드젠
  - 패키지 레벨 함수 호출 코드 생성

---

## pkg Go 모듈

### 모듈 구조

```
pkg/                                    # fullend 루트 go.mod의 서브패키지 (별도 go.mod 없음)
├── auth/
│   ├── hash_password.go            # bcrypt 해싱
│   ├── verify_password.go          # bcrypt 검증
│   └── issue_token.go              # JWT 발급
├── payment/
│   └── ...                         # 향후 확장
└── notification/
    └── ...                         # 향후 확장
```

### 사용 방식

프로젝트 `go.mod`에서:
```
go get github.com/geul-org/fullend
```

생성된 코드에서:
```go
import auth "github.com/geul-org/fullend/pkg/auth"

// 패키지 레벨 함수 호출
out, err := auth.HashPassword(auth.HashPasswordInput{Password: password})
```

### import 경로 결정

SSaC 스펙 파일의 import 선언이 진실:

```go
// 기본 (pkg 사용)
import auth "github.com/geul-org/fullend/pkg/auth"

// 커스텀 (프로젝트 로컬)
import auth "<module>/internal/auth"
```

SSaC 파서가 Go AST로 import를 읽고, 생성 코드에 그대로 옮긴다.
fullend는 커스텀일 때만 `specs/<project>/func/auth/` → `internal/auth/`로 복사.
기본이면 `go get`으로 의존성 해결.

---

## 변경 파일 목록

### 1. pkg 모듈 (신규)

| 파일 | 내용 |
|---|---|
| `pkg/auth/hash_password.go` | bcrypt 해싱 |
| `pkg/auth/verify_password.go` | bcrypt 검증 |
| `pkg/auth/issue_token.go` | JWT 발급 |

### 2. fullend 코드 수정

| 파일 | 변경 |
|---|---|
| `artifacts/internal/orchestrator/detect.go` | `KindFunc` 추가 (10번째 SSOT kind) |
| `artifacts/internal/orchestrator/validate.go` | `validateFunc()` — func 파일 파싱 + 교차 검증 |
| `artifacts/internal/orchestrator/gen.go` | `genFunc()` — 커스텀 복사 또는 import 경로 결정 |
| `artifacts/internal/orchestrator/status.go` | `statusFunc()` — func 현황 출력 |
| `artifacts/internal/funcspec/parser.go` (신규) | Go AST로 @func, @description, Input/Output struct 파싱 |
| `artifacts/internal/funcspec/parser_test.go` (신규) | 파서 테스트 |
| `artifacts/internal/crosscheck/func.go` (신규) | Func ↔ SSaC 교차 검증 |
| `artifacts/internal/crosscheck/crosscheck.go` | `CheckFuncs()` 호출 추가 |
| `artifacts/internal/gluegen/gluegen.go` | main.go 생성 시 Handler 초기화 (func 필드 제거 반영) |

### 3. dummy-lesson 스펙 수정

| 파일 | 변경 |
|---|---|
| `specs/dummy-lesson/service/auth/register.go` | `@func hashPassword` → `@func auth.hashPassword` + import 추가 |
| `specs/dummy-lesson/service/auth/login.go` | `@sequence password` → `@sequence call` + `@func auth.verifyPassword` + import 추가 |
| `specs/dummy-lesson/service/auth/login.go` | `@func issueToken` → `@func auth.issueToken` |
| `specs/dummy-study/service/login.go` | ssac 011에서 이미 변경됨 — 확인만 |

---

## 상세 설계

### 1. pkg/auth 기본 구현

#### `hash_password.go`
```go
package auth

import "golang.org/x/crypto/bcrypt"

// @func hashPassword
// @description 평문 비밀번호를 bcrypt 해시로 변환한다

type HashPasswordInput struct {
    Password string
}

type HashPasswordOutput struct {
    HashedPassword string
}

func HashPassword(in HashPasswordInput) (HashPasswordOutput, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
    return HashPasswordOutput{HashedPassword: string(hash)}, err
}
```

#### `verify_password.go`
```go
package auth

import "golang.org/x/crypto/bcrypt"

// @func verifyPassword
// @description 저장된 해시와 평문 비밀번호가 일치하는지 검증한다

type VerifyPasswordInput struct {
    PasswordHash string
    Password     string
}

type VerifyPasswordOutput struct{}

func VerifyPassword(in VerifyPasswordInput) (VerifyPasswordOutput, error) {
    err := bcrypt.CompareHashAndPassword([]byte(in.PasswordHash), []byte(in.Password))
    return VerifyPasswordOutput{}, err
}
```

#### `issue_token.go`
```go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

// @func issueToken
// @description 인증된 사용자 정보로 JWT 액세스 토큰을 발급한다

type IssueTokenInput struct {
    UserID int64
    Email  string
    Role   string
}

type IssueTokenOutput struct {
    AccessToken string
}

func IssueToken(in IssueTokenInput) (IssueTokenOutput, error) {
    claims := jwt.MapClaims{
        "user_id": in.UserID,
        "email":   in.Email,
        "role":    in.Role,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString([]byte("secret"))
    return IssueTokenOutput{AccessToken: signed}, err
}
```

### 2. `funcspec/parser.go` — Func 스펙 파서

Go AST로 파싱하여 추출:
- `@func` — 함수 식별자
- `@description` — 자연어 설명
- Input struct 필드 목록 (이름, 타입)
- Output struct 필드 목록 (이름, 타입)

```go
type FuncSpec struct {
    Package      string  // "auth"
    Name         string  // "hashPassword"
    Description  string  // "@description 값"
    InputFields  []Field // HashPasswordInput의 필드들
    OutputFields []Field // HashPasswordOutput의 필드들
    HasBody      bool    // TODO가 아닌 실제 구현이 있는지
}

type Field struct {
    Name string
    Type string
}
```

파싱 대상: `specs/<project>/func/` 또는 `pkg/` 모두 동일 형식.

### 3. `crosscheck/func.go` — 교차 검증

| 검증 | 조건 | 레벨 |
|---|---|---|
| func 구현 없음 | SSaC `@func pkg.X` → import가 로컬 경로인데 `specs/<project>/func/<pkg>/` 에 파일 없음 | ERROR + 구현 지침 |
| func 스펙 불일치 | SSaC `@func pkg.X` → func 파일은 있지만 `@func` 이름 불일치 | ERROR |
| Input 필드 일치 | SSaC `@param` ↔ FuncSpec InputFields | ERROR |
| Output 필드 일치 | SSaC `@result` ↔ FuncSpec OutputFields | ERROR |
| 본체 미구현 | `HasBody == false` (TODO) | WARNING |

로컬 vs 외부 판별:
- import 경로가 `github.com/geul-org/fullend/pkg/` → fullend 기본 제공 (검증 스킵, go get으로 해결)
- 그 외 → 로컬 커스텀, `specs/<project>/func/`에 구현 필수

### ERROR 시 구현 지침 프롬프트

구현이 없는 func에 대해 SSaC `@param`/`@result`에서 스켈레톤을 자동 유도하여 출력:

```
ERROR: @func billing.calculateRefund — 구현 없음

다음 파일을 작성하세요: specs/<project>/func/billing/calculate_refund.go

package billing

// @func calculateRefund
// @description <이 함수가 무엇을 하는지 한 줄로 설명>

type CalculateRefundInput struct {
    Reservation Reservation
}

type CalculateRefundOutput struct {
    Refund Refund
}

func CalculateRefund(in CalculateRefundInput) (CalculateRefundOutput, error) {
    // TODO: implement
    return CalculateRefundOutput{}, nil
}
```

LLM에 이 에러 메시지를 그대로 전달하면 바로 파일 생성 가능.

### 4. gen 흐름

```
genFunc():
  1. SSaC 스펙의 import에서 func 패키지 경로 확인
  2. 로컬 패키지인 경우 (import "<module>/internal/<pkg>"):
     → specs/<project>/func/<pkg>/*.go → artifacts/<project>/backend/internal/<pkg>/로 복사
  3. 외부 모듈인 경우 (import "github.com/geul-org/fullend/<pkg>"):
     → 복사 불필요. go get으로 의존성 해결
```

import 경로 결정은 SSaC 스펙 파일 작성자의 책임. fullend는 그에 따라 복사 또는 의존성 추가만 수행.

### 5. detect.go — KindFunc

```go
KindFunc SSOTKind = "Func"

// detect: func/<pkg>/*.go 패턴 (프로젝트 커스텀 존재 여부)
// func/ 디렉토리 없으면 Func kind 미감지 → 검증 스킵 (ERROR 아님)
// SSaC @func 참조가 있는데 func/ 없으면 crosscheck에서 잡음
```

### 6. glue-gen — main.go

패키지 함수이므로 Handler에 func 필드 주입 불필요. glue-gen은 기존 `func(args ...interface{})` 필드를 더 이상 생성하지 않음:
```go
Auth: &authsvc.Handler{
    UserModel: model.NewUserModel(conn),
    // HashPassword, IssueToken 필드 제거됨 — 패키지 함수로 직접 호출
},
```

---

## 의존성

| 패키지 | 용도 |
|---|---|
| `go/ast`, `go/parser` | Func 스펙 파일 파싱 |
| `github.com/geul-org/fullend` | 기본 func 구현 (Go 모듈) |

pkg 모듈 내부 의존성:
| 패키지 | 용도 |
|---|---|
| `golang.org/x/crypto/bcrypt` | hash/verify password |
| `github.com/golang-jwt/jwt/v5` | issue token |

---

## 검증 방법

```bash
# 1. fullend 테스트
cd ~/.clari/repos/fullend && go test ./...

# 2. validate 통과
go run ./artifacts/cmd/fullend/main.go validate specs/dummy-lesson

# 3. gen 성공
go run ./artifacts/cmd/fullend/main.go gen specs/dummy-lesson artifacts/dummy-lesson

# 4. 생성된 서버 빌드 (go get으로 pkg 의존성 해결)
cd artifacts/dummy-lesson/backend && go build ./cmd/main.go

# 5. 서버 기동 + Hurl 테스트
./main -dsn "..." &
hurl --test --variable host=http://localhost:8080 ../tests/*.hurl
```

## 참고

- pkg는 Go 모듈이므로 Go 생태계에 묶임. 향후 멀티 언어 지원 시 언어별 레지스트리로 분리 가능 (계약 In/Out 자체는 언어 중립)
- pkg 퍼블리시 전까지는 `go.mod replace`로 로컬 참조
