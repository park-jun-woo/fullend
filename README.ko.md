# fullend

9원소 이종(heterogeneous) 선언적 명세 매체 집합 {SSOT_i}_{i=1..9} 에 대하여
레이어 간 쌍별 정합 술어 ψ_{i,j} 의 전역 만족 ψ_{i,j}(SSOT_i, SSOT_j) = ⊤ 을
강제하는 동시에, 해당 집합을 Go/Gin 백엔드·TypeScript/React 프런트엔드·Hurl
프로브 코퍼스로 구성된 다언어 산출 표면 Γ(·) 에 결정적(deterministic)이고
비트 단위 재현 가능하도록 투사(project)하는 단일-바이너리 명세 기반
오케스트레이터이다.

## 1. 취지 (Abstract)

본 산출물은 Go 1.22+ 명령줄 프로그램으로서 그 관심사는 폐쇄세계 가정
(closed-world assumption) 하에서 9원소 선언적 명세 매체 파티션 전역에 걸친
참조 무결성(referential integrity)의 보존에 한정된다. 신규성이나 사용 편의
(ergonomics) 를 주장하지 않으며, 본 도구는 정점이 명세 종류(kind)이고 간선이
구문-의미 간 강제(coercion)인 범주론적 도식(diagram)의 기계적 실현에 그 존재
이유를 둔다. 프레임워크, 스캐폴드, 생산성 가속기를 탐색하는 독자는 다른
곳을 참조하기 바란다.

## 2. 기반 온톨로지 (Substrate Ontology)

적격(well-formed) 프로젝트의 정전(canonical) 배치는 `specs/` 디렉토리로서
그 직계 자식이 아래 9원소 집합을 구성한다. 선택적 `func/` 를 제외한 여하한
구성원의 결손은 §4.1 하 진단 가능한 불일치를 구성한다.

```
specs/
├── fullend.yaml                     ∈ Σ_config
├── api/openapi.yaml                 ∈ Σ_openapi            (OpenAPI 3.x; OAS §4 참조)
├── db/*.sql                         ∈ Σ_ddl ∪ Σ_query      (PostgreSQL DDL + sqlc 질의)
├── service/**/*.ssac                ∈ Σ_ssac               (SSaC: 주석 주도 시퀀스 DSL)
├── model/*.go                       ∈ Σ_model              (Go 구조체 선언; //@dto 표지)
├── func/<pkg>/*.go                  ∈ Σ_func               (선택; @call 대상 사용자 구현)
├── states/*.md                      ∈ Σ_fsm                (Mermaid stateDiagram)
├── policy/*.rego                    ∈ Σ_rego               (OPA Rego)
├── tests/{scenario,invariant}-*.hurl ∈ Σ_scenario           (Hurl 코퍼스)
└── frontend/*.html                  ∈ Σ_stml               (STML: HTML5 + data-* DSL)
```

## 3. 획득 (Acquisition)

```
go install github.com/park-jun-woo/fullend/cmd/fullend@latest
```

런타임 의존은 내장되지 않는다. 다음 동반 실행 파일이 생성 시점의 `PATH`
상에서 해소 가능해야 한다: `oapi-codegen` (v2.x), `sqlc` (>= 1.25.0),
`hurl` (>= 4.0, 하류 프로브 검증 용도). 내부 버전 고정(pinning) 전략은
제공되지 않는다.

## 4. 연산 동사 (Operational Verbs)

절차는 4원(quaternary) 동사로 매개화된다. 모든 동사는 단조 파싱 단계
`ParseAll` 을 공유하며, 단일 프로세스 수명 내 반복 호출에 대한 멱등성
(idempotence) 을 핵심 속성으로 한다.

### 4.1 `validate`

(a) 각 레이어 파서에 위임되는 레이어 내(intra-layer) 구문 검사, (b) 9원소에
걸친 레이어 간(cross-layer) 술어 ψ_{i,j} 의 이행 폐포(transitive closure),
(c) 선재(pre-existing) `artifacts/` 트리와의 계약-다이제스트 대조를 순차
수행한다. `--skip k,...` 수식자는 `k ∈ {openapi, ddl, ssac, model, stml,
states, policy, scenario, func}` 에 속하는 종류를 두 단계에서 공히 생략하되,
그 부재를 `Pass`/`Fail` 이 아닌 `Skip` 으로 기록한다.

