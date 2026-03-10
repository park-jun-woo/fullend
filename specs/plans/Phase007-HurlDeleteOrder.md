✅ 완료

# Phase 007: hurl-gen 삭제 순서 FK 의존성 정렬

## 목표

smoke.hurl의 DELETE 순서를 DDL FK 의존성 기반으로 정렬한다. 자식 테이블 → 부모 테이블 순서로 삭제해서 FK constraint 위반을 방지한다.

## 배경

현재 `buildScenarioOrder`에서 DELETE는 path depth DESC → path 알파벳순으로 정렬. depth가 같은 리소스(`/courses/{id}`, `/lessons/{id}`, `/reviews/{id}`)는 알파벳순이 되어 DeleteCourse가 먼저 실행되고 FK 위반 발생.

```
현재 순서 (FK 위반):
  DeleteCourse ← lessons, enrollments, reviews가 참조 중 → 500
  DeleteLesson
  DeleteReview

올바른 순서:
  DeleteReview   (reviews → courses FK)
  DeleteLesson   (lessons → courses FK)
  DeleteCourse   (부모 테이블, 자식 삭제 후)
```

## 해결 방안

### hurl.go `buildScenarioOrder` 수정

DELETE 정렬 시 DDL FK 의존 그래프를 사용:

1. `generateHurlTests`에 `specsDir` 파라미터 추가
2. DDL에서 FK 관계를 파싱 → 테이블 간 의존 그래프 구축
3. operationId에서 리소스 테이블을 추론 (path의 첫 번째 non-param segment → 복수형 테이블명)
4. DELETE 정렬: FK 의존 그래프에서 자식(FK를 가진 테이블) → 부모 순서 (위상 정렬)

### FK 의존 그래프 예시 (dummy-lesson)

```
courses ← lessons (lessons.course_id → courses.id)
courses ← enrollments (enrollments.course_id → courses.id)
courses ← reviews (reviews.course_id → courses.id)
users ← enrollments (enrollments.user_id → users.id)
enrollments ← payments (payments.enrollment_id → enrollments.id)
```

위상 정렬 결과 (삭제 순서):
```
payments → reviews → lessons → enrollments → courses → users
```

DELETE 엔드포인트가 있는 것만 필터링:
```
DeleteReview → DeleteLesson → DeleteCourse
```

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/gluegen/hurl.go` | 수정 — DELETE 정렬에 FK 의존 그래프 사용 |
| `internal/gluegen/gluegen.go` | 수정 — `generateHurlTests`에 specsDir 전달 |

## 의존성

- `parseDDLFiles` (model_impl.go) — 이미 FK 관계를 파싱하는 `ddlColumn.FKTable` 존재. 재사용.

## 검증 방법

```bash
go build ./cmd/fullend/
./fullend gen specs/dummy-lesson artifacts/dummy-lesson
# smoke.hurl에서 DELETE 순서 확인: DeleteReview → DeleteLesson → DeleteCourse
cd artifacts/dummy-lesson/backend && go build ./...
# 서버 기동 후 hurl --test 18/18 통과
```
