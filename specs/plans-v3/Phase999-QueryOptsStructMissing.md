# Phase016: QueryOpts struct 누락 — codegen 빌드 실패 수정

## 목표

`fullend gen`으로 생성된 backend가 `go build` 시 `undefined: QueryOpts` 에러로 실패하는 버그를 수정한다.

## 버그 요약

| 항목 | 내용 |
|------|------|
| 증상 | `artifacts/*/backend/internal/model/queryopts.go`에서 `undefined: QueryOpts` (4건) |
| 원인 | `generateQueryOpts()` 템플릿이 `QueryOpts` 타입을 사용하지만 정의하지 않음 |
| 발생 시점 | **처음부터** — 3개 커밋 모두 `type QueryOpts struct` 없음 |
| 발견 경로 | dummy-zenflow `fullend gen` → `go build` 실패 |

## 근본 원인 분석

`QueryOpts` struct의 책임 소재가 **SSaC와 fullend 사이에서 분리**되어 있다.

| 역할 | 담당 | 파일 |
|------|------|------|
| `type QueryOpts struct` **정의** | SSaC | `ssac/generator/go_interface.go:345-354` |
| `QueryOpts` **사용** (ParseQueryOpts, BuildSelectQuery, BuildCountQuery) | fullend | `internal/gluegen/queryopts.go` 템플릿 |

SSaC의 `go_interface.go`는 `needQueryOpts == true`일 때(pagination을 쓰는 모델이 있을 때) 인터페이스 파일에 `type QueryOpts struct`를 출력한다. 따라서 **pagination을 쓰는 프로젝트에서는 빌드가 성공**한다 — SSaC가 struct를 생성하므로.

그러나 **pagination을 안 쓰는 프로젝트**(zenflow처럼)에서는 SSaC가 `QueryOpts`를 생성하지 않는데, fullend는 `queryopts.go`를 **무조건 생성**하므로 빌드가 실패한다.

## 변경 파일

| 파일 | 변경 |
|------|------|
| `internal/gluegen/queryopts.go` | 템플릿에 `type QueryOpts struct` 정의 추가 |
| `internal/gluegen/queryopts_test.go` | 기존 테스트에 struct 존재 확인 추가 |

## 수정 내용

`internal/gluegen/queryopts.go`의 `generateQueryOpts()` 템플릿 문자열에서 `FilterConfig` struct 뒤, `ParseQueryOpts` 함수 앞에 추가:

```go
// QueryOpts holds parsed query parameters for pagination, sort, and filter.
type QueryOpts struct {
	Limit   int
	Offset  int
	Cursor  string
	SortCol string
	SortDir string
	Filters map[string]string
}
```

SSaC 측에서도 동일 struct를 생성하지만, 같은 package 내 중복 정의는 컴파일 에러이므로 **SSaC 측에서의 조건부 생성과 충돌하지 않는지 확인** 필요:

- SSaC `go_interface.go:345`는 `needQueryOpts == true`일 때만 생성
- fullend `queryopts.go`는 **무조건** 생성
- 둘 다 생성하면 **중복 정의 에러**

→ **해법**: fullend 템플릿에 struct를 추가하되, SSaC 측의 `needQueryOpts` 분기에서 fullend가 이미 생성한다는 전제로 **SSaC 측 생성을 제거**하거나, 반대로 fullend에서 조건부 생성.

**가장 단순한 해법**: fullend가 `queryopts.go`에 `QueryOpts` struct를 **항상** 포함시키고, SSaC 수정지시서를 보내 SSaC 측 `needQueryOpts` 분기를 제거한다.

## 실행 순서

1. `internal/gluegen/queryopts.go` 템플릿에 `type QueryOpts struct` 추가
2. `go test ./internal/gluegen/...` 통과 확인
3. `go test ./...` 전체 통과 확인
4. dummy-zenflow으로 검증: `fullend gen` → `go build` 성공 확인
5. dummy-gigbridge으로 검증: `fullend gen` → `go build` 성공 확인 (SSaC 중복 정의 여부)
6. gigbridge에서 중복 에러 발생 시 → SSaC 수정지시서 발송 (needQueryOpts 분기 제거)

## 검증

```bash
go test ./internal/gluegen/...
go test ./...
./fullend gen specs/zenflow artifacts/zenflow
cd artifacts/zenflow/backend && go build -o server ./cmd/
./fullend gen specs/gigbridge artifacts/gigbridge
cd artifacts/gigbridge/backend && go build -o server ./cmd/
```