### 4.2 `gen`

`validate` 와 결정적 산출 합성 절차 Γ 의 합성(composition)이며, 그 상
(image)은 `artifacts/` 트리이다. Γ 는 §3 의 외부 도구를 모듈로 순수(pure)
하게 정의되며 비트 수준 재현 가능성은 휴리스틱이 아닌 불변량(invariant)
이다. 이에 대한 위반은 회귀(regression)를 구성하며 그에 따라 조치된다.

### 4.3 `status`

순전히 정보성(informational) 집계를 산출한다. 파일 시스템 부작용 없음.

### 4.4 `chain`

입력으로 OpenAPI `operationId` 를 받아 그에 대한 이행 연결 폐포(transitive
connectivity closure) 에 속하는 SSOT 및 산출 노드 집합을 (kind, file, line,
summary, ownership) 5원조로 방출한다. 일상적 사용이 아닌 사후(post-hoc)
영향 분석을 의도한다.

### 4.5 `gen-model` (보조)

OpenAPI 문서(파일 또는 URI)를 수용하여 지정 출력 디렉토리 하에 Go HTTP
클라이언트 패키지(`package external`)를 산출하는 직교(orthogonal) 동사.
`gen` 과는 OpenAPI 적재기(loader)를 제외하면 코드 경로를 공유하지 아니하며
내부 DDL 주도 모델 합성 파이프라인과 무관하다.

## 5. 레이어 간 술어 열거 (Cross-Layer Predicate Enumeration)

아래는 위반 시 진단 가능한 불일치를 구성하는 주요 ψ_{i,j} 술어의 비망라적
열거이다. 술어의 arity 와 한정자(quantifier) 구조는 생략되며, 규범적
정식화(formulation)는 `pkg/crosscheck/` 참조.

- ψ(config, openapi):       미들웨어 식별자와 `components.securitySchemes` 의 합치
- ψ(openapi, ddl):           x-sort / x-filter 컬럼의 존재; x-include → 테이블 매핑
- ψ(ssac, ddl):              @result 타이핑 대 DDL 유도 구조 정의역; arg ↔ column 전사성
- ψ(fsm, ssac):               전이 이벤트 ↔ ServiceFunc 의 가드 함수 상에서의 전단사
- ψ(fsm, ddl):                상태 컬럼 정의역의 매장(embedding)
- ψ(fsm, openapi):            전이 이벤트 ↔ operationId 합치
- ψ(rego, ssac):              Rego allow 규칙 전항에서의 (action, resource) 발현
- ψ(rego, ddl):               @ownership 테이블/컬럼 존재
- ψ(rego, fsm):               @auth 부가 전이에 대한 allow 규칙 피복(coverage)
- ψ(scenario, openapi):       엔드포인트 존재
- ψ(queue):                    @publish ↔ @subscribe 토픽 합치; 페이로드 구조적 일치
- ψ(func, ssac):              @call 차수, 위치 타이핑, result/response 합치
- ψ(stml, ssac):              operationId 공참조 매개

임의 두 종류 간 술어의 부재는 의도된 비결합(intentional decoupling)을
뜻하며 누락을 의미하지 않는다.

## 6. 산출 표면 (Γ)

Γ 는 상호 서로소(disjoint)인 3 산출 하위 기반(substratum)을 표적으로 한다.

- **Go/Gin 백엔드** — SSaC→Go 핸들러 합성기, oapi-codegen 중개의 타입/서버
  골격, sqlc 중개의 질의 레이어, 그리고 트리 내 feature 군집(feature-grouped)
  Handler/Server 구성기의 협조로 산출된다. feature 군집은 `specs/service/`
  의 직계 서브디렉토리에 의해 유도된다.
- **TypeScript/React 프런트엔드** — STML→TSX 페이지 합성기와 OpenAPI 유래의
  최소 글루 표면(`App.tsx`, `main.tsx`, `api.ts`) 의 결합으로 산출된다.
