✅ 완료

# Phase 031: DDL → OpenAPI/SSaC 역방향 검증 + @archived

## 목표

DDL에 정의된 테이블·컬럼이 OpenAPI/SSaC에서 참조되지 않으면 WARNING을 출력한다. 보존용 테이블·컬럼은 `-- @archived` 태그로 WARNING을 억제한다.

## 배경

현재 crosscheck는 OpenAPI/SSaC → DDL 정방향(참조 대상이 DDL에 존재하는지)만 검사한다. 역방향이 없어서 DDL에 테이블을 만들어놓고 API에서 안 쓰는 "죽은 테이블"을 감지하지 못한다.

유지보수 과정에서 더 이상 API로 노출하지 않지만 고객 데이터 보존을 위해 남겨둔 테이블·컬럼은 `@archived`로 표시하여 WARNING 대상에서 제외한다.

SSaC 수정지시서017에서 SSaC 측 DDL 파서 수정을 요청했으나, `@archived`의 유일한 소비자가 fullend이므로 fullend에서 직접 구현하라는 회신을 받음.

## @archived 사용법

```sql
-- 테이블 수준: CREATE TABLE 직전 줄에 주석
-- @archived
CREATE TABLE legacy_notifications (
    id BIGSERIAL PRIMARY KEY,
    message TEXT
);

-- 컬럼 수준: 컬럼 정의 줄 끝에 인라인 주석
CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    old_category VARCHAR(100), -- @archived
    category VARCHAR(50) NOT NULL
);
```

## 검증 규칙

### 규칙 1: DDL 테이블 → OpenAPI/SSaC (WARNING)

DDL 테이블이 다음 중 하나라도 해당하면 "사용 중"으로 판정:
- SSaC `@model`에서 참조 (e.g. `@model Course.FindByID` → `courses` 테이블)
- SSaC `@result` 타입으로 참조 (e.g. `@result course Course` → `courses` 테이블)

`@archived` 테이블은 검증 대상에서 제외.

```
WARNING [DDL → SSaC]: DDL 테이블 "legacy_notifications"가 SSaC에서 참조되지 않습니다
  → 더 이상 사용하지 않는 테이블이면 DDL에 -- @archived를 추가하세요
```

### 규칙 2: DDL 컬럼 → OpenAPI 스키마 (WARNING)

DDL 테이블의 컬럼이 OpenAPI 스키마 properties에 없으면 WARNING. `@archived` 컬럼은 제외.

컬럼 매칭: DDL snake_case 컬럼명을 PascalCase로 변환하여 OpenAPI 스키마 properties와 비교 (e.g. `instructor_id` → `InstructorID`).

```
WARNING [DDL → OpenAPI]: DDL 컬럼 "courses.old_category"가 OpenAPI 스키마에 없습니다
  → 더 이상 사용하지 않는 컬럼이면 DDL에 -- @archived를 추가하세요
```

## 변경 항목

### A. DDL @archived 파서

| 파일 | 변경 |
|------|------|
| `internal/crosscheck/archived.go` (신규) | DDL 파일에서 `@archived` 테이블·컬럼 파싱. `ArchivedInfo` struct 반환 |

```go
type ArchivedInfo struct {
    Tables  map[string]bool            // "legacy_notifications" → true
    Columns map[string]map[string]bool // "courses" → {"old_category": true}
}

// ParseArchived는 DDL 디렉토리의 .sql 파일에서 @archived 태그를 파싱한다.
func ParseArchived(dbDir string) (*ArchivedInfo, error)
```

파싱 로직:
- `-- @archived` 단독 줄 → 다음 `CREATE TABLE`의 테이블명을 `Tables`에 등록
- 컬럼 정의 줄에 `-- @archived` 포함 → 해당 테이블·컬럼을 `Columns`에 등록

### B. DDL coverage crosscheck 규칙

| 파일 | 변경 |
|------|------|
| `internal/crosscheck/ddl_coverage.go` (신규) | DDL 테이블 → SSaC 참조 검증 |

```go
func CheckDDLCoverage(
    st *ssacvalidator.SymbolTable,
    funcs []ssacparser.ServiceFunc,
    doc *openapi3.T,
    archived *ArchivedInfo,
) []CrossError
```

- 규칙 1 (테이블): `st.DDLTables`의 각 테이블에 대해 `archived.Tables`에 있으면 스킵, SSaC 참조 집합에 없으면 WARNING
- 규칙 2 (컬럼): 사용 중인 테이블의 각 컬럼에 대해 `archived.Columns`에 있으면 스킵, OpenAPI 스키마 properties에 없으면 WARNING

### C. crosscheck.go Run() 연결

| 파일 | 변경 |
|------|------|
| `internal/crosscheck/crosscheck.go` | `CrossValidateInput`에 `ArchivedInfo` 필드 추가, `Run()`에서 `CheckDDLCoverage` 호출 |

### D. orchestrator에서 ArchivedInfo 전달

| 파일 | 변경 |
|------|------|
| `internal/orchestrator/orchestrator.go` | validate 단계에서 `ParseArchived(specsDir + "/db")` 호출, `CrossValidateInput.ArchivedInfo`에 전달 |

## 변경 파일 목록

| 파일 | 변경 |
|------|------|
| `internal/crosscheck/archived.go` | 신규 — `@archived` 파서 |
| `internal/crosscheck/ddl_coverage.go` | 신규 — DDL → SSaC coverage 검증 |
| `internal/crosscheck/crosscheck.go` | `CrossValidateInput`에 `ArchivedInfo` 추가, `Run()`에 호출 추가 |
| `internal/orchestrator/orchestrator.go` | `ParseArchived` 호출 + 전달 |

## 의존성

- SSaC 수정지시서017 회신 확인 완료 (fullend 이관)
- Phase030 완료 불필요 (병렬 진행 가능)

## 검증 방법

```bash
# 1. go test
go test ./internal/crosscheck/... -count=1

# 2. dummy-lesson에서 확인 (모든 테이블이 참조되므로 WARNING 없어야 함)
fullend validate specs/dummy-lesson

# 3. 임시로 @archived 없는 미사용 테이블을 DDL에 추가 → WARNING 확인
# 4. -- @archived 추가 → WARNING 사라지는지 확인
```
