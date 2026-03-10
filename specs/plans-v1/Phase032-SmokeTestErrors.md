# Phase 032: 스모크 테스트 오류 기록

## dummy-lesson 빌드 오류 (2026-03-10)

### 오류 1: SSaC codegen — @param source 변수 참조 미해결

**파일:** `artifacts/dummy-lesson/backend/internal/service/course/get_course.go:43`

```
undefined: instructorID
not enough arguments in call to h.UserModel.FindByID
```

**원인:** SSaC codegen이 `@param InstructorID course`에서 source `course`를 이전 @result 변수로 인식하지 못함.
- `course.InstructorID`로 생성해야 하는데 `instructorID`(미정의 변수)로 생성
- `FindByID(int64, model.QueryOpts)` 시그니처에 `opts` 인자 누락

**대응:** SSaC 수정지시서019 발송. 회신 대기 중.

**SSaC spec 변경 (get_course.go):**
```go
// @sequence get
// @model User.FindByID
// @param InstructorID course
// @result instructor User
```
이 시퀀스를 추가하여 instructor 조회를 포함시킴 (Phase030 H 해소).
