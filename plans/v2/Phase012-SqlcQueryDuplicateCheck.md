✅ 완료

# Phase 012: sqlc 쿼리 이름 중복 검증

## 목표

`fullend validate` 단계에서 `db/queries/*.sql` 파일 간 `-- name:` 중복을 ERROR로 잡는다.
현재는 `sqlc generate` 시점에서야 실패하며, validate가 통과해버려 원인 파악이 늦어진다.

## 배경

sqlc는 모든 쿼리 파일을 하나의 Go 패키지로 합치므로, 서로 다른 파일에 동일한 `-- name: FindByID`가 있으면 컴파일 에러가 된다. fullend의 네이밍 컨벤션은 `ModelPrefix + Method` (예: `GigFindByID`, `UserCreate`)이지만, 이를 validate에서 강제하지 않고 있다.

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `internal/orchestrator/validate.go` | `validateDDL()` 내에서 sqlc 쿼리 이름 중복 검사 추가 |

## 구현 내용

`validateDDL()` 함수에 다음 로직 추가:

1. `db/queries/*.sql` 파일을 스캔하여 모든 `-- name: XXX :cardinality` 라인 파싱
2. 이름별로 출현 파일을 기록
3. 동일 이름이 2개 이상 파일에서 발견되면 ERROR 추가:
   ```
   db/queries: "FindByID" 이름이 중복됩니다 (users.sql, gigs.sql) — sqlc는 전역 네임스페이스이므로 ModelPrefix를 붙이세요 (예: UserFindByID, GigFindByID)
   ```

## 의존성

없음 (fullend 내부 변경만)

## 검증 방법

1. 중복 이름이 있는 쿼리 파일로 `fullend validate` → ERROR 출력 확인
2. prefix 적용 후 `fullend validate` → ERROR 없음 확인
3. `go test ./internal/orchestrator/...` 통과
