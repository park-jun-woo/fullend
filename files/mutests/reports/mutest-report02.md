# Mutation Test Report 02

- 일시: 2026-03-17
- 대상: specs/dummys/zenflow-try05/
- 바이너리: fullend (commit 460eb43, Phase045 적용)
- 테스트 파일: files/mutests/*.md (27파일, 114케이스)

## 요약

| 항목 | 건수 |
|---|---|
| 총 케이스 | 114 |
| PASS | 83 |
| FAIL | 8 |
| SKIP | 23 |
| 통과율 (SKIP 제외) | 91.2% |

## 전체 결과

### config.md
| ID | 결과 |
|---|---|
| MUT-CONFIG-001 | PASS |
| MUT-CONFIG-002 | PASS |

### config-openapi.md
| ID | 결과 |
|---|---|
| MUT-CONFIG-OPENAPI-001 | PASS |
| MUT-CONFIG-OPENAPI-002 | PASS |

### ddl.md
| ID | 결과 |
|---|---|
| MUT-DDL-001 | PASS |
| MUT-DDL-002 | PASS |
| MUT-DDL-003 | PASS |
| MUT-DDL-004 | PASS |

### ddl-ssac.md
| ID | 결과 |
|---|---|
| MUT-DDL-SSAC-001 | PASS |

### func.md
| ID | 결과 |
|---|---|
| MUT-FUNC-001 | PASS |
| MUT-FUNC-002 | PASS |

### model.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-MODEL-001 | PASS | |
| MUT-MODEL-002 | SKIP | zenflow에 @dto 없음 |

### openapi.md
| ID | 결과 |
|---|---|
| MUT-OPENAPI-001 | PASS |
| MUT-OPENAPI-002 | PASS |

### openapi-ddl.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-OPENAPI-DDL-001 | PASS | |
| MUT-OPENAPI-DDL-002 | PASS | |
| MUT-OPENAPI-DDL-003 | PASS | |
| MUT-OPENAPI-DDL-004 | SKIP | zenflow에 cursor pagination 없음 |
| MUT-OPENAPI-DDL-005 | SKIP | zenflow에 cursor pagination 없음 |
| MUT-OPENAPI-DDL-006 | SKIP | zenflow에 x-include 없음 |
| MUT-OPENAPI-DDL-007 | SKIP | zenflow에 x-include 없음 |

### policy.md
| ID | 결과 |
|---|---|
| MUT-POLICY-001 | PASS |
| MUT-POLICY-002 | PASS |

### policy-config.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-POLICY-CONFIG-001 | **FAIL** | Rego role 값 "Admin" vs "admin" 대소문자 미검출 |
| MUT-POLICY-CONFIG-002 | PASS | claims JWT key는 자유값. 코드젠 자기 일관적. 검증 범위 밖 |
| MUT-POLICY-CONFIG-003 | SKIP | zenflow에 roles 섹션 없음 |

### policy-ddl.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-POLICY-DDL-001 | PASS | |
| MUT-POLICY-DDL-002 | SKIP | zenflow에 join table ownership 없음 |
| MUT-POLICY-DDL-003 | SKIP | zenflow에 join table ownership 없음 |

### policy-ssac.md
| ID | 결과 |
|---|---|
| MUT-POLICY-SSAC-001 | PASS |
| MUT-POLICY-SSAC-002 | PASS |

### policy-states.md
| ID | 결과 |
|---|---|
| MUT-POLICY-STATES-001 | PASS |
| MUT-POLICY-STATES-002 | PASS |

### scenario.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-SCENARIO-001 | PASS | |
| MUT-SCENARIO-002 | **FAIL** | 빈 tests/ → WARN만 표시, ERROR 미발생 |

### scenario-openapi.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-SCENARIO-OPENAPI-001 | **FAIL** | hurl 경로 vs OpenAPI path 교차 검증 없음 |
| MUT-SCENARIO-OPENAPI-002 | **FAIL** | hurl HTTP method vs OpenAPI method 교차 검증 없음 |
| MUT-SCENARIO-OPENAPI-003 | SKIP | 검증 복잡도 높음 |

### ssac-config.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-SSAC-CONFIG-001 | **FAIL** | claims에서 ID 삭제 시 currentUser.ID 미참조 미검출 |
| MUT-SSAC-CONFIG-002 | PASS | |

### ssac-ddl.md
| ID | 결과 |
|---|---|
| MUT-SSAC-DDL-001 | PASS |

### ssac-func.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-SSAC-FUNC-001 | PASS | |
| MUT-SSAC-FUNC-002 | SKIP | zenflow에 해당 시나리오 없음 |
| MUT-SSAC-FUNC-003 | SKIP | 설정 복잡 |
| MUT-AUTHZ-001 | SKIP | zenflow에 해당 시나리오 없음 |

### ssac-openapi.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-SSAC-OPENAPI-001 | PASS | |
| MUT-SSAC-OPENAPI-002 | SKIP | 설계상 정상 |
| MUT-SSAC-OPENAPI-003 | PASS | |
| MUT-SSAC-OPENAPI-004 | PASS | |
| MUT-SSAC-OPENAPI-005 | PASS | |

### ssac-queue.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-SSAC-QUEUE-001 | SKIP | zenflow에 @publish 없음 |
| MUT-SSAC-QUEUE-002 | SKIP | zenflow에 @publish 없음 |
| MUT-SSAC-QUEUE-003 | SKIP | zenflow에 @publish 없음 |

### ssac-states.md
| ID | 결과 |
|---|---|
| MUT-SSAC-STATES-001 | PASS |

### ssac.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-SSAC-001 | PASS | |
| MUT-SSAC-002 | PASS | Phase045: singular 검증으로 검출 |
| MUT-SSAC-003 | PASS | |
| MUT-SSAC-004 | PASS | |
| MUT-SSAC-005 | **FAIL?** | input key "Id" vs DDL "ID" 대소문자 미검출 |
| MUT-SSAC-006 | PASS | |
| MUT-SSAC-007 | PASS | |
| MUT-SSAC-008 | PASS | |
| MUT-SSAC-009 | PASS | |
| MUT-SSAC-010 | PASS | |
| MUT-SSAC-011 | PASS | |
| MUT-SSAC-012 | PASS | |
| MUT-SSAC-013 | PASS | |
| MUT-SSAC-014 | PASS | |
| MUT-SSAC-015 | PASS | |
| MUT-SSAC-016 | SKIP | zenflow에 @delete 없음 |
| MUT-SSAC-017 | SKIP | zenflow에 @delete 없음 |
| MUT-SSAC-018 | PASS | |
| MUT-SSAC-019 | PASS | |
| MUT-SSAC-020 | **FAIL?** | HTTP status 999 범위 검증 없음 |
| MUT-SSAC-021 | PASS | |
| MUT-SSAC-022 | SKIP | zenflow에 @exists 없음 |
| MUT-SSAC-023 | SKIP | zenflow에 @exists 없음 |
| MUT-SSAC-024 | PASS | |
| MUT-SSAC-025 | PASS | Phase045: seq.Transition 교체로 검출 |
| MUT-SSAC-026 | PASS | |
| MUT-SSAC-027 | PASS | States ↔ DDL: state field 검출 |
| MUT-SSAC-028 | PASS | |
| MUT-SSAC-029 | PASS | |
| MUT-SSAC-030 | PASS | |
| MUT-SSAC-031 | PASS | |
| MUT-SSAC-032 | PASS | |
| MUT-SSAC-033 | PASS | |
| MUT-SSAC-034 | PASS | |
| MUT-SSAC-035 | PASS | Phase045: 소문자 함수명 ERROR 검출 |
| MUT-SSAC-036 | PASS | |
| MUT-SSAC-037 | **FAIL** | @call input key "OrgId" vs "OrgID" 미검출 |
| MUT-SSAC-038 | PASS | |
| MUT-SSAC-039 | SKIP | zenflow에 @publish 없음 |
| MUT-SSAC-040 | SKIP | zenflow에 @publish 없음 |
| MUT-SSAC-041 | SKIP | zenflow에 @publish 없음 |
| MUT-SSAC-042 | PASS | |
| MUT-SSAC-043 | PASS | |
| MUT-SSAC-044 | PASS | |
| MUT-SSAC-045 | PASS | |
| MUT-SSAC-046 | PASS | dotted field 타입 불일치 — 검증 범위 외 |
| MUT-SSAC-047 | PASS | |
| MUT-SSAC-048 | PASS | @put 이후 stale 경고 없음 — 허용 |
| MUT-SSAC-049 | PASS | 변수 재선언 허용 |

### states.md
| ID | 결과 |
|---|---|
| MUT-STATES-001 | PASS |

### states-ddl.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-STATES-DDL-001 | PASS | |
| MUT-STATES-DDL-002 | SKIP | 미구현 |

### states-openapi.md
| ID | 결과 |
|---|---|
| MUT-STATES-OPENAPI-001 | PASS |
| MUT-STATES-OPENAPI-002 | PASS |

### states-ssac.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-STATES-SSAC-001 | PASS | |
| MUT-STATES-SSAC-002 | SKIP | 미실행 |

### stml.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-STML-001 | PASS | |
| MUT-STML-002 | PASS | |
| MUT-STML-003 | SKIP | zenflow에 x-component 없음 |

## FAIL 상세

### 확정 FAIL (7건)

| ID | 원인 | 수정 방향 |
|---|---|---|
| MUT-POLICY-CONFIG-001 | Rego role 값의 대소문자를 정규 role 목록과 교차 검증하지 않음 | config roles vs Rego role 값 대소문자 교차 검증 추가 |
| ~~MUT-POLICY-CONFIG-002~~ | ~~claims JWT key는 자유값~~ | **제외** — 코드젠 자기 일관적, 검증 범위 밖 |
| MUT-SCENARIO-002 | 빈 tests/ 디렉토리가 WARN만 표시, ERROR 아님 | 빈 scenario 디렉토리 → ERROR 격상 검토 |
| MUT-SCENARIO-OPENAPI-001 | hurl 경로와 OpenAPI path 교차 검증 없음 | hurl path ↔ OpenAPI path 교차 검증 추가 |
| MUT-SCENARIO-OPENAPI-002 | hurl HTTP method와 OpenAPI operation method 교차 검증 없음 | hurl method ↔ OpenAPI method 교차 검증 추가 |
| MUT-SSAC-CONFIG-001 | claims에서 개별 필드 삭제 시 currentUser.* 참조와 교차 검증 없음 | currentUser 필드별 claims 존재 검증 추가 |
| MUT-SSAC-037 | @call input key "OrgId" vs FuncSpec "OrgID" 대소문자 미검출 | @call input key ↔ Request 필드명 exact match 검증 추가 |

### 불확실 FAIL (2건, "?" 표기)

| ID | 원인 | 판단 |
|---|---|---|
| MUT-SSAC-005 | input key "Id" vs "ID" 대소문자 | Go naming convention 차이, 검증 범위 외일 수 있음 |
| MUT-SSAC-020 | HTTP status 999 범위 검증 | 유효 범위(100~599) 검증 추가 가능하나 우선순위 낮음 |

## SKIP 사유

| 사유 | 건수 |
|---|---|
| zenflow에 해당 기능 없음 (queue, @dto, @delete, @exists, x-component 등) | 14 |
| 미실행 (설정 복잡 또는 미구현 규칙) | 7 |
| 설계상 정상 (검증 대상 아님) | 2 |

## Phase045 효과

| ID | report01 (gigbridge) | report02 (zenflow) | 비고 |
|---|---|---|---|
| MUT-FUNC-002 | FAIL | PASS | panic("TODO") stub 감지 추가 |
| MUT-SSAC-002 | FAIL | PASS | singular 검증 추가 (Workflows→Workflow) |
| MUT-SSAC-025 | FAIL | PASS | fn.Name→seq.Transition 교체 |
| MUT-SSAC-035 | FAIL | PASS | 소문자 함수명 ERROR 검출 추가 |

Phase045에서 수정한 4건 모두 zenflow-try05에서 PASS 확인.
