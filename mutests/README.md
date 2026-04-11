# mutests

SSOT 변이 테스트. SSOT 명세에 의도적 결함을 주입하고 `fullend validate`가 검출하는지 확인한다.

## 실행 방법

1. 대상 프로젝트 specs를 복사한다
2. `.md` 파일의 "변경" 항목대로 SSOT를 수정한다
3. `fullend validate <specs-dir>`를 실행한다
4. "기대" 항목과 실제 출력을 비교하여 PASS/FAIL을 판정한다

## 케이스 형식

```markdown
### MUT-{SOURCE}-{TARGET}-{NNN}: 설명
- 대상: `specs/project/path/to/file`
- 변경: 원본 → 변이
- 기대: ERROR|WARNING — 기대하는 검출 메시지
- 결과: PASS|FAIL|SKIP — 실제 결과 + 비고
```

## 파일 구조

### 단일 SSOT 자체 검증

| 파일 | 대상 | 케이스 |
|------|------|--------|
| `config.md` | fullend.yaml | 2 |
| `openapi.md` | OpenAPI | 2 |
| `ddl.md` | DDL | 4 |
| `ssac.md` | SSaC | 49 |
| `stml.md` | STML | 3 |
| `states.md` | States | 1 |
| `policy.md` | Rego | 2 |
| `scenario.md` | Hurl | 2 |
| `func.md` | Func | 2 |
| `model.md` | Model | 2 |

### SSOT 간 교차 검증

| 파일 | 방향 | 케이스 |
|------|------|--------|
| `ssac-openapi.md` | SSaC ↔ OpenAPI | 5 |
| `ssac-ddl.md` | SSaC → DDL | 1 |
| `ssac-states.md` | SSaC → States | 1 |
| `ssac-func.md` | SSaC → Func | 4 |
| `ssac-config.md` | SSaC → Config | 2 |
| `ssac-queue.md` | SSaC Queue | 3 |
| `openapi-ddl.md` | OpenAPI ↔ DDL | 7 |
| `config-openapi.md` | Config ↔ OpenAPI | 2 |
| `states-ssac.md` | States → SSaC | 2 |
| `states-ddl.md` | States → DDL | 2 |
| `states-openapi.md` | States → OpenAPI | 2 |
| `ddl-ssac.md` | DDL → SSaC | 1 |
| `scenario-openapi.md` | Scenario → OpenAPI | 3 |
| `policy-ssac.md` | Policy ↔ SSaC | 2 |
| `policy-ddl.md` | Policy → DDL | 3 |
| `policy-config.md` | Policy → Config | 3 |
| `policy-states.md` | Policy → States | 2 |

## 보고서

| 파일 | 일시 | 결과 |
|------|------|------|
| `reports/mutest-report01.md` | 2026-03-17 | 114건: 87 PASS, 9 FAIL, 18 SKIP (90.6%) |
| `reports/mutest-report02.md` | — | 2차 실행 |
| `reports/mutest-report03.md` | — | 3차 실행 |

## 판정 기준

| 결과 | 의미 |
|------|------|
| PASS | 변이가 기대한 레벨(ERROR/WARNING)로 검출됨 |
| FAIL | 변이를 검출하지 못함 또는 잘못된 레벨로 검출 |
| SKIP | 대상 프로젝트에 해당 기능 없음 또는 미실행 |
| FAIL? | 검출 여부 불확실 (추가 확인 필요) |
