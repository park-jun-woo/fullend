# Fullend Validate 검증 현황

## SSOT 소비 관계

```
DDL ←── OpenAPI ←── Config (middleware → securitySchemes)
 ↑         ↑
 ├── SSaC ─┘←── Policy (action/resource → @auth)
 │    ↑              ↑
 │    ├── States     ├── DDL (@ownership)
 │    ├── Func       ├── Config (roles, claims)
 │    └── Config     └── States
 │
 ├── States (state field → DDL column)
 └── DDL → SSaC (coverage)

Scenario → OpenAPI
STML → OpenAPI (stml/validator에서 수행, Phase029에서 crosscheck 이동 예정)
```

### 단방향

| 소비자 | 대상 | Check 함수 | 검증 내용 |
|---|---|---|---|
| OpenAPI | DDL | `CheckOpenAPIDDL` | property↔column, x-include/x-sort/x-filter↔FK/index, ghost/missing property |
| SSaC | DDL | `CheckSSaCDDL` | @result type↔DDL table, param type↔DDL column |
| SSaC | Func | `CheckFuncs` | @call pkg.Function 존재, 시그니처 일치 |
| SSaC | Func | `CheckAuthz` | @auth inputs↔pkg/authz CheckRequest fields |
| SSaC | Config | `CheckClaims` | currentUser.X↔claims 정의 |
| Policy | SSaC | `CheckPolicy` | Rego action/resource↔SSaC @auth |
| Policy | DDL | `CheckPolicy` | @ownership table.column↔DDL column |
| Policy | Config | `CheckClaimsRego` | Rego input.claims.X↔claims 정의 |
| Policy | Config | `CheckRoles` | Rego input.role↔roles 목록 |
| Policy | States | `CheckPolicy` | Rego 상태 참조↔States 정의 |
| Config | OpenAPI | `CheckMiddleware` | middleware↔securitySchemes |
| Scenario | OpenAPI | `CheckHurlFiles` | Hurl path/method↔OpenAPI endpoint |

### 양방향

| 쌍 | Check 함수 | A→B | B→A |
|---|---|---|---|
| SSaC ↔ OpenAPI | `CheckSSaCOpenAPI` | func name→operationId, @response→response schema, @empty/@exists→error response | operationId→SSaC func 존재 |
| SSaC ↔ States | `CheckStates` | @state diagramID→diagram 존재, func name→유효 전이 | transition event→SSaC func 존재, event without @state→WARNING |
| Func → SSaC | `CheckFuncCoverage` | func spec→SSaC @call 참조 존재 (coverage) | — |
| SSaC ↔ DDL | `CheckDDLCoverage` | (CheckSSaCDDL에서 커버) | DDL table→SSaC에서 사용됨 (coverage) |
| States ↔ DDL | `CheckStates` | @state input field→DDL column | (역방향 없음) |

### 단독 검증

| SSOT | Check 함수 | 검증 내용 |
|---|---|---|
| DDL | `CheckSensitiveColumns` | 칼럼명 sensitive 패턴 매칭 + @sensitive 어노테이션 대조 |
| SSaC | `CheckQueue` | @publish topic↔@subscribe topic 내부 일관성 |
| States | `statemachine.Parse` | 상태명 case-insensitive 중복 검출 |

---

## 개선안

### 완료된 개선 (Phase027)

- ~~**States → OpenAPI 제거**~~ — 전이적 중복 제거 완료 (CheckStates #4 삭제)
- ~~**Func → SSaC 커버리지**~~ — `CheckFuncCoverage` 추가 완료 (WARNING 수준)

### 추가 후보 (1건)

**STML → OpenAPI 크로스체크 (Phase029)**

현황: STML validator가 내부에서 OpenAPI 교차 검증을 수행 중. crosscheck에 중복 구현하면 이중 출력.

계획: STML validator의 OpenAPI 의존성을 crosscheck로 이동 (Phase029-STMLCrosscheck).