- **Hurl 프로브 코퍼스** — OpenAPI-×-FSM-×-정책의 곱으로부터 유도되며,
  리소스 및 상태 전이 의존성 상의 위상 정렬 스케줄에 따라 순서화된
  결정적 smoke 시퀀스.

Γ 는 컴파일 가능성에 요구되는 최소치를 넘어 ORM 레이어, 서버 프레임워크,
빌드 도구, 번들러 설정의 산출을 의도적으로 삼간다.

## 7. 내장 호출 대상 (pkg/)

SSaC `@call` 지점에서 사용 가능한 고정 호출 대상 집합이 벤더링되어 있다.
구현의 안정성(stability)은 미지정(unspecified)이다. 대안적 의미론을
요구하는 프로젝트는 `specs/<project>/func/<pkg>/` 하에 섀도잉(shadowing)
구현을 제공해야 한다.

네임스페이스: `auth` (bcrypt/JWT/재설정), `crypto` (AES-256-GCM, TOTP),
`storage` (S3 호환), `mail` (SMTP), `text` (유니코드 슬러그, HTML 정제,
자소(grapheme) 고려 절단), `image` (OG/썸네일 래스터화).

## 8. 내장 모델 인터페이스 (pkg/)

DDL 로부터 유도 불가능한 I/O 를 위하여 `SessionModel`, `CacheModel`,
`FileModel`, Pub/Sub 싱글턴이 패키지 범위 `@model` 인터페이스로 제공된다.
백엔드 선택(PostgreSQL, 인메모리, S3, 로컬 디스크)은 `fullend.yaml` 에
위임된다. 규범적 계약은 각 해당 패키지 참조. 편의 문서는 여기에 중복되지
않는다.

## 9. 레이어 간 검증 근거

레이어 내 검증기(SSaC, STML 등)는 필요 조건이나 충분 조건이 아니다.
fullend 의 역할은 레이어의 데카르트 곱 전역에 걸친 정합성의 강제이며,
개별 레이어의 적격성(well-formedness)은 전제된다. 레이어 내 부적격에
대한 위반은 책임 파서의 진단 소유권(diagnostic ownership)으로 보고되며
fullend 가 합성하지 않는다.

## 10. 런타임 프로빙

산출된 `artifacts/tests/smoke.hurl` 은 생성된 백엔드의 기동 인스턴스에
대하여 Hurl 로 실행 가능하다. 백엔드 프로세스의 오케스트레이션은 제공되지
않으며, 그 기동과 해제는 사용자의 관심사이다.

```
hurl --test --variable host=http://localhost:8080 artifacts/<project>/tests/smoke.hurl
```

`specs/<project>/tests/` 하에 저작된 scenario 및 invariant 코퍼스는 산출
트리로 원형 보존되어 전달된다.

## 11. 구조적 비망 (Architectural Notes)

SSaC 와 STML 은 역사적으로 독립 저장소로 유지되었으나 현재 본 트리에
`internal/ssac/` 및 `internal/stml/` 로 융합되어 있다. 상류 저장소
(`park-jun-woo/ssac`, `park-jun-woo/stml`)는 본 트리의 하위 트리의
파일 단위 복제본(mirror)이며 독립적 진화 궤적을 보유하지 아니한다.
모든 SSOT 의 획득은 단일 `ParseAll()` 진입점을 거쳐 공유 `ParsedSSOTs`
구조를 구체화(materialize)하며 모든 동사가 이를 소비한다.

## 12. 감사 (Acknowledgments)

본 도구의 존재는 다음에 의존한다: OpenAPI Initiative, sqlc, Open Policy
Agent, Mermaid, oapi-codegen, kin-openapi, Hurl, React, React Router,
TanStack Query, React Hook Form, Vite, Tailwind CSS, TypeScript, Gin,
`lib/pq`. 전기(前記) 중 어느 것에 대해서도 기여, 파생, 경쟁의 주장이나
함의는 없다.

## 13. License

MIT. `LICENSE` 참조.
