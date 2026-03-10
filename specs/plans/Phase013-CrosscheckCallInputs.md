# Phase 013: Func ↔ SSaC crosscheck — @call Inputs 대응

## 배경

ssac 수정지시서007에서 `@call` 인자가 positional `(arg1, arg2)` → named field `({Key: val})` 로 변경됨.
ssac 파서가 `seq.Args` 대신 `seq.Inputs` (map[string]string)에 저장하게 되면서,
fullend crosscheck의 `Func ↔ SSaC` 검증이 `seq.Args`를 참조하여 항상 0개로 카운트 → 오탐 ERROR 발생.

## 목표

`internal/crosscheck/func.go`의 @call 검증 로직을 `seq.Inputs` 기반으로 전환.

## 변경 파일

### 1. `internal/crosscheck/func.go` — 검증 로직 전환 ✅ 완료

**Rule 1: 필드 수 비교**
- 변경 전: `countNonLiteralArgs(seq.Args)` vs `len(spec.RequestFields)`
- 변경 후: `len(seq.Inputs)` vs `len(spec.RequestFields)`

**Rule 2: 필드명 일치 검증 (신규)**
- `seq.Inputs`의 key가 FuncSpec Request struct 필드명에 존재하는지 확인
- 기존 positional 방식에서는 순서 기반 타입 비교였으나, named field에서는 이름 기반 존재 검증으로 변경

**Rule 3: Result ↔ Response** — 변경 없음

**Rule 4: Source 변수 정의 여부**
- 변경 전: `seq.Args`의 `arg.Source` 순회
- 변경 후: `seq.Inputs`의 value에서 `strings.SplitN(value, ".", 2)`로 source 추출

**삭제된 함수:**
- `countNonLiteralArgs()` — `seq.Args` 카운트, 불필요
- `checkPositionalTypes()` — positional 순서 기반 타입 비교, named field에서는 불필요
- `resolveArgType()` — `checkPositionalTypes` 전용, 불필요

**미사용 잔존 함수 (향후 정리 가능):**
- `resolveDDLColumnType()`, `resolveOpenAPIFieldType()`, `openAPITypeToGo()`, `typesCompatible()` — 현재 호출처 없음. 향후 named field 타입 검증 확장 시 재활용 가능하므로 당장 삭제하지 않음.

**`generateSkeleton()` 수정:**
- `seq.Args` 순회 → `seq.Inputs` key 순회로 변경

### 2. `specs/dummy-gigbridge/service/` — @call 문법 마이그레이션 ✅ 완료

| 파일 | 변경 |
|---|---|
| `auth/register.go` | `auth.HashPassword(request.Password)` → `auth.HashPassword({Password: request.Password})` |
| `auth/login.go` | `auth.VerifyPassword(...)` → `({PasswordHash: user.PasswordHash, Password: request.Password})`, `auth.IssueToken(...)` → `({UserID: user.ID, Email: user.Email, Role: user.Role})` |
| `proposal/accept_proposal.go` | `billing.HoldEscrow(...)` → `({GigID: gig.ID, Amount: gig.Budget, ClientID: gig.ClientID})` |
| `gig/approve_work.go` | `billing.ReleaseFunds(...)` → `({GigID: gig.ID, Amount: gig.Budget, FreelancerID: gig.FreelancerID})` |

## 의존성

- ssac 수정지시서007 완료 필수 (seq.Inputs 사용)

## 검증 방법

```bash
go build ./internal/crosscheck/
go test ./...
go run ./cmd/fullend validate specs/dummy-gigbridge
```

- `Func ↔ SSaC` @call 관련 ERROR가 사라져야 함
- @call Inputs 필드명과 Request 필드명 불일치 시 ERROR 출력 확인

## 상태: ✅ 완료
