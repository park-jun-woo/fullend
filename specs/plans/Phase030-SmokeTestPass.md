# Phase 030: 스모크 테스트 전체 통과

## 목표

dummy-lesson `fullend gen` → `go build` → `hurl --test smoke.hurl` 전체 통과.

## 현황

Phase 029에서 18개 엔드포인트 중 9개 통과 (Register ~ ListCourses). GetCourse에서 `$.instructor` 없어서 실패.

## 발견된 문제

### 1. OpenAPI ↔ SSaC 교차 검증 누락

OpenAPI 응답 스키마에 `instructor` 필드가 있지만 SSaC GetCourse 스펙에는 instructor 조회 시퀀스가 없다. `fullend validate`가 이 불일치를 감지하지 못한다.

**수정**: crosscheck에 "OpenAPI 응답 필드 vs SSaC @var 교차 검증" 규칙 추가.

### 2. GetCourse에서 instructor 미반환

SSaC GetCourse 스펙에 instructor 조회가 없어 응답에 instructor가 포함되지 않는다.

**수정 방안** (택 1):
- A. SSaC GetCourse 스펙에 instructor 조회 시퀀스 추가 — `x-include`의 `instructor_id:users.id` 매핑을 활용하여 서비스 레이어에서 instructor를 조회하고 응답에 포함
- B. OpenAPI GetCourse 응답 스키마에서 instructor 제거 — instructor는 List에서만 제공

### 3. FindByID + include 지원

현재 model_impl의 include 로직은 List 메서드에만 적용된다. FindByID 같은 단일 조회에서도 관련 리소스를 포함해야 하면 별도 처리가 필요하다.

**수정**: GetCourse 서비스 코드에서 course 조회 후 별도로 instructor를 조회하거나, model_impl의 FindByID에도 include 헬퍼를 호출하는 패턴 추가.

### 4. 나머지 스모크 테스트 실패 가능성

9번째 이후 엔드포인트(ListMyEnrollments, ListMyPayments, ListLessons, ListReviews, UpdateCourse, UpdateLesson, DeleteCourse, DeleteLesson, DeleteReview)는 아직 미검증. 추가 문제 발견 시 이 계획에 추가.

## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| `internal/crosscheck/crosscheck.go` | OpenAPI 응답 필드 vs SSaC @var 교차 검증 규칙 추가 |
| `specs/dummy-lesson/service/course/get_course.go` | instructor 조회 시퀀스 추가 (방안 A 선택 시) |
| `internal/gluegen/model_impl.go` | FindByID에 include 지원 추가 (필요 시) |
| 기타 | 스모크 테스트 중 발견되는 추가 문제 |

## 의존성

- SSaC 수정지시서015 (QueryOpts Includes 제거) — 병렬 진행 가능

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
