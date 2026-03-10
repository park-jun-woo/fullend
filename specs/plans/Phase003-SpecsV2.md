# Phase 003: dummy 프로젝트 SSaC v2 스펙 재작성 + 스모크 테스트

## 목표

dummy-lesson (18개), dummy-study (11개) SSaC 스펙 파일을 v2 문법으로 재작성하고, `fullend validate` → `fullend gen` → `go build` → `hurl --test` 전체 파이프라인 통과를 확인한다.

## 배경

Phase 001~002에서 crosscheck/gluegen이 SSaC v2 타입을 처리할 수 있게 되었으나, 실제 스펙 파일은 아직 v1 문법이다. 스펙 파일을 v2로 전환하고 end-to-end 검증한다.

## 문법 변환 규칙

### v1 → v2 변환 패턴

```go
// v1
// @sequence get
// @model Course.FindByID
// @param CourseID request
// @result course Course

// v2
// @get Course course = Course.FindByID(request.CourseID)
```

```go
// v1
// @sequence guard nil course
// @message "강의를 찾을 수 없습니다"

// v2
// @empty course "강의를 찾을 수 없습니다"
```

```go
// v1
// @sequence authorize
// @action update
// @resource course
// @id CourseID

// v2
// @auth "update" "course" {id: request.CourseID} "권한 없음"
```

```go
// v1
// @sequence call
// @func auth.hashPassword
// @param Password request
// @result hashedPassword string

// v2
// @call string hashedPassword = auth.HashPassword(request.Password)
```

```go
// v1
// @sequence guard state course
// @param course.Status

// v2
// @state course {status: course.Status} "PublishCourse" "상태 전이 불가"
```

```go
// v1
// @sequence response json
// @var course
// @var lessons

// v2
// @response {
//   course: course,
//   lessons: lessons
// }
```

## 변경 파일 목록

### dummy-lesson (18개)

| 파일 | 시퀀스 수 | 복잡도 |
|------|----------|--------|
| `service/auth/login.go` | ~5 | 고 (call 2개) |
| `service/auth/register.go` | ~4 | 중 |
| `service/course/list_courses.go` | ~2 | 저 |
| `service/course/get_course.go` | ~5 | 고 (변수 참조) |
| `service/course/create_course.go` | ~3 | 중 |
| `service/course/update_course.go` | ~4 | 중 (authorize) |
| `service/course/delete_course.go` | ~4 | 중 (authorize) |
| `service/course/publish_course.go` | ~5 | 고 (state) |
| `service/lesson/*` (4개) | ~3 각 | 중 |
| `service/review/*` (3개) | ~3 각 | 중 |
| `service/enrollment/*` (2개) | ~3 각 | 중 |
| `service/payment/list_my_payments.go` | ~2 | 저 |

### dummy-study (SSaC 프로젝트 내)

| 파일 | 비고 |
|------|------|
| `specs/dummy-study/service/*.go` (7개) | SSaC 프로젝트 내 파일 — 수정지시서 필요 여부 확인 |
| `specs/dummy-study/func/billing/*.go` (1개) | Func spec — 변경 불필요 예상 |

> dummy-study는 SSaC 프로젝트(`~/.clari/repos/ssac/specs/dummy-study/`) 내에 위치. SSaC v2 전환 시 SSaC 측에서 이미 변환했을 가능성 있음 → 확인 필요.

## 의존성

- Phase 001 완료 (crosscheck가 v2 타입 처리)
- Phase 002 완료 (gluegen이 v2 타입으로 코드 생성)
- SSaC v2 CLI (`ssac gen`)

## 검증 방법

```bash
# 1. validate
./fullend validate specs/dummy-lesson

# 2. gen
./fullend gen specs/dummy-lesson artifacts/dummy-lesson

# 3. build
cd artifacts/dummy-lesson/backend && go mod tidy && go build ./...

# 4. hurl test (서버 기동 후)
hurl --test --variable host=http://localhost:8080 artifacts/dummy-lesson/tests/*.hurl
```
