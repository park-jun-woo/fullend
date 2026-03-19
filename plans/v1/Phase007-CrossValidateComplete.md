# ✅ 완료 — Phase 7 — 교차 검증 규칙 완성

## 목표
fullend 고유 교차 검증 규칙을 완성한다. ssac DDLTable에 추가된 ForeignKey, Index 정보를 활용하고, 함수명 ↔ operationId 매칭을 추가한다.

## 선행 조건
- Phase 4 완료
- ssac DDLTable에 ForeignKeys, Indexes 필드 추가 완료

## 변경 파일

| 파일 | 작업 |
|---|---|
| `artifacts/internal/crosscheck/openapi_ddl.go` | 수정. x-include FK 검증 강화, x-sort 인덱스 경고 추가 |
| `artifacts/internal/crosscheck/ssac_openapi.go` | 생성. 함수명 ↔ operationId 매칭 |
| `artifacts/internal/crosscheck/crosscheck.go` | 수정. ssac_openapi 규칙 추가 |

## 규칙 상세

### 규칙 1: x-include ↔ DDL FK (기존 강화)

현재: 테이블 존재만 확인 (naive)
변경: 실제 FK 관계로 연결된 테이블인지 확인

```
x-include:
  allowed: [Room]

→ "rooms" 테이블이 존재하는가? (현재)
→ 해당 엔드포인트의 주 테이블에서 "rooms"로 FK가 있는가? (강화)
```

주 테이블 추론: operationId로 ssac 함수를 찾고, 첫 @model의 모델명에서 테이블명 도출.
주 테이블을 추론할 수 없으면 기존 방식(테이블 존재 확인)으로 fallback.

### 규칙 2: x-sort ↔ DDL Index (신규)

x-sort.allowed 컬럼에 인덱스가 없으면 WARNING (성능 경고).

```
x-sort:
  allowed: [StartAt, CreatedAt]

→ start_at: idx_reservations_room_time에 포함 → OK
→ created_at: 인덱스 없음 → WARNING "x-sort column created_at has no index"
```

컬럼이 인덱스의 첫 번째 컬럼이거나 단독 인덱스에 포함되면 OK.
복합 인덱스의 2번째 이후 컬럼만 있으면 WARNING (선두 컬럼이 아님).

### 규칙 3: SSaC 함수명 → operationId 매칭 (신규)

모든 SSaC 함수명이 OpenAPI operationId로 존재하는지 확인.

```
SSaC: func Login(...)        → Operations["Login"] 존재? ✓
SSaC: func Orphan(...)       → Operations["Orphan"] 존재? ✗ ERROR
```

### 규칙 4: operationId → SSaC 함수명 (역방향, 신규)

모든 OpenAPI operationId에 대응하는 SSaC 함수가 존재하는지 확인.

```
OpenAPI: operationId: Login           → SSaC func Login 존재? ✓
OpenAPI: operationId: AdminDashboard  → SSaC func AdminDashboard 존재? ✗ WARNING
```

WARNING으로 처리: OpenAPI에 정의되어 있지만 아직 구현하지 않은 엔드포인트일 수 있음.

## 입력 데이터

Phase 4에서 이미 보존하는 중간 결과를 그대로 사용:
- `openapi3.T` — kin-openapi로 로드한 OpenAPI spec
- `ssacvalidator.SymbolTable` — DDLTables (ForeignKeys, Indexes 포함), Operations
- `[]ssacparser.ServiceFunc` — 파싱된 서비스 함수

## 검증 방법

dummy-study 프로젝트로 확인:
- x-include: reservations 테이블에서 rooms, users로 FK 존재 → 정상
- x-sort: StartAt은 인덱스 있음, CreatedAt은 인덱스 없음 → WARNING
- 함수명 ↔ operationId: 7개 함수 = 7개 operationId → 전체 매칭
- `go build` 성공, 기존 테스트 결과 유지
