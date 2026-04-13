# ✅ Phase018 — DDLPipelineIntegration (완료: 2026-04-14, 미커밋)

> DDL 적용 자동화. 위상정렬된 `schema.sql` 통합 산출 + `DEFAULT N FK` 패턴 용 nobody seed 자동 주입.

## 배경

Phase017 실측에서 드러난 인프라 결함:

- `dummys/*/specs/db/*.sql` 을 알파벳 순 적용 시 FK 위반 (gigs 가 users 보다 먼저 와서 실패)
- `gigs.freelancer_id BIGINT NOT NULL DEFAULT 0 REFERENCES users(id)` 패턴 — id=0 사용자가 있어야 INSERT 동작. 현재는 spec 이 수동으로 `INSERT (id=0, ...)` seed 를 작성해야 하고 그마저 CHECK 와 충돌

## 목표

**fullend gen 이 DDL 적용 순서와 seed 를 모두 책임**. 사용자는 `psql -f artifacts/backend/db/schema.sql` 한 번으로 DB 초기화 가능.

## 전제

- Phase013 완료 ✅
- Phase016/017 과 **독립** (순서 자유). 단 Phase017 완료 후 진행하면 Phase017 1-5 에서 임시 spec 수정을 되돌릴 수 있음

## 범위

### 1. DDL 위상정렬

**대상**: `pkg/parser/ddl` 가 이미 CREATE TABLE 의 FK 참조를 파싱. 이를 DAG 로 구성해 위상정렬 결과로 테이블 적용 순서 결정.

**작업**:
- `pkg/generate/db/` (또는 적절한 위치) 신설
- `sortTablesByFK(tables []ddl.Table) []string` — 위상정렬
- 사이클 감지 시 에러 (순환 FK 는 PostgreSQL 에선 deferrable 로 해결 가능하지만 v5 범위 밖)

### 2. `schema.sql` 통합 산출

**산출 경로**: `artifacts/backend/db/schema.sql`

**내용**:
1. 위상정렬된 순서로 각 CREATE TABLE + INDEX 이어 붙임
2. seed INSERT 는 모든 CREATE 이후 집합
3. 헤더 주석: 생성 시각, 원본 spec 경로

**작업**:
- `pkg/generate/db/generate_schema.go` 신설
- `orchestrator` 에 `gen_schema` 스텝 추가 (sqlc 와 별개)

### 3. `DEFAULT N FK` 패턴 용 nobody seed 자동 주입

**감지 조건**:
- CREATE TABLE 컬럼이 `DEFAULT <N>` 이고 같은 컬럼에 `REFERENCES <table>(id)` 존재
- N 이 정수 (보통 0)

**자동 주입**:
- 참조되는 테이블 (예: users) 에 `INSERT INTO users (id, ...) VALUES (N, ...)` 주입
- 필요한 컬럼값은 타입 기반 zero-value (string 은 "nobody@system" 같은 placeholder, NOT NULL 이면 CHECK 규칙 충족 값)
- 사용자가 이미 같은 id 를 seed 한 경우 skip (ON CONFLICT DO NOTHING 사용)

**주의**: 
- CHECK 제약 위반하지 않는 값 선택 (예: role 의 경우 CHECK 목록 첫 항목)
- 자동 주입은 opt-in — manifest 에 `db.auto_nobody_seed: true` 같은 설정

### 4. `fullend gen` 단계 통합

**변경**:
- orchestrator 의 codegen 단계에 `gen_schema` 신설
- 리포트: `✓ schema-gen    schema.sql generated (N tables, M seeds)`

### 5. dummy spec 정리 (Phase017 1-5 후속)

**대상**:
- `dummys/gigbridge/specs/db/users.sql` 의 수동 nobody INSERT 제거 (auto seed 사용)
- `dummys/zenflow/specs/db/` 도 동일 검토

---

## 작업 순서

### Step 1. pkg/parser/ddl 능력 확인

FK 정보가 `ddl.Table` 에 이미 파싱되는지 확인. 없으면 파서 확장.

### Step 2. 위상정렬 + schema.sql

`pkg/generate/db/` 신설, 위상정렬 + 병합 로직.

### Step 3. nobody seed 자동 주입

`DEFAULT N FK` 감지 + INSERT 생성. CHECK 목록 참조.

### Step 4. orchestrator 통합

`gen_schema` 단계 추가. 리포트 출력.

### Step 5. 검증

