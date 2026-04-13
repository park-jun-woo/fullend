# Phase018 — DDLPipelineIntegration Report

2026-04-14 · DDL 위상정렬 + schema.sql 통합 산출 + auto nobody seed

## 목표 및 달성

| 목표 | 결과 |
|------|------|
| `fullend gen` 이 DDL 적용 순서 자동화 | ✅ |
| 단일 `schema.sql` 로 DB 초기화 | ✅ (gigbridge, zenflow) |
| `DEFAULT N FK` 패턴 nobody seed 자동 주입 | ✅ (gigbridge) |
| opt-in 플래그 (`backend.db.auto_nobody_seed`) | ✅ |

## 구현 개요

### 1. 위상정렬 (Kahn's algorithm)

FK 참조 그래프를 DAG 로 구성해 Kahn's algorithm 으로 위상정렬. 결정적 출력 (같은 depth 내 알파벳순). 순환 FK 는 에러.

**gigbridge 예시 순서**: users → gigs → proposals → transactions

### 2. schema.sql 통합

`pkg/generate/db/generate_schema.go` 가 FK 위상순서로 각 `<specs>/db/<table>.sql` 을 읽어 병합. 헤더 + `-- ---- <table> ----` 섹션 주석 + auto seed 블록.

**산출**: `<artifacts>/backend/db/schema.sql`

### 3. auto nobody seed

**감지**: 컬럼이 `DEFAULT <int> REFERENCES <ref>(id)` 이고 `<ref>` 에 id=N seed 없음.

**생성 규칙**:
- id 컬럼 → 지정 N
- CHECK IN 목록 컬럼 → 첫 값
- DEFAULT 가 있는 다른 컬럼 → SQL 생략 (DB 가 채움)
- 나머지 컬럼 → 타입 기반:
  - `int*`, `float64` → 0
  - `bool` → false
  - `email` 포함 문자열 → `nobody-<table>-<col>@autoseed.local` (UNIQUE 충돌 회피)
  - 기타 문자열 → `nobody-<table>-<col>`
- `ON CONFLICT DO NOTHING` 으로 idempotent

**gigbridge 예시**:
```sql
INSERT INTO users (id, email, password_hash, role, name) VALUES
(0, 'nobody-users-email@autoseed.local', 'nobody-users-password_hash', 'client', 'nobody-users-name')
ON CONFLICT DO NOTHING;
```

### 4. opt-in

`dummys/*/specs/fullend.yaml`:
```yaml
backend:
  db:
    auto_nobody_seed: true
```

기본 off. 기존 프로젝트가 의도적으로 sentinel 시드 관리하는 경우 보호.

### 5. sentinel 검증 상호작용

기존 `internal/orchestrator/check_sentinel_record.go` 가 FK DEFAULT 0 + sentinel INSERT 부재를 ERROR 로 차단했었음. auto seed 활성 시 이 검증을 skip (validate_ddl.go 에서 manifest 플래그 조회). Phase016 X-79 (WARNING) 는 유지.

## 실측 검증 (docker pg)

### gigbridge

```bash
$ psql -f schema.sql
CREATE DATABASE / CREATE TABLE x4 / CREATE INDEX x5 / INSERT 0 1
$ hurl smoke.hurl
Success (12 request(s) in 192 ms) - 12/12 통과 ✓
```

### zenflow

```bash
$ psql -f schema.sql
CREATE DATABASE / CREATE TABLE x5 / CREATE INDEX x5
# auto seed 없음 (DEFAULT FK 패턴 없음)
$ INSERT INTO organizations ... (여전히 수동 — zenflow 는 root anchor org 필요)
```

## filefunc

- 신규 파일 20개 모두 F1/F2/A1/A3/A10/Q1 준수
- baseline 37 유지

## 다음 Phase

v5 로드맵의 버그·검증·자동화 축은 **본 Phase 로 마감**.

남은 v5 트랙: Phase014 (internal/* 삭제), Phase015 (template/잔여 정리). 두 Phase 모두 Phase018 에 의존 없음.

**Phase019 후보** (v6 영역, 별도 로드맵):
- zenflow ListWorkflows 500 (Phase017 에서 기록) — SSaC `@get` 의 org 스코프 자동 주입 설계
- hurl smoke 재생성 (schema.sql 사용 버전 업데이트)
- Phase016 X-79 규칙의 auto seed 인식 튜닝
