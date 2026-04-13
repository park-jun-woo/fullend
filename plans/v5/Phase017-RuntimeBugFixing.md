# Phase017 — RuntimeBugFixing

> 2026-04-14 docker + hurl 실측에서 드러난 **런타임 버그** 수정. 검증 규칙 보강 (Phase016) 과 생성 기능 enhancement (Phase018) 는 분리.

## 배경

Phase013 까지 완료 후 실제 파이프라인 실행 (docker PostgreSQL + generated backend + hurl smoke) 로 노출된 결함:

- gigbridge smoke 12/12 통과 (수동 DB 패치 후)
- zenflow smoke → **ActivateWorkflow 에서 403**
- `OPA_POLICY_PATH` 디렉토리 미지원
- 생성 main.go DSN 기본값 환경 의존적
- gigbridge `users.sql` seed INSERT 가 CHECK 위반

## 목표

**수동 DB 패치 없이** gigbridge / zenflow smoke 전 구간 통과.

## 전제

- Phase013 완료 ✅
- docker + PostgreSQL + hurl 테스트 환경

Phase014/015/017/018 과 **독립**.

---

## 범위

### 1-1. zenflow `ActivateWorkflow` 403 원인 추적 + 수정

**증상**: register → login → CreateWorkflow (통과) → ActivateWorkflow 에서 403.

**추정 원인** (실조사 필요):
- **(a)** JWT claim 디코딩 시 `org_id` 가 float64 로 복원 → OPA Rego `data.owners.workflow[id] == input.claims.org_id` 비교 실패
- **(b)** `loadOwners` 반환 map 구조 `owners[resource][id] = ownerID` 가 Rego 의 `data.owners.workflow[input.resource_id]` 기대값과 타입 불일치
- **(c)** OPA 가 받는 `input.resource_id` 가 string ("2") vs key 가 int64 (2)

**작업**:
1. 서버 debug 로그 + OPA trace 로 실제 실패 지점 식별
2. `pkg/auth` JWT decode 경로에서 숫자 claim 타입 확인 (`jwt.RegisteredClaims` vs 커스텀)
3. `pkg/authz/load_owners.go` 반환 map 의 key/value 타입 확인
4. 원인별 수정
   - (a) → JWT parser 에 숫자 claim 명시적 int64 cast
   - (b) → loadOwners 의 key 를 string 으로 통일 (Rego 관례)
   - (c) → Rego input 조립 시 `fmt.Sprint(id)` 사용 일관성

**완료 기준**: ActivateWorkflow 200 반환, 이후 smoke 단계 진행.

### 1-2. OPA_POLICY_PATH 디렉토리 지원 + fallback

**증상**: env 에 디렉토리 지정 시 `is a directory` 에러. `.rego` 파일 경로만 수용.

**수정 대상**: `pkg/authz/init_authz.go`

**작업**:
- env 값이 디렉토리면 `filepath.Glob("*.rego")` 로 전체 로드
- env 미지정 시 실행 바이너리 기준 `./internal/authz/` 또는 현재 디렉토리 자동 탐색
- 생성된 `cmd/main.go` 에 필수 env 안내 주석 추가 (generator 수정)

**완료 기준**: env 지정 없이도 기본 경로에서 .rego 로드, 서버 기동 성공.

### 1-3. 생성 main.go DSN 기본값 개선

**증상**: 기본 dsn `postgres://localhost:5432/app?sslmode=disable` — 프로젝트 무관 하드코딩.

**수정 대상**: `pkg/generate/gogin/main_template.go` (또는 동등 위치)

**작업**:
- 기본값 우선순위: `DATABASE_URL` env > `--dsn` flag > fallback
- fallback 을 `postgres://localhost:5432/{모듈명}?sslmode=disable` 로 (모듈명 = `fs.Manifest.Backend.Module` 마지막 세그먼트)
- `DATABASE_URL` env 처리 코드 emit

**완료 기준**: 생성된 main.go 에 `os.Getenv("DATABASE_URL")` 분기 존재, dummy 별 기본 DB 이름 다름.

### 1-5. gigbridge `users.sql` seed CHECK 위반 해소

**증상**: `INSERT ... role='system'` 이 `CHECK (role IN ('client','freelancer','admin'))` 위반.

