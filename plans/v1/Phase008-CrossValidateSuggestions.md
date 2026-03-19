# ✅ Phase 008: 교차 검증 수정 제안 (Cross-Validate Suggestions)

## 목표

교차 검증에서 불일치 발견 시 **구체적 수정 제안**을 함께 출력한다.
LLM 없이 `fmt.Sprintf` 기반 템플릿으로 생성한다.

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `artifacts/internal/crosscheck/types.go` | CrossError에 `Suggestion string` 필드 추가 |
| `artifacts/internal/crosscheck/openapi_ddl.go` | checkXSort, checkXFilter, checkXInclude에 Suggestion 생성 |
| `artifacts/internal/crosscheck/ssac_ddl.go` | checkResultType, checkParamTypes에 Suggestion 생성 |
| `artifacts/internal/crosscheck/ssac_openapi.go` | 양방향 매칭에 Suggestion 생성 |
| `artifacts/internal/reporter/reporter.go` | Suggestion이 있으면 `  → 제안:` 줄 출력 |

## 수정 제안 규칙

| 규칙 | Suggestion |
|---|---|
| x-sort 컬럼 없음 | `DDL에 추가: ALTER TABLE {table} ADD COLUMN {col} -- TODO: 타입 지정;` |
| x-sort 인덱스 없음 | `DDL에 추가: CREATE INDEX idx_{table}_{col} ON {table}({col});` |
| x-filter 컬럼 없음 | `DDL에 추가: ALTER TABLE {table} ADD COLUMN {col} -- TODO: 타입 지정;` |
| x-include 테이블 없음 | `DDL에 추가: CREATE TABLE {table}s (...);` |
| x-include FK 없음 | `DDL에 추가: ALTER TABLE {src} ADD COLUMN {dst}_id BIGINT REFERENCES {dst}(id);` |
| @result 타입 테이블 없음 | `DDL에 추가: CREATE TABLE {table} (...); 또는 model에 // @dto 선언` |
| @param 컬럼 없음 | `DDL에 추가: ALTER TABLE {table} ADD COLUMN {col} -- TODO: 타입 지정;` |
| SSaC → OpenAPI 없음 | `OpenAPI에 추가: operationId: {name}` |
| OpenAPI → SSaC 없음 | `SSaC에 추가: func {opID}(w http.ResponseWriter, r *http.Request) {}` |

### 제약

- 타입을 알 수 없는 경우 `-- TODO: 타입 지정` placeholder 사용
- 테이블명은 기존 `modelToTable()`, `resolveTableName()` 로직 재활용
- x-sort 인덱스 제안은 테이블명 특정 가능 (컬럼이 속한 테이블)

## 출력 예시

```
✗ Cross        2 errors, 1 warning
    [ERROR] x-sort ↔ DDL: GET /courses (ListCourses) — x-sort column "price" (→ price) not found in any DDL table
      → 제안: DDL에 추가: ALTER TABLE courses ADD COLUMN price -- TODO: 타입 지정;
    [WARNING] x-sort ↔ DDL index: GET /courses (ListCourses) — x-sort column "created_at" has no index
      → 제안: DDL에 추가: CREATE INDEX idx_courses_created_at ON courses(created_at);
```

## 의존성

- ssac, stml 변경 없음 — fullend 내부 완결

## 검증 방법

1. `go build ./...` 성공
2. `go test ./...` 통과
3. `go run ./artifacts/cmd/fullend validate specs/dummy-lesson/` 실행 시 WARNING에 `→ 제안:` 줄 출력 확인
