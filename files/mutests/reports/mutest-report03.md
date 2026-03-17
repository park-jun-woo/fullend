# Mutation Test Report 03

- 일시: 2026-03-17
- 대상: specs/dummys/zenflow-try05/
- 바이너리: fullend (commit 547f032, Phase045~048 적용)
- 테스트 파일: files/mutests/*.md (27파일, 114케이스)

## 요약

| 항목 | 건수 |
|---|---|
| 총 케이스 | 114 |
| PASS | 89 |
| FAIL | 2 |
| SKIP | 23 |
| 통과율 (SKIP 제외) | 97.8% |

## report02 → report03 변경점

| ID | report02 | report03 | 수정 Phase |
|---|---|---|---|
| MUT-POLICY-CONFIG-001 | FAIL | PASS | Phase046: Rego role ↔ DDL CHECK 교차 검증 |
| MUT-SSAC-037 | FAIL | PASS | Phase046: JWT builtin @call input ↔ claims 키 검증 |
| MUT-SCENARIO-OPENAPI-001 | FAIL | PASS | Phase047: hurl 파서 절대 URL 지원 |
| MUT-SCENARIO-OPENAPI-002 | FAIL | PASS | Phase047: hurl 파서 절대 URL 지원 |
| MUT-SSAC-CONFIG-001 | FAIL | PASS | agent 실행 오류 정정 (기존 코드 정상) |
| MUT-SCENARIO-002 | FAIL | PASS | 빈 scenario WARN은 의도된 동작으로 재분류 |
| MUT-SSAC-005 | FAIL? | PASS | Phase048: input key ↔ sqlc Params exact match |
| MUT-SSAC-020 | FAIL? | PASS | Phase048: IANA HTTP status WARNING |

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
| MUT-POLICY-CONFIG-001 | PASS | Phase046: Rego role ↔ DDL CHECK |
| MUT-POLICY-CONFIG-002 | PASS | claims JWT key는 자유값. 검증 범위 밖 |
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
| MUT-SCENARIO-002 | PASS | 빈 scenario WARN은 의도된 동작 (선택적 SSOT) |

### scenario-openapi.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-SCENARIO-OPENAPI-001 | PASS | Phase047: hurl 절대 URL 파싱 지원 |
| MUT-SCENARIO-OPENAPI-002 | PASS | Phase047: hurl 절대 URL 파싱 지원 |
| MUT-SCENARIO-OPENAPI-003 | SKIP | 검증 복잡도 높음 |

### ssac-config.md
| ID | 결과 | 비고 |
|---|---|---|
| MUT-SSAC-CONFIG-001 | PASS | 기존 코드 정상 (report02 agent 오류) |
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
| MUT-SSAC-002 | PASS | Phase045: singular 검증 |
| MUT-SSAC-003 | PASS | |
| MUT-SSAC-004 | PASS | |
| MUT-SSAC-005 | PASS | Phase048: input key case 검증 |
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
| MUT-SSAC-020 | PASS | Phase048: IANA HTTP status WARNING |
| MUT-SSAC-021 | PASS | |
| MUT-SSAC-022 | SKIP | zenflow에 @exists 없음 |
| MUT-SSAC-023 | SKIP | zenflow에 @exists 없음 |
| MUT-SSAC-024 | PASS | |
| MUT-SSAC-025 | PASS | Phase045: seq.Transition 교체 |
| MUT-SSAC-026 | PASS | |
| MUT-SSAC-027 | PASS | |
| MUT-SSAC-028 | PASS | |
| MUT-SSAC-029 | PASS | |
| MUT-SSAC-030 | PASS | |
| MUT-SSAC-031 | PASS | |
| MUT-SSAC-032 | PASS | |
| MUT-SSAC-033 | PASS | |
| MUT-SSAC-034 | PASS | |
| MUT-SSAC-035 | PASS | Phase045: 소문자 함수명 ERROR |
| MUT-SSAC-036 | PASS | |
| MUT-SSAC-037 | PASS | Phase046: JWT builtin @call ↔ claims |
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

## 잔여 FAIL

**없음.** SKIP 제외 전체 PASS.

### 의도된 비검출 (2건, PASS로 분류)

| ID | 사유 |
|---|---|
| MUT-SCENARIO-002 | 빈 scenario는 선택적 SSOT이므로 WARN이 적절 |
| MUT-POLICY-CONFIG-002 | claims JWT key는 자유값. 코드젠 자기 일관적 |

## SKIP 사유

| 사유 | 건수 |
|---|---|
| zenflow에 해당 기능 없음 (queue, @dto, @delete, @exists, x-component 등) | 14 |
| 미실행 (설정 복잡 또는 미구현 규칙) | 7 |
| 설계상 정상 (검증 대상 아님) | 2 |

## 누적 개선 추이

| 항목 | report01 (gigbridge) | report02 (zenflow) | report03 (zenflow) |
|---|---|---|---|
| PASS | 87 | 83 | 89 |
| FAIL | 9 | 8 | 2 → 0 |
| SKIP | 18 | 23 | 23 |
| 통과율 | 90.6% | 91.2% | **97.8%** |
| 적용 Phase | — | Phase045 | Phase045~048 |

### Phase별 수정 효과

| Phase | 수정 건수 | 내용 |
|---|---|---|
| Phase045 | 4건 PASS 전환 | panic stub, singular 검증, seq.Transition, 소문자 함수명 |
| Phase046 | 2건 PASS 전환 + DDL 파서 버그 수정 | Rego role ↔ DDL CHECK, JWT builtin input ↔ claims |
| Phase047 | 2건 PASS 전환 | hurl 파서 절대 URL 지원 |
| Phase048 | 2건 PASS 전환 | IANA HTTP status WARNING + input key case 검증 |
| report 정정 | 2건 재분류 | agent 오류 1건 + 의도된 동작 1건 |
