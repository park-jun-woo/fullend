# Phase 004: dummy-lesson 스모크 테스트 (validate → gen → build)

## 목표

SSaC v2 수정(수정지시서001) 반영 후, dummy-lesson 전체 파이프라인을 통과시킨다.

## 배경

- Phase 001~002: crosscheck/gluegen SSaC v2 마이그레이션 완료
- Phase 003: dummy-lesson 18개 SSaC 스펙 v2 재작성 + CurrentUser 모델 추가 완료
- SSaC 수정지시서001: `@get` 0-arg 허용, `@delete` 0-arg WARNING, `!` 접미사 WARNING 억제

Phase 003에서 validate 실패한 `Course.List()` Args 누락 에러가 SSaC 수정으로 해결되었으므로, 최신 SSaC를 반영하고 전체 파이프라인을 검증한다.

## 단계

### A. SSaC 최신 반영

```bash
go mod tidy
go build ./cmd/fullend/
```

### B. validate 통과 확인

```bash
./fullend validate specs/dummy-lesson
```

Phase 003에서 남아있던 Cross warnings 4개 확인:
- `enroll_course.go` — payments 테이블에 payment_method 컬럼 누락 → DDL 수정
- `lessons.created_at` — OpenAPI Lesson 스키마에 없음 → DDL `-- @archived` 또는 OpenAPI 추가
- `users.password_hash` — OpenAPI User 스키마에 없음 → DDL `-- @archived` (보안 필드)
- `users.created_at` — OpenAPI User 스키마에 없음 → DDL `-- @archived` 또는 OpenAPI 추가

### C. gen 실행

```bash
./fullend gen specs/dummy-lesson artifacts/dummy-lesson
```

### D. build 확인

```bash
cd artifacts/dummy-lesson/backend && go mod tidy && go build ./...
```

### E. 서버 기동 + hurl 테스트 (선택)

```bash
# DB 기동 확인
psql -h localhost -p 15432 -U postgres -d dummy_lesson -c "SELECT 1"

# 서버 기동
cd artifacts/dummy-lesson/backend && go run . &

# hurl 테스트
hurl --test --variable host=http://localhost:8080 artifacts/dummy-lesson/tests/*.hurl
```

## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| `go.mod` / `go.sum` | SSaC 최신 참조 (`go mod tidy`) |
| `specs/dummy-lesson/db/payments.sql` | `payment_method` 컬럼 추가 (WARN 해결) |
| `specs/dummy-lesson/db/courses.sql` 등 | `-- @archived` 주석 추가 (WARN 해결, 필요 시) |

## 의존성

- Phase 003 완료 (SSaC v2 스펙 재작성)
- SSaC 수정지시서001 완료 (`@get` 0-arg 허용)

## 검증 방법

```bash
./fullend validate specs/dummy-lesson   # ERROR 0개
./fullend gen specs/dummy-lesson artifacts/dummy-lesson
cd artifacts/dummy-lesson/backend && go mod tidy && go build ./...
```
