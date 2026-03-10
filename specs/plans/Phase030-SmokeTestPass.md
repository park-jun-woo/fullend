# Phase 030: 스모크 테스트 전체 통과

## 목표

dummy-lesson `fullend gen` → `go build` → `hurl --test smoke.hurl` 전체 통과.

## 현황

18개 엔드포인트 중 9개 통과 (Register ~ ListCourses). GetCourse에서 `$.instructor` 없어서 실패.

## 완료된 수정 (스모크 테스트 중 발견)

### A. gluegen — transformSource receiver 중복 방지

| 파일 | 변경 |
|------|------|
| `internal/gluegen/gluegen.go` | `transformSource`에서 `\nfunc ` 뒤에 `(`가 이미 있으면(receiver 존재) 치환 건너뛰기. 기존 artifacts에 재실행 시 `(s *Server) (s *Server)` 중복 receiver 버그 수정 |

### B. authzgen — via 매핑 create 시 부모 테이블 직접 조회

| 파일 | 변경 |
|------|------|
| `internal/gluegen/authzgen.go` | `lookupOwner`에 `action` 파라미터 추가. `via` ownership 매핑에서 action이 `create`이면 자식 테이블 JOIN 대신 부모 테이블에서 직접 소유자 조회. 아직 자식 레코드가 없는 create 시 403 발생 버그 수정 |

### C. dummy-lesson model — @dto 추가

| 파일 | 변경 |
|------|------|
| `specs/dummy-lesson/model/session.go` | `IssueTokenResponse`, `HashPasswordResponse` DTO 추가. 교차 검증 경고 해소 |

### D. hurl-gen — filter 테스트 값 불일치 수정

| 파일 | 변경 |
|------|------|
| `internal/gluegen/hurl.go` | x-filter 쿼리 파라미터 값을 `test` → `test_string`으로 수정. `generateDummyValue`가 문자열 필드에 `test_string`을 생성하므로 filter 값도 일치시킴 |

### E. hurl-gen — state transition 순서 수정

| 파일 | 변경 |
|------|------|
| `internal/gluegen/hurl.go` | PUT 중 request body가 없는 것(PublishCourse 등 상태 전환)을 `transitionSteps`로 분리하여 create 직후, read 이전에 배치. 기존에는 update 그룹에 속해 read 뒤에 실행되어 `published = TRUE` 조건의 List가 빈 결과 반환 |

### F. x-include를 런타임 쿼리 파라미터에서 코드젠 메타데이터로 변경

`x-include`는 관련 테이블 JOIN을 코드젠 시 항상 적용하기 위한 메타데이터이므로, 런타임 `?include=` 쿼리 파라미터로 노출하면 안 된다.

| 파일 | 변경 |
|------|------|
| `internal/gluegen/queryopts.go` | `IncludeConfig` struct 삭제, `QueryOptsConfig`에서 `Include` 필드 삭제, `ParseQueryOpts`에서 `?include=` 파싱 로직 삭제 |
| `internal/gluegen/model_impl.go` | include 로직에서 `containsStr(opts.Includes, ...)` 조건 제거 → include 항상 실행 |
| `internal/gluegen/hurl.go` | `?include=` 쿼리 파라미터 생성 제거 |
| `internal/gluegen/gluegen.go` | `buildQueryOptsConfig`에서 `x-include` → `Include:` 설정 생성 제거 |
| SSaC 수정지시서015 | `QueryOpts`에서 `Includes []string` 필드 제거 ✅ 실행 완료

### G. OpenAPI ↔ SSaC 역방향 검증 추가

SSaC 수정지시서016 ✅ 실행 완료. `validateResponse`에 역방향 검증(OpenAPI response → SSaC @var), `validateRequest`에 역방향 검증(OpenAPI request → SSaC @param) 추가.

### H. GetCourse에서 instructor 미반환 — 방향 결정 완료

방안 A 채택. SSaC GetCourse 스펙에 instructor 조회 시퀀스(`@model User.FindByID` + `@param course.InstructorID`) 추가, `@var instructor` 포함. x-include나 FindByID include 패턴 없이 SSaC 시퀀스만으로 해결. I항(FindByID + include) 불필요.

SSaC 수정지시서018 ✅ 실행 완료 — `fullend validate` 시 누락된 @var에 대해 조치 방법을 안내.

**차후 확인 예정**: `specs/dummy-lesson/service/course/get_course.go`에 instructor 조회 시퀀스 실제 추가 + `fullend validate` 에러 메시지 소비 확인.

## 미해결 문제

### K. DDL → OpenAPI/SSaC 역방향 검증 누락

Phase031로 분리.


## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| SSaC 수정지시서016 | `validateResponse`·`validateRequest` 역방향 검증 추가 ✅ 실행 완료 |
| `internal/crosscheck/ddl_coverage.go` | DDL → OpenAPI/SSaC 역방향 검증 (WARNING) + `@archived` 파싱 |
| `specs/dummy-lesson/service/course/get_course.go` | instructor 조회 시퀀스 추가 (방안 A) |
| 기타 | 스모크 테스트 중 발견되는 추가 문제 |

## 의존성

- SSaC 수정지시서015 (QueryOpts Includes 제거) ✅ 실행 완료
- SSaC 수정지시서016 (validateRequest·validateResponse 역방향 검증) ✅ 실행 완료
- SSaC 수정지시서017 (@archived DDL 파싱) → fullend 이관 (fullend에서 직접 구현)
- SSaC 수정지시서018 (역방향 에러 수정 가이드) ✅ 실행 완료

## 검증 방법

```bash
# 1. fullend validate — 새 crosscheck 규칙 동작 확인
fullend validate specs/dummy-lesson

# 2. gen + build
fullend gen specs/dummy-lesson artifacts/dummy-lesson
cd artifacts/dummy-lesson/backend && go mod tidy && go build ./...

# 3. 스모크 테스트 전체 통과
hurl --test --variable host=http://localhost:8080 smoke.hurl

# 4. go test
go test ./internal/...
```