```bash
rm -rf dummys/gigbridge/artifacts
FULLEND_LOCAL_PATH=$(pwd) go run ./cmd/fullend gen dummys/gigbridge/specs dummys/gigbridge/artifacts
# 단일 파일로 DB 초기화
PGPASSWORD=test psql -h localhost -p 15432 -U postgres -d gigbridge -f dummys/gigbridge/artifacts/backend/db/schema.sql
# smoke 통과 확인
```

### Step 6. dummy spec 정리 + 커밋

---

## 주의사항

### R1. sqlc 호환 유지

sqlc 가 여전히 `specs/db/*.sql` 을 입력으로 사용하므로 원본 개별 파일 구조는 유지. `schema.sql` 은 추가 산출, 대체 아님.

### R2. auto seed 는 opt-in

기본 off. `fullend.yaml` 에 `backend.db.auto_nobody_seed: true` 로 활성. 기본 off 이유: 기존 프로젝트가 의도적으로 id=0 시드를 안 쓸 수도 있음.

### R3. CHECK 제약 대응

seed 값 선택 시 같은 컬럼의 CHECK 목록에서 유효 값 선택. 목록 없으면 빈 문자열 / 0.

### R4. 순환 FK

v5 범위 밖. 감지 시 ERROR 반환, 사용자가 수동으로 schema.sql 작성하도록 안내.

### R5. Phase017 1-5 되돌리기

본 Phase 완료 시 gigbridge users.sql 의 수동 nobody INSERT 를 제거. auto seed 가 대체. 커밋 분리.

---

## 완료 조건 (Definition of Done)

- [x] `pkg/generate/db/` 신설 — 위상정렬 (Kahn) + schema.sql 생성 + auto seed
- [x] `DEFAULT N FK` 패턴 감지 + nobody seed INSERT 자동 주입
- [x] `fullend.yaml` 의 `backend.db.auto_nobody_seed` opt-in 플래그
- [x] orchestrator `gen_schema` 단계 통합 (`✓ schema-gen    schema.sql generated (N tables, M seeds)`)
- [x] gigbridge `schema.sql` 1회 실행 → CREATE 4 + INSERT 1 (auto nobody) + INDEX 5 → smoke 12/12
- [x] zenflow `schema.sql` — CREATE 5 + INDEX (DEFAULT FK 패턴 없음 → auto seed 생략)
- [x] Phase017 1-5 임시 spec 수정 되돌리기 (gigbridge users.sql 수동 INSERT 제거, auto seed 가 대체)
- [x] `go build / vet / test ./pkg/...` 통과
- [x] `filefunc validate` baseline 37 유지 (신규 파일 위반 0)
- [x] sentinel 검증과 auto seed 상호작용: manifest.backend.db.auto_nobody_seed=true 시 기존 `check_sentinel_record` skip (false positive 방지)
- [ ] 커밋: `feat(generate): DDL 위상정렬 + schema.sql 통합 + auto nobody seed (Phase018)`

### 부산물

- `pkg/generate/db/` 16 파일 신설: 위상정렬 (`sort_tables_by_fk`, `build_fk_graph`, `topo_sort`, `pick_zero_degree`, `pending_tables`, `decrement_dependents`, `apply_zero_degree_batch`), schema 조립 (`generate_schema`, `config`, `assemble_schema_sql`, `read_ddl_file_for_table`, `index_tables_by_name`), auto seed (`build_auto_nobody_seeds`, `collect_required_seed_ids`, `sorted_seed_keys`, `parse_seed_key`, `seed_already_exists`, `build_seed_insert_stmt`, `seed_value_for`, `nobody_placeholder_for`)
- `pkg/parser/manifest/db_config.go` + `backend.go` (DB 필드 추가)
- `internal/orchestrator/gen_schema.go`, `run_codegen_steps.go` 등록
- `internal/orchestrator/{check_column_line,check_ddl_nullable_columns,validate_ddl}.go` — auto seed 활성 시 sentinel 검증 skip
- `dummys/gigbridge/specs/fullend.yaml` + `db/users.sql` — opt-in 플래그 + 수동 seed 제거

## 의존

- Phase013 완료 ✅
- Phase017 완료 **권장** (1-5 되돌리기 단계 때문). 독립 진행도 가능하나 spec 이 임시 수정된 상태에서 작업하게 됨
- Phase016 선행 완료 시: 본 Phase auto seed 로 Phase016 의 2-2 규칙 의미 재조정 (false positive 튜닝) 필요

## 다음 Phase

v5 로드맵의 **버그·검증·자동화 축은 본 Phase 로 마감**. 이후는 Phase014/015 (구조 정리 마감) + v6 (기능 확장).
