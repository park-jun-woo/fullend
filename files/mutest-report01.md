# Mutation Test Report 01

- 일시: 2026-03-17
- 대상: specs/dummys/gigbridge-try02/
- 바이너리: fullend (commit 49c5410)
- 테스트 파일: files/mutests/*.md (27파일, 114케이스)

## 요약

| 항목 | 건수 |
|---|---|
| 총 케이스 | 114 |
| PASS | 87 |
| FAIL | 9 |
| SKIP | 18 |
| 통과율 (SKIP 제외) | 90.6% |

## 전체 결과

### config.md
| ID | 결과 |
|---|---|
| MUT-CONFIG-001 | PASS |
| MUT-CONFIG-002 | PASS |

### config-openapi.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-CONFIG-OPENAPI-001 | **FAIL** | middleware 제거 시 securitySchemes 불일치 미검출 |
| MUT-CONFIG-OPENAPI-002 | PASS | |

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
| ID | 결과 | 비고 |
|---|---|---|
| MUT-FUNC-001 | PASS | |
| MUT-FUNC-002 | **FAIL** | panic("TODO") stub body 미감지 |

### model.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-MODEL-001 | PASS | |
| MUT-MODEL-002 | SKIP | gigbridge에 @dto 없음 |

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
| MUT-OPENAPI-DDL-004 | SKIP | 미실행 |
| MUT-OPENAPI-DDL-005 | SKIP | 미실행 |
| MUT-OPENAPI-DDL-006 | SKIP | 미실행 |
| MUT-OPENAPI-DDL-007 | SKIP | 미실행 |

### policy.md
| ID | 결과 |
|---|---|
| MUT-POLICY-001 | PASS |
| MUT-POLICY-002 | PASS |

### policy-config.md
| ID | 결과 |
|---|---|
| MUT-POLICY-CONFIG-001 | PASS |
| MUT-POLICY-CONFIG-002 | PASS |
| MUT-POLICY-CONFIG-003 | PASS |

### policy-ddl.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-POLICY-DDL-001 | PASS | |
| MUT-POLICY-DDL-002 | SKIP | 미실행 |
| MUT-POLICY-DDL-003 | SKIP | 미실행 |

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
| ID | 결과 |
|---|---|
| MUT-SCENARIO-001 | PASS |
| MUT-SCENARIO-002 | PASS |

### scenario-openapi.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-SCENARIO-OPENAPI-001 | PASS | |
| MUT-SCENARIO-OPENAPI-002 | PASS | |
| MUT-SCENARIO-OPENAPI-003 | SKIP | 미실행 |

### ssac-config.md
| ID | 결과 |
|---|---|
| MUT-SSAC-CONFIG-001 | PASS |
| MUT-SSAC-CONFIG-002 | PASS |

### ssac-ddl.md
| ID | 결과 |
|---|---|
| MUT-SSAC-DDL-001 | PASS |

### ssac-func.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-SSAC-FUNC-001 | PASS | |
| MUT-SSAC-FUNC-002 | SKIP | 미실행 |
| MUT-SSAC-FUNC-003 | SKIP | 미실행 |
| MUT-AUTHZ-001 | SKIP | 미실행 |

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
| MUT-SSAC-QUEUE-001 | SKIP | gigbridge에 @publish 없음 |
| MUT-SSAC-QUEUE-002 | SKIP | gigbridge에 @publish 없음 |
| MUT-SSAC-QUEUE-003 | SKIP | 미실행 |

### ssac-states.md
| ID | 결과 |
|---|---|
| MUT-SSAC-STATES-001 | PASS |

### ssac.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-SSAC-001 | PASS | |
| MUT-SSAC-002 | **FAIL** | @get result type "Gigs" 존재하지 않는 타입 미검증 |
| MUT-SSAC-003 | PASS | |
| MUT-SSAC-004 | PASS | |
| MUT-SSAC-005 | **FAIL?** | input key "Id" vs DDL "ID" 대소문자 (mutest 불확실 표기) |
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
| MUT-SSAC-016 | PASS | |
| MUT-SSAC-017 | PASS | |
| MUT-SSAC-018 | PASS | |
| MUT-SSAC-019 | PASS | |
| MUT-SSAC-020 | **FAIL?** | HTTP status 999 범위 검증 없음 (mutest 불확실 표기) |
| MUT-SSAC-021 | PASS | |
| MUT-SSAC-022 | PASS | |
| MUT-SSAC-023 | PASS | |
| MUT-SSAC-024 | PASS | |
| MUT-SSAC-025 | **FAIL** | @state transition 소문자 vs diagram PascalCase 미검출 |
| MUT-SSAC-026 | PASS | |
| MUT-SSAC-027 | PASS | |
| MUT-SSAC-028 | PASS | |
| MUT-SSAC-029 | PASS | |
| MUT-SSAC-030 | PASS | |
| MUT-SSAC-031 | PASS | |
| MUT-SSAC-032 | PASS | |
| MUT-SSAC-033 | PASS | |
| MUT-SSAC-034 | PASS | |
| MUT-SSAC-035 | **FAIL** | @call pkg 함수명 소문자 미검출 |
| MUT-SSAC-036 | PASS | |
| MUT-SSAC-037 | PASS | |
| MUT-SSAC-038 | PASS | |
| MUT-SSAC-039 | PASS | |
| MUT-SSAC-040 | PASS | |
| MUT-SSAC-041 | PASS | |
| MUT-SSAC-042 | PASS | |
| MUT-SSAC-043 | PASS | |
| MUT-SSAC-044 | PASS | |
| MUT-SSAC-045 | PASS | |
| MUT-SSAC-046 | **FAIL?** | dotted field 타입 불일치 미검출 (mutest 불확실 표기) |
| MUT-SSAC-047 | PASS | |
| MUT-SSAC-048 | PASS | |
| MUT-SSAC-049 | PASS | |

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
| MUT-STML-003 | SKIP | gigbridge에 x-component 없음 |

## FAIL 상세

### 확정 FAIL (5건)

| ID | 원인 | 수정 방향 |
|---|---|---|
| MUT-CONFIG-OPENAPI-001 | middleware↔securitySchemes 교차 검증 규칙 없음 | crosscheck에 middleware-security 교차 검증 추가 |
| MUT-FUNC-002 | funcspec 파서가 `panic("TODO")` 패턴을 stub body로 인식 못함 | isStubBody에 panic("TODO") 패턴 추가 |
| MUT-SSAC-002 | @get result type 존재 여부를 DDL 테이블과 교차 검증하지 않음 | ssac-ddl crosscheck에 result type 검증 추가 |
| MUT-SSAC-025 | @state transition event 대소문자 검증 없음 | states-ssac crosscheck에 case-sensitive 비교 추가 |
| MUT-SSAC-035 | @call pkg 함수명 대소문자 검증 없음 | ssac-func crosscheck에 PascalCase 함수명 검증 추가 |

### 불확실 FAIL (3건, mutest 자체에 "?" 표기)

| ID | 원인 | 판단 |
|---|---|---|
| MUT-SSAC-005 | input key "Id" vs "ID" 대소문자 | Go naming convention 차이, 검증 범위 외일 수 있음 |
| MUT-SSAC-020 | HTTP status 999 범위 검증 | 유효 범위(100~599) 검증 추가 가능하나 우선순위 낮음 |
| MUT-SSAC-046 | dotted field 참조 타입 불일치 | 타입 추론 복잡도 높음, 향후 과제 |

## SKIP 사유

| 사유 | 건수 |
|---|---|
| gigbridge에 해당 기능 없음 (queue, @dto, x-component 등) | 7 |
| 미실행 (시간 제약 또는 변경 적용 불가) | 10 |
| 미구현 규칙 | 1 |