**수정 옵션**:
- **(a)** spec 수정 — `role='admin'` 으로 변경 (단기)
- **(b)** spec 에서 INSERT 제거 — fullend 가 nobody seed 자동 주입 (Phase018 의존)

**본 Phase 선택**: **(a)**. Phase018 완료 시 spec 의 INSERT 자체를 제거로 재정리.

**완료 기준**: `psql -f users.sql` 성공, DDL 적용 시 에러 없음.

---

## 작업 순서

### Step 1. docker PostgreSQL 기동 + 1-1 원인 추적

```bash
docker run -d --name fullend-pg-test -e POSTGRES_PASSWORD=test -p 15432:5432 postgres:16-alpine
# zenflow 재생성, 기동, hurl 실행하며 OPA_DEBUG=1 등으로 trace 수집
```

추적 산출물: `reports/phase016-zenflow-403-trace.md` (임시)

### Step 2. 원인별 1-1 수정

Step 1 결과에 따라 `pkg/auth` or `pkg/authz` 수정. 테스트 추가 (claim 타입 보존, owners map key 타입).

### Step 3. 1-2 / 1-3 / 1-5

독립적 수정 3건. 각각 단위 검증.

### Step 4. 통합 smoke

```bash
# 재생성 + 기동
rm -rf dummys/gigbridge/artifacts dummys/zenflow/artifacts
FULLEND_LOCAL_PATH=$(pwd) go run ./cmd/fullend gen dummys/gigbridge/specs dummys/gigbridge/artifacts
FULLEND_LOCAL_PATH=$(pwd) go run ./cmd/fullend gen dummys/zenflow/specs dummys/zenflow/artifacts

# DB 초기화 (DDL 순서 수동 — Phase018 에서 자동화)
for f in users gigs proposals transactions; do
  PGPASSWORD=test psql -h localhost -p 15432 -U postgres -d gigbridge -f dummys/gigbridge/specs/db/${f}.sql
done

# 서버 기동 + hurl
# gigbridge smoke 12/12 + zenflow smoke 전체 통과 확인
```

### Step 5. 정리 + 커밋

- docker 컨테이너 정리
- `reports/bugfix-phase016.md` 생성 (원인 기록 + 수정 내용)
- 커밋: `fix(runtime): Phase017 런타임 버그 4종 수정`

---

## 주의사항

### R1. 원인 조사 우선

1-1 은 실제 원인 모름. 추정으로 수정 금지 — trace 로 핀포인트 확정 후 수정.

### R2. gigbridge 회귀 금지

1-1 수정이 JWT/authz 공통 경로를 건드리므로 gigbridge smoke 재확인 필수.

### R3. 1-5 는 Phase018 에서 재조정

1-5 (a) 안으로 spec 을 수정했다가 Phase018 (nobody seed 자동 주입) 완료 시 해당 INSERT 를 spec 에서 제거할 것.

### R4. 회귀 검사

각 수정 후 `go build / vet / test ./pkg/...` 통과 확인.

---

## 완료 조건 (Definition of Done)

- [ ] 1-1 zenflow ActivateWorkflow 403 원인 확정 + 수정
- [ ] 1-2 OPA_POLICY_PATH 디렉토리 지원 + fallback
- [ ] 1-3 DSN 기본값 개선 (`DATABASE_URL` env + 모듈명 fallback)
- [ ] 1-5 gigbridge users.sql seed CHECK 위반 해소
- [ ] 수동 DB 패치 없이 gigbridge smoke 12/12 통과
- [ ] zenflow smoke 전체 통과 (register → workflow 플로우 완주)
- [ ] `go build / vet / test ./pkg/...` 통과
- [ ] `filefunc validate` — 신규 파일 위반 0
- [ ] `reports/bugfix-phase016.md` 생성
- [ ] 커밋: `fix(runtime): Phase017 런타임 버그 4종 수정`

## 의존

- Phase013 완료 ✅
- Phase016 선행 **권장** — validate 결과가 spec 결함 리스트를 자동으로 내려줘서 조사 범위 축소
- Phase014/015/018 과 **독립**

## 다음 Phase

- **Phase018** — DDLPipelineIntegration (schema.sql 통합 + nobody seed 자동 주입, 본 Phase 1-5 의 임시 spec 수정 되돌림)
