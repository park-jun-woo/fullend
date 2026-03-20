# pkg/crosscheck

Toulmin defeats graph 기반 SSOT 간 교차 검증.

## 규칙 목록

| 규칙 | Source | Target | 설명 |
|------|--------|--------|------|
| OpenAPI ↔ DDL | OpenAPI | DDL | x-sort/filter/include 필드, 유령/누락 property |
| SSaC ↔ DDL | SSaC | DDL | @result/@param 타입 일치 |
| SSaC input key case | SSaC | DDL | 입력 키 snake_case 일관성 |
| SSaC ↔ OpenAPI | SSaC | OpenAPI | 함수명=operationId, @response 필드 매칭 |
| States ↔ SSaC/DDL | States | SSaC | 상태 전이 → 함수 매핑, 가드 DDL 참조 |
| Policy ↔ SSaC/DDL | Policy | SSaC | Rego 규칙 → 유효한 리소스/컬럼 참조 |
| Scenario → OpenAPI | Scenario | OpenAPI | Hurl 경로 → API 엔드포인트 매칭 |
| SSaC → Func | SSaC | Func | @func/@call → func spec 매칭, 입출력 타입 |
| SSaC → JWT Claims | SSaC | Config | JWT 클레임 참조 일치 |
| Config → OpenAPI | Config | OpenAPI | 미들웨어 → 엔드포인트 참조 |
| SSaC → Config | SSaC | Config | 클레임 사용 → fullend.yaml 일치 |
| Policy → Config (claims) | Policy | Config | Rego claims → fullend.yaml 일치 |
| DDL Coverage | DDL | SSaC | 모든 테이블이 SSaC에서 참조되는지 |
| Queue | SSaC | Config | publish ↔ subscribe 토픽 짝 |
| SSaC → Authz | SSaC | Func | @auth 입력 → CheckRequest 필드 |
| DDL Sensitive | DDL | — | 민감 컬럼 어노테이션 검사 |
| Func Coverage | Func | SSaC | 모든 func spec이 @call에서 사용되는지 |
| Policy → Config (roles) | Policy | Config | Rego 역할 → fullend.yaml 일치 |
| Policy → DDL (roles) | Policy | DDL | DDL 제약조건 역할 → 정책 일치 |
| OpenAPI Constraints | OpenAPI | DDL | 제약조건 어노테이션 → DDL 일치 |

## Toulmin 매핑

```
claim   = 검증 대상 (ServiceFunc, DDLTable, Operation, ...)
ground  = 교차 검증에 필요한 타 SSOT 데이터 (Fullstack + SymbolTable)
backing = 검증 기준/설정 (claims, roles, queue config 등)

warrant  = 기본 규칙 ("타입이 일치해야 한다")
rebuttal = 예외 ("pkg 모델이면 DDL 불필요")
```

## 검증 흐름

```
Fullstack (파싱 결과) + SymbolTable
  → 규칙별 Toulmin Graph 구성
  → Graph.Evaluate(claim, ground) per SSOT 항목
  → verdict + evidence 반환
```
